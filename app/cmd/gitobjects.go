package cmd

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
