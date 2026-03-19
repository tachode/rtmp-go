package amf3

import "fmt"

// UnexpectedRefTypeError is returned when an object reference resolves to
// an unexpected type.
type UnexpectedRefTypeError struct {
	Expected string
	Got      any
}

func (e *UnexpectedRefTypeError) Error() string {
	return fmt.Sprintf("AMF3 object reference: expected %s, got %T", e.Expected, e.Got)
}
