package amf3

import (
	"fmt"
	"io"
)

// Object represents the AMF 3 object type (§3.12).
// Objects have trait information describing their class, and can be:
//   - Anonymous: no class name
//   - Typed: has a registered class name
//   - Dynamic: can have extra name/value pairs beyond sealed members
//   - Externalizable: completely custom serialization
type Object struct {
	Traits         *TraitInfo
	SealedMembers  map[string]any // sealed member name -> value
	DynamicMembers map[string]any // dynamic member name -> value (only if traits.IsDynamic)
	External       []byte         // raw bytes for externalizable objects
}

func init() { RegisterType(new(Object)) }

func (v Object) Get(key string) (obj any, found bool) {
	obj, found = v.SealedMembers[key]
	if !found {
		obj, found = v.DynamicMembers[key]
	}
	return
}

func (v *Object) isObjectRefType() {}

func (v Object) Type() Type { return ObjectMarker }

func (v *Object) Read(r *Reader) error {
	ref, header, isRef, err := r.readObjectRef()
	if err != nil {
		return err
	}
	if isRef {
		if o, ok := ref.(*Object); ok {
			*v = *o
			return nil
		}
		return &UnexpectedRefTypeError{"Object", ref}
	}

	// Add to object reference table before reading contents
	r.addObjectRef(v)

	// Parse trait information from header bits
	// header has already been shifted right by 1 (object ref bit removed)
	v.Traits, err = r.readTraits(header)
	if err != nil {
		return err
	}

	if v.Traits.IsExternalizable {
		// For externalizable objects, the remaining bytes are opaque.
		// We read until EOF or until the caller knows how many bytes to expect.
		// In practice, the caller should know the format. We store raw bytes.
		// Since we can't know the length, we don't read here — we leave it
		// to the caller to handle via the External field.
		// For now, we signal this condition:
		return fmt.Errorf("AMF3 externalizable object '%s': custom deserialization not supported", v.Traits.ClassName)
	}

	// Read sealed member values in the same order as trait member names
	v.SealedMembers = make(map[string]any, len(v.Traits.Members))
	for _, name := range v.Traits.Members {
		val, err := r.ReadValue()
		if err != nil {
			return err
		}
		v.SealedMembers[name] = val
	}

	// Read dynamic members if the trait is dynamic
	if v.Traits.IsDynamic {
		v.DynamicMembers = make(map[string]any)
		for {
			key, err := r.readString()
			if err != nil {
				return err
			}
			if key == "" {
				break
			}
			val, err := r.ReadValue()
			if err != nil {
				return err
			}
			v.DynamicMembers[key] = val
		}
	}

	return nil
}

// readTraits reads trait information from the header bits.
// The header has already been shifted right by 1 (readObjectRef consumed bit 0).
//
// Original U29 bit layout:
//
//	bit 0: object-ref flag (already consumed, was 1)
//	bit 1: trait-ref flag (0 = trait reference, 1 = inline traits)
//	bit 2: externalizable flag (only if inline traits)
//	bit 3: dynamic flag (only if inline, non-externalizable)
//	bits 4+: sealed member count (only if inline, non-externalizable)
//
// So the header we receive has:
//
//	bit 0 = original bit 1 (trait-ref flag)
//	bit 1 = original bit 2 (ext flag)
//	bit 2 = original bit 3 (dynamic flag)
//	bits 3+ = original bits 4+ (member count)
func (r *Reader) readTraits(header uint32) (*TraitInfo, error) {
	// Bit 0: trait reference flag
	if header&1 == 0 {
		index := header >> 1
		return r.getTraitRef(index)
	}
	header >>= 1 // consume trait-ref flag

	trait := &TraitInfo{}

	// Bit 0 (of shifted header): externalizable flag
	if header&1 == 1 {
		trait.IsExternalizable = true
		// Read class name
		var err error
		trait.ClassName, err = r.readString()
		if err != nil {
			return nil, err
		}
		r.addTraitRef(trait)
		return trait, nil
	}
	header >>= 1 // consume ext flag (0)

	// Bit 0 (of shifted header): dynamic flag
	trait.IsDynamic = header&1 == 1
	header >>= 1 // consume dynamic flag

	// Remaining bits: sealed member count
	memberCount := header

	// Read class name
	var err error
	trait.ClassName, err = r.readString()
	if err != nil {
		return nil, err
	}

	// Read sealed member names
	trait.Members = make([]string, memberCount)
	for i := uint32(0); i < memberCount; i++ {
		trait.Members[i], err = r.readString()
		if err != nil {
			return nil, err
		}
	}

	r.addTraitRef(trait)
	return trait, nil
}

