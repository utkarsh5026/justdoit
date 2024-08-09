package objects

// BlobObject represents a Git blob object in the Git object model.
// A blob object is used to store the contents of a file.
//
// Fields:
// - format: The type of the Git object, which is BlobType for blob objects.
// - data: The raw byte data of the blob, which contains the file contents.
type BlobObject struct {
	data []byte
}

func Blob() *BlobObject {
	return &BlobObject{}
}

func (b *BlobObject) Serialize() ([]byte, error) {
	return b.data, nil
}

func (b *BlobObject) Deserialize(data []byte) error {
	b.data = data
	return nil
}

func (b *BlobObject) Format() GitObjectType {
	return BlobType
}

func (b *BlobObject) SetData(data []byte) {
	b.data = data
}
