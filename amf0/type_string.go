// Code generated by "stringer -type=Type"; DO NOT EDIT.

package amf0

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NumberMarker-0]
	_ = x[BooleanMarker-1]
	_ = x[StringMarker-2]
	_ = x[ObjectMarker-3]
	_ = x[MovieclipMarker-4]
	_ = x[NullMarker-5]
	_ = x[UndefinedMarker-6]
	_ = x[ReferenceMarker-7]
	_ = x[EcmaArrayMarker-8]
	_ = x[ObjectEndMarker-9]
	_ = x[StrictArrayMarker-10]
	_ = x[DateMarker-11]
	_ = x[LongStringMarker-12]
	_ = x[UnsupportedMarker-13]
	_ = x[RecordsetMarker-14]
	_ = x[XmlDocumentMarker-15]
	_ = x[TypedObjectMarker-16]
	_ = x[AvmplusObjectMarker-17]
}

const _Type_name = "NumberMarkerBooleanMarkerStringMarkerObjectMarkerMovieclipMarkerNullMarkerUndefinedMarkerReferenceMarkerEcmaArrayMarkerObjectEndMarkerStrictArrayMarkerDateMarkerLongStringMarkerUnsupportedMarkerRecordsetMarkerXmlDocumentMarkerTypedObjectMarkerAvmplusObjectMarker"

var _Type_index = [...]uint16{0, 12, 25, 37, 49, 64, 74, 89, 104, 119, 134, 151, 161, 177, 194, 209, 226, 243, 262}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
