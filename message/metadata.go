package message

type MetadataFields struct {
	Timestamp uint32
	StreamId  uint32
	Length    uint32
}

func (m *MetadataFields) Metadata() *MetadataFields {
	return m
}
