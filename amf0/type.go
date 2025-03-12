package amf0

// AMF 0 Data Types from amf0_spec_121207.pdf §3
type Type uint8

//go:generate stringer -type=Type
const (
	NumberMarker        Type = 0x00
	BooleanMarker       Type = 0x01
	StringMarker        Type = 0x02
	ObjectMarker        Type = 0x03
	MovieclipMarker     Type = 0x04 // reserved, not supported
	NullMarker          Type = 0x05
	UndefinedMarker     Type = 0x06
	ReferenceMarker     Type = 0x07
	EcmaArrayMarker     Type = 0x08
	ObjectEndMarker     Type = 0x09
	StrictArrayMarker   Type = 0x0A
	DateMarker          Type = 0x0B
	LongStringMarker    Type = 0x0C
	UnsupportedMarker   Type = 0x0D
	RecordsetMarker     Type = 0x0E // reserved, not supported
	XmlDocumentMarker   Type = 0x0F
	TypedObjectMarker   Type = 0x10
	AvmplusObjectMarker Type = 0x11
)
