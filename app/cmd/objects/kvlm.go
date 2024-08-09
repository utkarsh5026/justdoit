package objects

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
)

const (
	NewLine = '\n'
	Space   = ' '
)

// KvlmParse parses a raw byte slice starting from a given position and populates an OrderedDict with key-value pairs.
// It recursively processes the input to handle multi-line values and nested structures.
//
// Parameters:
// - raw: A byte slice containing the raw data to be parsed.
// - start: An integer representing the starting position in the raw byte slice.
// - dict: A pointer to an OrderedDict where the parsed key-value pairs will be stored. If nil, a new OrderedDict is created.
//
// Returns:
// - *ordereddict.OrderedDict: A pointer to the populated OrderedDict
func KvlmParse(raw []byte, start int, dict *ordereddict.OrderedDict) *ordereddict.OrderedDict {
	if dict == nil {
		dict = ordereddict.New()
	}

	if start >= len(raw) {
		return dict
	}

	// Find the next space and the next newline
	spc := bytes.IndexByte(raw[start:], ' ')
	nl := bytes.IndexByte(raw[start:], '\n')

	// Base case: if newline appears first (or there's no space at all)
	if spc < 0 || (nl >= 0 && nl < spc) {
		if nl != start {
			return dict
		}
		dict.Set("", raw[start+1:])
		return dict
	}

	// Recursive case: we read a key-value pair and recurse for the next
	key := string(raw[start : start+spc])

	// Find the end of the value
	end := start
	for {
		end = bytes.IndexByte(raw[end+1:], '\n')
		if end < 0 || len(raw) <= end+2 || raw[end+2] != ' ' {
			break
		}
		end += 1 // Adjust for the slice in IndexByte
	}

	// Handle case where no newline is found
	if end < 0 {
		end = len(raw)
	} else {
		end += start + 1
	}

	// Grab the value and drop the leading space on continuation lines
	value := bytes.Replace(raw[start+spc+1:end], []byte("\n "), []byte("\n"), -1)

	// Don't overwrite existing data contents
	if existingValue, exists := dict.Get(key); exists {
		switch v := existingValue.(type) {
		case [][]byte:
			dict.Set(key, append(v, value))
		default:
			dict.Set(key, [][]byte{v.([]byte), value})
		}
	} else {
		dict.Set(key, value)
	}

	// Recurse for the next key-value pair
	return KvlmParse(raw, end+1, dict)
}

// KvlmSerialize serializes the key-value list with message (kvlm) into a byte slice.
// It iterates over the OrderedDict and formats each key-value pair into a byte buffer.
// The message itself (if present) is appended at the end.
//
// Parameters:
// - kvlm: An OrderedDict containing the key-value pairs to be serialized.
//
// Returns:
// - []byte: A byte slice containing the serialized key-value pairs and message.
func KvlmSerialize(kvlm *ordereddict.OrderedDict) ([]byte, error) {
	if kvlm == nil {
		return nil, fmt.Errorf("input OrderedDict is nil")
	}
	var ret bytes.Buffer
	var errs []string

	// Output fields
	kvlm.Range(func(k string, v interface{}) bool {
		// Skip the message itself
		if k == "" {
			return true
		}
		switch val := v.(type) {
		case []byte:
			if err := writeKeyValue(&ret, k, val); err != nil {
				errs = append(errs, fmt.Sprintf("error writing key '%s': %v", k, err))
			}
		case [][]byte:
			for _, vb := range val {
				if err := writeKeyValue(&ret, k, vb); err != nil {
					errs = append(errs, fmt.Sprintf("error writing key '%s': %v", k, err))
				}
			}
		default:
			errs = append(errs, fmt.Sprintf("unsupported type for key '%s': %T", k, v))
		}

		return true
	})

	// Append message
	if message, exists := kvlm.Get(""); exists {
		if msg, ok := message.([]byte); ok {
			ret.WriteByte(NewLine)
			ret.Write(msg)
			ret.WriteByte(NewLine)
		} else {
			errs = append(errs, fmt.Sprintf("invalid message type: %T", message))
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("error serializing kvlm: %v", errs)
	}

	return ret.Bytes(), nil
}

// writeKeyValue writes a key-value pair to the provided buffer.
// It writes the key, followed by a space, the value, and a newline character.
// If any write operation fails, an error is returned.
//
// Parameters:
// - buffer: A pointer to a bytes.Buffer where the key-value pair will be written.
// - key: A string representing the key to be written.
// - value: A byte slice representing the value to be written.
//
// Returns:
// - error: An error if any write operation fails, otherwise nil.
func writeKeyValue(buffer *bytes.Buffer, key string, value []byte) error {
	if _, err := buffer.WriteString(key); err != nil {
		return fmt.Errorf("error writing key: %v", err)
	}
	if err := buffer.WriteByte(Space); err != nil {
		return fmt.Errorf("error writing space after key: %v", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(value))
	firstLine := true
	for scanner.Scan() {
		if !firstLine {
			if err := buffer.WriteByte(NewLine); err != nil {
				return fmt.Errorf("error writing newline: %v", err)
			}
			if err := buffer.WriteByte(Space); err != nil {
				return fmt.Errorf("error writing space: %v", err)
			}
		}
		if _, err := buffer.Write(scanner.Bytes()); err != nil {
			return fmt.Errorf("error writing value: %v", err)
		}
		firstLine = false
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning value: %v", err)
	}

	if err := buffer.WriteByte(NewLine); err != nil {
		return fmt.Errorf("error writing final newline: %v", err)
	}

	return nil
}
