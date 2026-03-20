package message

type MetadataFields struct {
	Timestamp uint32
	StreamId  uint32
	Length    uint32
	context   *Context
}

func (m *MetadataFields) Metadata() *MetadataFields {
	return m
}

func (m *MetadataFields) SetContext(c *Context) {
	m.context = c
}