func (v Object) Write(w *Writer) error {
	if v.Traits != nil && v.Traits.IsExternalizable {
		return v.writeExternalizable(w)
	}

	return v.writeNormal(w)
}

func (v Object) writeNormal(w *Writer) error {
	trait := v.Traits
	if trait == nil {
		// Default to anonymous dynamic object
		trait = &TraitInfo{IsDynamic: true}
	}

	// Try trait reference
	traitKey := trait.key()
	if idx, ok := w.writeTraitRef(traitKey); ok {
		// U29O-traits-ref: ...xxxx10 (bit 0=1 instance, bit 1=0 trait ref)
		// remaining bits = trait ref index
		header := uint32(idx<<2) | 0x01
		err := writeU29(w.w, header)
		if err != nil {
			return err
		}
	} else {
		w.addTraitRef(traitKey)
		// U29O-traits: ...xxxx0110 if not dynamic, ...xxxx1110 if dynamic
		// Bits: [member count] [dynamic] [not-ext=0] [inline-traits=1] [instance=1]
		memberCount := uint32(len(trait.Members))
		header := (memberCount << 4)
		if trait.IsDynamic {
			header |= 0x08 // dynamic flag
		}
		header |= 0x03 // bits: inline-traits=1, instance=1

		err := writeU29(w.w, header)
		if err != nil {
			return err
		}

		// Write class name
		err = w.writeString(trait.ClassName)
		if err != nil {
			return err
		}

		// Write sealed member names
		for _, name := range trait.Members {
			err = w.writeString(name)
			if err != nil {
				return err
			}
		}
	}

	// Write sealed member values in order
	for _, name := range trait.Members {
		val := v.SealedMembers[name]
		err := w.WriteValue(val)
		if err != nil {
			return err
		}
	}

	// Write dynamic members if dynamic
	if trait.IsDynamic {
		for key, val := range v.DynamicMembers {
			err := w.writeString(key)
			if err != nil {
				return err
			}
			err = w.WriteValue(val)
			if err != nil {
				return err
			}
		}
		// Terminate with empty string
		err := w.writeString("")
		if err != nil {
			return err
		}
	}

	return nil
}

func (v Object) writeExternalizable(w *Writer) error {
	trait := v.Traits

	// U29O-traits-ext: bits ...0111
	// [0 bits for member count] [ext=1] [inline-traits=1] [instance=1]
	header := uint32(0x07)
	err := writeU29(w.w, header)
	if err != nil {
		return err
	}

	err = w.writeString(trait.ClassName)
	if err != nil {
		return err
	}

	// Write raw external data
	_, err = io.Copy(w.w, newByteReader(v.External))
	return err
}

// key returns a canonical string representation of the trait for reference deduplication.
func (t *TraitInfo) key() string {
	s := t.ClassName
	if t.IsDynamic {
		s += "|D"
	}
	if t.IsExternalizable {
		s += "|E"
	}
	for _, m := range t.Members {
		s += "|" + m
	}
	return s
}

// byteReader wraps a []byte to implement io.Reader.
type byteReader struct {
	data []byte
	pos  int
}

func newByteReader(data []byte) *byteReader {
	return &byteReader{data: data}
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
