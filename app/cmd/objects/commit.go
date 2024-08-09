package objects

import "github.com/utkarsh5026/justdoit/app/ordereddict"

// CommitObject represents a Git commit object in the Git object model.
// A commit object contains metadata about the commit, such as the author, committer, and commit message.
// All the metadata is stored in a key-value list.
//
// Fields:
// - format: The type of the Git object, which is CommitType for commit objects.
// - kvlm: An OrderedDict that stores key-value pairs of the commit metadata
type CommitObject struct {
	kvlm *ordereddict.OrderedDict
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
	return nil
}

func (c *CommitObject) Format() GitObjectType {
	return CommitType
}

func (c *CommitObject) SetData(data []byte) {
	c.kvlm = KvlmParse(data, 0, nil)
}
