package cmd

import "github.com/utkarsh5026/justdoit/app/ordereddict"

// BlobObject represents a Git blob object in the Git object model.
// A blob object is used to store the contents of a file.
//
// Fields:
// - format: The type of the Git object, which is BlobType for blob objects.
// - data: The raw byte data of the blob, which contains the file contents.
type BlobObject struct {
	format GitObjectType
	data   []byte
}

func Blob() *BlobObject {
	return &BlobObject{
		format: BlobType,
	}
}

func (b *BlobObject) Serialize() ([]byte, error) {
	return b.data, nil
}

func (b *BlobObject) Deserialize(data []byte) error {
	b.data = data
	return nil
}

func (b *BlobObject) Format() GitObjectType {
	return b.format
}

func (b *BlobObject) SetData(data []byte) {
	b.data = data
}

// CommitObject represents a Git commit object in the Git object model.
// A commit object contains metadata about the commit, such as the author, committer, and commit message.
// All the metadata is stored in a key-value list.
//
// Fields:
// - format: The type of the Git object, which is CommitType for commit objects.
// - kvlm: An OrderedDict that stores key-value pairs of the commit metadata
type CommitObject struct {
	format GitObjectType
	kvlm   *ordereddict.OrderedDict
}

func Commit() *CommitObject {
	return &CommitObject{
		format: CommitType,
		kvlm:   ordereddict.New(),
	}
}

func (c *CommitObject) Serialize() ([]byte, error) {
	return KvlmSerialize(c.kvlm)
}

func (c *CommitObject) Deserialize(data []byte) error {
	c.kvlm = KvlmParse(data, 0, c.kvlm)
	return nil
}

func (c *CommitObject) Format() GitObjectType {
	return c.format
}

func (c *CommitObject) SetData(data []byte) {
	c.kvlm = KvlmParse(data, 0, nil)
}
