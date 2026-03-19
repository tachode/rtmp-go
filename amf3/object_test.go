package amf3_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestObject_ReadWrite_AnonymousDynamic(t *testing.T) {
	// Anonymous dynamic object with one dynamic member "key" -> "value"
	//
	// U29O-traits: bits = [0 members][dynamic=1][ext=0][inline=1][instance=1]
	// = (0 << 4) | 0x08 | 0x03 = 0x0B
	var data bytes.Buffer
	data.WriteByte(byte(amf3.ObjectMarker))
	data.Write(u29Encode(0x0B)) // inline, dynamic, 0 sealed members
	data.Write(utf8vrEmpty())   // class name = "" (anonymous)
	data.Write(utf8vr("key"))   // dynamic member name
	data.WriteByte(byte(amf3.StringMarker))
	data.Write(utf8vr("value")) // dynamic member value
	data.Write(utf8vrEmpty())   // end dynamic members

	r := amf3.NewReader(&data)
	val, err := r.ReadValue()
	assert.NoError(t, err)

	obj, ok := val.(*amf3.Object)
	assert.True(t, ok)
	assert.Equal(t, "", obj.Traits.ClassName)
	assert.True(t, obj.Traits.IsDynamic)
	assert.False(t, obj.Traits.IsExternalizable)
	assert.Equal(t, 0, len(obj.Traits.Members))
	assert.Equal(t, amf3.String("value"), obj.DynamicMembers["key"])
}

func TestObject_ReadWrite_TypedSealed(t *testing.T) {
	// Typed, non-dynamic object "MyClass" with 2 sealed members: "name", "age"
	//
	// U29O-traits: (2 << 4) | 0x00 | 0x03 = 0x23
	// = 2 members, not dynamic, not ext, inline, instance
	var data bytes.Buffer
	data.WriteByte(byte(amf3.ObjectMarker))
	data.Write(u29Encode(0x23))   // inline, non-dynamic, 2 sealed members
	data.Write(utf8vr("MyClass")) // class name
	data.Write(utf8vr("name"))    // member name 1
	data.Write(utf8vr("age"))     // member name 2
	// Member values (in same order as names)
	data.WriteByte(byte(amf3.StringMarker))
	data.Write(utf8vr("Alice"))
	data.WriteByte(byte(amf3.IntegerMarker))
	data.Write(u29Encode(30))

	r := amf3.NewReader(&data)
	val, err := r.ReadValue()
	assert.NoError(t, err)

	obj, ok := val.(*amf3.Object)
	assert.True(t, ok)
	assert.Equal(t, "MyClass", obj.Traits.ClassName)
	assert.False(t, obj.Traits.IsDynamic)
	assert.Equal(t, []string{"name", "age"}, obj.Traits.Members)
	assert.Equal(t, amf3.String("Alice"), obj.SealedMembers["name"])
	assert.Equal(t, amf3.Integer(30), obj.SealedMembers["age"])
}

func TestObject_RoundTrip(t *testing.T) {
	t.Run("anonymous dynamic", func(t *testing.T) {
		obj := &amf3.Object{
			Traits: &amf3.TraitInfo{
				IsDynamic: true,
			},
			SealedMembers:  map[string]any{},
			DynamicMembers: map[string]any{"foo": amf3.String("bar")},
		}

		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(obj)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		readObj, ok := val.(*amf3.Object)
		assert.True(t, ok)
		assert.True(t, readObj.Traits.IsDynamic)
		assert.Equal(t, amf3.String("bar"), readObj.DynamicMembers["foo"])
	})

	t.Run("typed sealed", func(t *testing.T) {
		obj := &amf3.Object{
			Traits: &amf3.TraitInfo{
				ClassName: "Person",
				Members:   []string{"name", "age"},
			},
			SealedMembers: map[string]any{
				"name": amf3.String("Bob"),
				"age":  amf3.Integer(25),
			},
		}

		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(obj)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		readObj, ok := val.(*amf3.Object)
		assert.True(t, ok)
		assert.Equal(t, "Person", readObj.Traits.ClassName)
		assert.Equal(t, []string{"name", "age"}, readObj.Traits.Members)
		assert.Equal(t, amf3.String("Bob"), readObj.SealedMembers["name"])
		assert.Equal(t, amf3.Integer(25), readObj.SealedMembers["age"])
	})
}

func TestObject_TraitReference(t *testing.T) {
	// Write two objects with the same traits - second should use trait reference
	obj1 := &amf3.Object{
		Traits: &amf3.TraitInfo{
			ClassName: "Point",
			Members:   []string{"x", "y"},
		},
		SealedMembers: map[string]any{
			"x": amf3.Integer(10),
			"y": amf3.Integer(20),
		},
	}

	obj2 := &amf3.Object{
		Traits: &amf3.TraitInfo{
			ClassName: "Point",
			Members:   []string{"x", "y"},
		},
		SealedMembers: map[string]any{
			"x": amf3.Integer(30),
			"y": amf3.Integer(40),
		},
	}

	var buf bytes.Buffer
	w := amf3.NewWriter(&buf)
	err := w.WriteValue(obj1)
	assert.NoError(t, err)
	err = w.WriteValue(obj2)
	assert.NoError(t, err)

	// Read both back
	r := amf3.NewReader(&buf)
	val1, err := r.ReadValue()
	assert.NoError(t, err)
	readObj1, ok := val1.(*amf3.Object)
	assert.True(t, ok)
	assert.Equal(t, "Point", readObj1.Traits.ClassName)
	assert.Equal(t, amf3.Integer(10), readObj1.SealedMembers["x"])

	val2, err := r.ReadValue()
	assert.NoError(t, err)
	readObj2, ok := val2.(*amf3.Object)
	assert.True(t, ok)
	assert.Equal(t, "Point", readObj2.Traits.ClassName)
	assert.Equal(t, amf3.Integer(30), readObj2.SealedMembers["x"])
	assert.Equal(t, amf3.Integer(40), readObj2.SealedMembers["y"])
}
