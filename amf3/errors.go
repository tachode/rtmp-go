package amf3

import (
	"errors"
	"fmt"
)

// ErrNilWriter is returned when WriteValue is called on a Writer whose
// underlying io.Writer is nil (e.g., after being released via SetWriter(nil)).
var ErrNilWriter = errors.New("amf3: Writer has no underlying io.Writer (set one with SetWriter)")

// ErrNilReader is returned when ReadValue is called on a Reader whose
// underlying io.Reader is nil (e.g., after being released via SetReader(nil)).
var ErrNilReader = errors.New("amf3: Reader has no underlying io.Reader (set one with SetReader)")

// UnexpectedRefTypeError is returned when an object reference resolves to
// an unexpected type.
type UnexpectedRefTypeError struct {
	Expected string
	Got      any
}

func (e *UnexpectedRefTypeError) Error() string {
	return fmt.Sprintf("AMF3 object reference: expected %s, got %T", e.Expected, e.Got)
}
