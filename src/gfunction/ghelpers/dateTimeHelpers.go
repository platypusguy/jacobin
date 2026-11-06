package ghelpers

import (
	"jacobin/src/object"
	"jacobin/src/types"
)

// CreateTemporalUnit allocates a new ChronoUnit/TemporalUnit object with the given name.
func CreateTemporalUnit(name string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameChronoUnit)
	nameObj := object.StringObjectFromGoString(name)
	obj.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: nameObj}
	return obj
}
