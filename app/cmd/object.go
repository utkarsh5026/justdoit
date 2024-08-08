package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
)

type GitObjectType uint

const (
	BlobType GitObjectType = iota
	CommitType
	TreeType
	TagType
)

// String returns the string representation of the GitObjectType.
func (got GitObjectType) String() string {
	switch got {
	case BlobType:
		return "blob"
	case CommitType:
		return "commit"
	case TreeType:
		return "tree"
	case TagType:
		return "tag"
	default:
		return ""
	}
}

// typeFromString converts a string to a GitObjectType.
// If the string is not a valid object type, an error is returned.
func typeFromString(str string) (GitObjectType, error) {
	switch str {
	case "blob":
		return BlobType, nil
	case "commit":
		return CommitType, nil
	case "tree":
		return TreeType, nil
	case "tag":
		return TagType, nil
	default:
		return 0, fmt.Errorf("invalid object type: %s", str)
	}
}

type GitObject interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	Format() GitObjectType
}

// ObjectManager provides methods for reading and writing Git objects.
type ObjectManager struct {
	repo *GitRepository
}

// NewObjectManager creates a new ObjectManager with the given GitRepository.
func NewObjectManager(repo *GitRepository) *ObjectManager {
	return &ObjectManager{repo: repo}
}

// WriteObject serializes a GitObject, computes its SHA-1 hash, and writes it to the repository.
// It returns the SHA-1 hash of the written object or an error if the operation fails.
//
// Parameters:
// - obj: The GitObject to be written.
//
// Returns:
// - string: The SHA-1 hash of the written object.
// - error: An error if the operation fails.
func (om *ObjectManager) WriteObject(obj GitObject) (string, error) {
	data, err := obj.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize object: %w", err)
	}

	content := om.prepareObject(obj.Format(), data)
	sha := om.calculateSHA(content)
	path := getGitFilePath(om.repo, true, ObjectDir, sha[:2], sha[2:])

	if err := om.writeFile(path, content); err != nil {
		return "", fmt.Errorf("failed to write object: %w", err)
	}
	return sha, nil
}

// ReadObject reads a Git object from the repository using its SHA-1 hash.
// It returns the deserialized GitObject or an error if the operation fails.
//
// Parameters:
// - sha: The SHA-1 hash of the object to be read.
//
// Returns:
// - GitObject: The deserialized GitObject.
// - error: An error if the operation fails.
func (om *ObjectManager) ReadObject(sha string) (GitObject, error) {
	objectPath := getGitFilePath(om.repo, false, ObjectDir, sha[:2], sha[2:])
	content, err := om.readFile(objectPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	objectType, data, err := om.parseObject(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object: %w", err)
	}

	object, err := om.createObject(objectType)
	if err != nil {
		return nil, err
	}

	if err := object.Deserialize(data); err != nil {
		return nil, fmt.Errorf("failed to deserialize object: %w", err)
	}

	return object, nil
}

// prepareObject constructs the serialized Git object by adding the object type and size header.
func (om *ObjectManager) prepareObject(obType GitObjectType, data []byte) []byte {
	header := fmt.Sprintf("%s %d\x00", obType, len(data))
	return append([]byte(header), data...)
}

// calculateSHA computes the SHA-1 hash of the given content.
// It returns the hash as a hexadecimal string.
func (om *ObjectManager) calculateSHA(content []byte) string {
	hash := sha1.New()
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil))
}

// writeFile compresses the given content and writes it to the specified path.
// It returns an error if the operation fails.
//
// Parameters:
// - path: The file path where the content should be written.
// - content: The byte slice containing the content to be written.
//
// Returns:
// - error: An error if the operation fails.
func (om *ObjectManager) writeFile(path string, content []byte) error {
	var buff bytes.Buffer
	writer := zlib.NewWriter(&buff)

	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close object writer: %w", err)
	}

	return os.WriteFile(path, buff.Bytes(), 0644)
}

// readFile reads and decompresses the content from the specified file path.
// It returns the decompressed content as a byte slice or an error if the operation fails.
//
// Parameters:
// - path: The file path from which the content should be read.
//
// Returns:
// - []byte: The decompressed content read from the file.
// - error: An error if the operation fails.
func (om *ObjectManager) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer closeFile(file)

	var buff bytes.Buffer
	reader, err := zlib.NewReader(file)

	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			fmt.Println("failed to close zlib reader:", err)
		}
	}(reader)

	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}

	if _, err := io.Copy(&buff, reader); err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return buff.Bytes(), nil
}

// parseObject parses the given content to extract the Git object type and its data.
// It returns the GitObjectType, the object data as a byte slice, or an error if the operation fails.
//
// Parameters:
// - content: The byte slice containing the serialized Git object.
//
// Returns:
// - GitObjectType: The type of the Git object.
// - []byte: The data of the Git object.
// - error: An error if the operation fails.
func (om *ObjectManager) parseObject(content []byte) (GitObjectType, []byte, error) {
	nullIndex := bytes.IndexByte(content, 0)

	var ot GitObjectType
	if nullIndex == -1 {
		return ot, nil, fmt.Errorf("invalid object format")
	}

	header := string(content[:nullIndex])
	parts := bytes.SplitN([]byte(header), []byte(" "), 2)
	if len(parts) != 2 {
		return ot, nil, fmt.Errorf("invalid object header")
	}

	ot, err := typeFromString(string(parts[0]))

	if err != nil {
		return ot, nil, fmt.Errorf("invalid object type: %w", err)
	}

	size, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return ot, nil, fmt.Errorf("invalid object size: %w", err)
	}

	data := content[nullIndex+1:]
	if len(data) != size {
		return ot, nil, fmt.Errorf("object size mismatch")
	}

	return ot, data, nil
}

func (om *ObjectManager) createObject(ot GitObjectType) (GitObject, error) {
	switch ot {
	case BlobType:
		return Blob(), nil

	default:
		return nil, fmt.Errorf("unsupported object type: %s", ot)
	}
}

func (om *ObjectManager) FindObject(sha string, ot GitObjectType, follow bool) string {
	return sha
}
