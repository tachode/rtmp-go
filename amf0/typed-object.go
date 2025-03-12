package amf0

import "io"

type TypedObject struct {
	ClassName string
	Object    Object
}

func init() { RegisterType(new(TypedObject)) }

func (v TypedObject) Type() Type { return TypedObjectMarker }

func (v *TypedObject) Read(r io.Reader) error {
	var err error
	v.ClassName, err = readString[uint16](r)
	if err != nil {
		return err
	}
	return v.Object.Read(r)
}

func (v TypedObject) Write(w io.Writer) error {
	err := writeString[uint16](w, v.ClassName)
	if err != nil {
		return err
	}
	return v.Object.Write(w)
}
