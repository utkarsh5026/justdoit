package objects

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	ShaSize = 20 // The size of SHA-1 hash in Git tree.
)

// GitTreeLeaf represents a single entry in a Git tree object.
//
// A Git tree object is a fundamental component of Git's object model, representing the structure of a directory in a Git repository.
// Each entry in a tree object can be a file (blob), another directory (tree), a symbolic link, or a submodule.
//
// Fields:
// - mode: A string representing the file mode and type. Here The first two or three digits represent the type of the object and the last three digits represent the Unix file permissions. Common modes include:
//   - "100644": Regular file with read/write permissions for the owner, and read-only permissions for others.
//   - "100755": Executable file with read/write/execute permissions for the owner, and read/execute permissions for others.
//   - "040000": Directory (tree).
//   - "120000": Symbolic link.
//   - "160000": Gitlink (submodule).
//
// - sha: A string representing the 20-byte SHA-1 hash of the object this entry points to. This hash uniquely identifies the object in the Git repository.
// - path: A string representing the file or directory name relative to the root of the repository.
type GitTreeLeaf struct {
	mode string
	sha  string
	path string
}

type GitTree struct {
	entries []*GitTreeLeaf
}

func (tr *GitTreeLeaf) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("mode: %s - ", tr.mode))
	builder.WriteString(fmt.Sprintf("sha: %s - ", tr.sha))
	builder.WriteString(fmt.Sprintf("path: %s", tr.path))
	return builder.String()
}

// Type determines the Git object type of the tree leaf based on its mode.
//
// The function extracts the type string from the mode and maps it to a GitObjectType.
// It returns an error if the mode does not correspond to a valid object type.
//
// Returns:
// - GitObjectType: The type of the Git object (e.g., BlobType, TreeType, CommitType).
// - error: An error if the mode does not correspond to a valid object type.
func (tr *GitTreeLeaf) Type() (GitObjectType, error) {
	var typeStr string
	if len(tr.mode) == 5 {
		typeStr = tr.mode[:1]
	} else {
		typeStr = tr.mode[:2]
	}

	switch typeStr {
	case "10":
		return BlobType, nil
	case "04":
		return TreeType, nil
	case "16":
		return CommitType, nil
	case "12":
		return BlobType, nil

	default:
		return 0, fmt.Errorf("invalid object type: %s", typeStr)
	}
}

// Sha returns the SHA-1 hash of the tree leaf.
func (tr *GitTreeLeaf) Sha() string {
	return tr.sha
}

// Name returns the name of the tree leaf.
func (tr *GitTreeLeaf) Name() string {
	return tr.path
}

// Mode returns the mode of the tree leaf.
func (tr *GitTreeLeaf) Mode() string {
	return tr.mode
}

func Tree() *GitTree {
	return &GitTree{}
}

func (tr *GitTree) Serialize() ([]byte, error) {
	return treeSerialize(tr)
}

func (tr *GitTree) Deserialize(raw []byte) error {
	leaves, err := parseTree(raw)
	if err != nil {
		return err
	}

	tr.entries = leaves
	return nil
}

func (tr *GitTree) Format() GitObjectType {
	return TreeType
}

func (tr *GitTree) SetData(data []byte) {
	tr.entries = nil
	entries, err := parseTree(data)

	if err != nil {
		return
	}

	tr.entries = entries
}

func (tr *GitTree) Entries() []*GitTreeLeaf {
	return tr.entries
}

var InvalidTreeEntry = func(problem string) error {
	return fmt.Errorf("invalid tree entry: %s", problem)
}

// parseSingleTreeEntry parses a single tree entry from a raw byte slice starting at a given position.
//
// Tree entries in Git are defined as follows:
// - Mode: A string representing the file mode (e.g., "100644" for a regular file).
// - Path: The file path relative to the root of the repository.
// - SHA-1: A 20-byte SHA-1 hash of the object.
//
// The format of a tree entry is: <mode><space><path>\0<sha-1>
//
// Parameters:
// - raw: A byte slice containing the raw tree data.
// - start: An integer representing the starting position in the raw byte slice.
//
// Returns:
// - int: The position in the byte slice after the parsed entry.
// - *GitTreeLeaf: A pointer to a GitTreeLeaf struct containing the parsed mode, SHA-1, and path.
// - error: An error if the tree entry is invalid.
func parseSingleTreeEntry(raw []byte, start int) (int, *GitTreeLeaf, error) {
	x := bytes.IndexByte(raw[start:], Space)
	if x < 0 {
		return 0, nil, InvalidTreeEntry("missing mode")
	}

	if x != 5 && x != 6 {
		return 0, nil, InvalidTreeEntry("invalid mode length")
	}

	mode := raw[start : start+x]
	if len(mode) == 5 {
		mode = append([]byte{Space}, mode...)
	}

	nullIdx := bytes.IndexByte(raw[start+x+1:], 0)
	if nullIdx < 0 {
		return 0, nil, InvalidTreeEntry("missing null terminator")
	}

	nullIdx += start + x + 1 // Adjust for the slice
	path := string(raw[start+x+1 : nullIdx])
	if len(raw[nullIdx+1:]) < ShaSize {
		return 0, nil, InvalidTreeEntry("missing sha")
	}

	shaStart := nullIdx + 1
	sha := hex.EncodeToString(raw[shaStart : shaStart+ShaSize])

	return shaStart + ShaSize, &GitTreeLeaf{string(mode), sha, path}, nil
}

// parseTree parses a raw byte slice containing multiple tree entries and returns a slice of GitTreeLeaf pointers.
//
// Parameters:
// - raw: A byte slice containing the raw tree data.
//
// Returns:
// - []*GitTreeLeaf: A slice of pointers to GitTreeLeaf structs containing the parsed mode, SHA-1, and path for each entry.
// - error: An error if any of the tree entries are invalid.
func parseTree(raw []byte) ([]*GitTreeLeaf, error) {
	var leaves []*GitTreeLeaf
	var start int

	for start < len(raw) {
		end, leaf, err := parseSingleTreeEntry(raw, start)
		if err != nil {
			return nil, err
		}

		leaves = append(leaves, leaf)
		start = end
	}

	return leaves, nil
}

// treeSerialize serializes a GitTree into a byte slice.
//
// The function sorts the tree entries based on their paths, with directories
// sorted after files. It then constructs a byte buffer containing the serialized
// tree entries in the format: <mode><space><path>\0<sha-1>
//
// Parameters:
// - tree: A pointer to a GitTree struct containing the tree entries to be serialized.
//
// Returns:
// - []byte: A byte slice containing the serialized tree data.
// - error: An error if any of the tree entries are invalid.
func treeSerialize(tree *GitTree) ([]byte, error) {
	sortFunc := func(leaf *GitTreeLeaf) string {
		if leaf.mode[:2] == "10" {
			return leaf.path
		}

		return leaf.path + "/"
	}

	sort.Slice(tree.entries, func(i, j int) bool {
		return sortFunc(tree.entries[i]) < sortFunc(tree.entries[j])
	})

	var buffer bytes.Buffer

	for _, item := range tree.entries {
		buffer.WriteString(item.mode)
		buffer.WriteByte(Space)
		buffer.WriteString(item.path)
		buffer.WriteByte(0)

		sha, err := hex.DecodeString(item.sha)
		if err != nil {
			return nil, err
		}

		buffer.Write(sha)
	}

	return buffer.Bytes(), nil
}
