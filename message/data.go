package message

// Data is the common interface for AMF0 and AMF3 data messages.
type Data interface {
	Message
	GetHandler() string
	GetParameters() []any
}
