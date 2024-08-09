package objects

import (
	"bytes"
	"fmt"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
	"strconv"
	"time"
)

// GitSignature represents the signature information in a Git commit.
// In Git, signatures are used to identify the author and committer of a commit.
//
// The signature in a Git commit appears in the following format:
//
//	name <email> timestamp timezone
//
// For example:
//
//	John Doe john@example.com 1623456789 +0100
type GitSignature struct {
	Name  string    // The name of the author or committer.
	Email string    // The email address of the author or committer.
	When  time.Time // The timestamp of the commit.
}

// GitCommit represents a Git commit object. It contains all the essential
// information stored in a Git commit.
//
// In an actual Git commit, the information is stored in the following format:
//
//	tree <sha1>
//	parent <sha1>
//	author <name> <email> <timestamp> <timezone>
//	committer <name> <email> <timestamp> <timezone>
//
//	<commit message>
type GitCommit struct {
	// Tree is the SHA-1 hash of the tree object associated with this commit.
	Tree string

	// Parents is a slice of SHA-1 hashes of the parent commit(s). Most commits have one parent, merge commits have multiple.
	Parents []string

	// Author is the signature of the person who originally created the work in the commit.
	Author GitSignature

	// Committer is the signature of the person who actually created the commit. Often the same as the author, but can differ in cases like applying patches.
	Committer GitSignature

	// Message is the commit message that describes the changes made in the commit.
	Message string
}

var InvalidSignature = func(sign string) error {
	return fmt.Errorf("invalid signature %s", sign)
}

var CommitKeyMissing = func(key string) error {
	return fmt.Errorf("missing key %s in commit object", key)
}

// CommitObject represents a Git commit object in the Git object model.
// A commit object contains metadata about the commit, such as the author, committer, and commit message.
// All the metadata is stored in a key-value list.
//
// Fields:
// - kvlm: An OrderedDict that stores key-value pairs of the commit metadata
// - commit: A pointer to a GitCommit struct that represents the commit metadata
type CommitObject struct {
	kvlm   *ordereddict.OrderedDict
	commit *GitCommit
}

func Commit() *CommitObject {
	return &CommitObject{
		kvlm: ordereddict.New(),
	}
}

func (c *CommitObject) Serialize() ([]byte, error) {
	return KvlmSerialize(c.kvlm)
}

func (c *CommitObject) Deserialize(data []byte) error {
	c.kvlm = KvlmParse(data, 0, c.kvlm)
	commit, err := createCommitFromKVLM(c.kvlm)

	if err != nil {
		return err
	}
	c.commit = commit
	return nil
}

func (c *CommitObject) Format() GitObjectType {
	return CommitType
}

func (c *CommitObject) SetData(data []byte) {
	c.kvlm = KvlmParse(data, 0, nil)
}

// parseSignature parses a Git signature from a byte slice.
//
// A Git signature consists of a name, an email, and a timestamp. This function splits the input byte slice
// into these components and returns a GitSignature struct.
//
// Parameters:
// - sign: A byte slice containing the signature to be parsed.
//
// Returns:
// - A GitSignature struct containing the parsed name, email, and timestamp.
// - An error if the signature is invalid or if the timestamp cannot be parsed.
func parseSignature(sign []byte) (GitSignature, error) {
	var gitSign GitSignature

	parts := bytes.Split(sign, []byte{Space})
	if len(parts) < 3 {
		return gitSign, InvalidSignature(string(sign))
	}

	cnt := len(parts)
	name := string(bytes.Join(parts[:cnt-2], []byte{Space})) // Join the parts except the last two
	email := string(parts[cnt-2])
	email = email[1 : len(email)-1] // Remove the angle brackets

	timestamp, err := strconv.ParseInt(string(parts[cnt-1]), 10, 64)

	if err != nil {
		return gitSign, err
	}

	gitSign = GitSignature{
		Name:  name,
		Email: email,
		When:  time.Unix(timestamp, 0),
	}
	return gitSign, nil
}

// createCommitFromKVLM creates a GitCommit object from a key-value list of metadata.
//
// This function extracts the necessary fields from the provided OrderedDict and constructs a GitCommit object.
// It handles the tree, parents, author, committer, and message fields.
//
// Parameters:
// - kvlm: An OrderedDict containing the key-value pairs of the commit metadata.
//
// Returns:
// - A pointer to a GitCommit struct containing the parsed commit metadata.
// - An error if any required key is missing or if there is an error during parsing.
func createCommitFromKVLM(kvlm *ordereddict.OrderedDict) (*GitCommit, error) {
	commit := GitCommit{}
	tree, exists := kvlm.Get("tree")
	if !exists {
		return nil, CommitKeyMissing("tree")
	}
	commit.Tree = string(tree.([]byte))

	parents, exists := kvlm.Get("parent")
	if exists {
		switch p := parents.(type) {
		case []byte:
			commit.Parents = []string{string(p)}
		case [][]byte:
			commit.Parents = make([]string, len(p))
			for i, parent := range p {
				commit.Parents[i] = string(parent)
			}
		}
	}

	author, exists := kvlm.Get("author")
	if !exists {
		return nil, CommitKeyMissing("author")
	}
	authorSign, err := parseSignature(author.([]byte))
	if err != nil {
		return nil, err
	}
	commit.Author = authorSign

	committer, exists := kvlm.Get("committer")
	if !exists {
		return nil, CommitKeyMissing("committer")
	}
	committerSign, err := parseSignature(committer.([]byte))
	if err != nil {
		return nil, err
	}
	commit.Committer = committerSign

	message, exists := kvlm.Get("")
	if !exists {
		return nil, CommitKeyMissing("message")
	}
	commit.Message = string(message.([]byte))

	return &commit, nil
}
