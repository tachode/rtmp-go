package amf3

// TraitInfo describes the characteristics of an Object's class (§3.12).
type TraitInfo struct {
	ClassName        string
	IsDynamic        bool
	IsExternalizable bool
	Members          []string // sealed member names
}
