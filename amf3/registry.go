package amf3

// typeRegistry contains a mapping from AMF 3 type markers to prototypical instances.
var typeRegistry map[Type]MutableValue

func RegisterType(v MutableValue) {
	if typeRegistry == nil {
		typeRegistry = make(map[Type]MutableValue)
	}
	typeRegistry[v.Type()] = v
}
