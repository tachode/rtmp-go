package amf3

// AMF 3 Data Types from amf3_spec_05_05_08.pdf §3.1
type Type uint8

//go:generate stringer -type=Type
const (
	UndefinedMarker   Type = 0x00
	ByteArrayMarker   Type = 0x0C
	XmlMarker         Type = 0x0B
	ObjectMarker      Type = 0x0A
	ArrayMarker       Type = 0x09
	DateMarker        Type = 0x08
	XmlDocumentMarker Type = 0x07
	StringMarker      Type = 0x06
	DoubleMarker      Type = 0x05
	IntegerMarker     Type = 0x04
	TrueMarker        Type = 0x03
	FalseMarker       Type = 0x02
	NullMarker        Type = 0x01
)
