package message

type Message interface {
	Type() Type
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
	Metadata() *MetadataFields
}
