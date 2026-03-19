package amf3

// Array represents the AMF 3 array type (§3.11).
// AMF 3 arrays have two parts:
//   - Associative portion: string-keyed name/value pairs (like an ECMA array)
//   - Dense portion: ordered numeric-indexed values
//
// The binary representation writes:
//  1. U29A-value header (count of dense portion, low bit=1)
//  2. Name/value pairs terminated by empty string (associative portion)
//  3. Dense values in order
type Array struct {
	Associative map[string]any
	Dense       []any
}

func init() { RegisterType(new(Array)) }

func (v *Array) isObjectRefType() {}

func (v Array) Type() Type { return ArrayMarker }

func (v *Array) Read(r *Reader) error {
	ref, header, isRef, err := r.readObjectRef()
	if err != nil {
		return err
	}
	if isRef {
		if a, ok := ref.(*Array); ok {
			*v = *a
			return nil
		}
		return &UnexpectedRefTypeError{"Array", ref}
	}

	denseCount := header // remaining bits after low-bit flag = dense count

	// Add to object reference table (before reading contents, per spec)
	r.addObjectRef(v)

	// Read associative portion: name/value pairs until empty string
	v.Associative = make(map[string]any)
	for {
		key, err := r.readString()
		if err != nil {
			return err
		}
		if key == "" {
			break // empty string terminates the associative portion
		}
		val, err := r.ReadValue()
		if err != nil {
			return err
		}
		v.Associative[key] = val
	}

	// Read dense portion
	v.Dense = make([]any, denseCount)
	for i := uint32(0); i < denseCount; i++ {
		v.Dense[i], err = r.ReadValue()
		if err != nil {
			return err
		}
	}

	return nil
}

func (v Array) Write(w *Writer) error {
	denseCount := uint32(len(v.Dense))
	// U29A-value: low bit = 1, remaining bits = dense count
	err := writeU29(w.w, (denseCount<<1)|1)
	if err != nil {
		return err
	}

	// Write associative portion
	for key, val := range v.Associative {
		err = w.writeString(key)
		if err != nil {
			return err
		}
		err = w.WriteValue(val)
		if err != nil {
			return err
		}
	}
	// Terminate associative portion with empty string
	err = w.writeString("")
	if err != nil {
		return err
	}

	// Write dense portion
	for _, val := range v.Dense {
		err = w.WriteValue(val)
		if err != nil {
			return err
		}
	}

	return nil
}
