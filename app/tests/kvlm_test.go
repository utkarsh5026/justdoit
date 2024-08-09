package test

import (
	"bytes"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
	"reflect"
	"testing"
)

func TestKvlmSerialize(t *testing.T) {
	tests := []struct {
		name     string
		input    *ordereddict.OrderedDict
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    ordereddict.New(),
			expected: []byte{},
			wantErr:  false,
		},
		{
			name: "Single key-value pair",
			input: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				return od
			}(),
			expected: []byte("key1 value1\n"),
			wantErr:  false,
		},
		{
			name: "Multiple key-value pairs",
			input: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				od.Set("key2", []byte("value2"))
				return od
			}(),
			expected: []byte("key1 value1\nkey2 value2\n"),
			wantErr:  false,
		},
		{
			name: "Multi-line value",
			input: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("line1\nline2"))
				return od
			}(),
			expected: []byte("key1 line1\n line2\n"),
			wantErr:  false,
		},
		{
			name: "With message",
			input: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				od.Set("", []byte("This is a message"))
				return od
			}(),
			expected: []byte("key1 value1\n\nThis is a message\n"),
			wantErr:  false,
		},
		{
			name: "Multiple values for same key",
			input: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", [][]byte{[]byte("value1"), []byte("value2")})
				return od
			}(),
			expected: []byte("key1 value1\nkey1 value2\n"),
			wantErr:  false,
		},
		{
			name: "Unsupported value type",
			input: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", 123) // Integer is not supported
				return od
			}(),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := objects.KvlmSerialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("kvlmSerialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("kvlmSerialize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestKvlmParse(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *ordereddict.OrderedDict
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: ordereddict.New(),
		},
		{
			name:  "Single key-value pair",
			input: []byte("key value\n"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key", []byte("value"))
				return od
			}(),
		},
		{
			name:  "Multiple key-value pairs",
			input: []byte("key1 value1\nkey2 value2\nkey3 value3\n"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				od.Set("key2", []byte("value2"))
				od.Set("key3", []byte("value3"))
				return od
			}(),
		},
		{
			name:  "Multi-line value",
			input: []byte("key1 value1\nkey2 line1\n line2\n line3\nkey3 value3\n"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				od.Set("key2", []byte("line1\nline2\nline3"))
				od.Set("key3", []byte("value3"))
				return od
			}(),
		},
		{
			name:  "Message at the end",
			input: []byte("key1 value1\nkey2 value2\n\nThis is a message"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				od.Set("key2", []byte("value2"))
				od.Set("", []byte("This is a message"))
				return od
			}(),
		},
		{
			name:  "Duplicate keys",
			input: []byte("key1 value1\nkey1 value2\nkey1 value3\n"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", [][]byte{[]byte("value1"), []byte("value2"), []byte("value3")})
				return od
			}(),
		},
		{
			name:  "No newline at end",
			input: []byte("key value"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key", []byte("value"))
				return od
			}(),
		},
		{
			name:  "Only message",
			input: []byte("\nThis is only a message"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("", []byte("This is only a message"))
				return od
			}(),
		},
		{
			name:  "Multi-line value without final newline",
			input: []byte("key1 value1\nkey2 line1\n line2\n line3"),
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key1", []byte("value1"))
				od.Set("key2", []byte("line1\nline2\nline3"))
				return od
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objects.KvlmParse(tt.input, 0, nil)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("KvlmParse() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestKvlmParse_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		start    int
		expected *ordereddict.OrderedDict
	}{
		{
			name:  "Start from middle",
			input: []byte("ignore this\nkey value\n"),
			start: 12,
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key", []byte("value"))
				return od
			}(),
		},
		{
			name:  "No newline at end",
			input: []byte("key value"),
			start: 0,
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("key", []byte("value"))
				return od
			}(),
		},
		{
			name:  "Only message",
			input: []byte("\nThis is only a message"),
			start: 0,
			expected: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("", []byte("This is only a message"))
				return od
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objects.KvlmParse(tt.input, tt.start, nil)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("KvlmParse() = %v, want %v", result, tt.expected)
			}
		})
	}
}
