/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/src/object"
	"jacobin/src/types"
)

// A java/lang/Class object (or JLC object) is a standard object with the
// following fields:
//   name: *object.Object, which is the name of the class as a Java string object
//      note: the object's KlassName is always java/lang/Class
//   $klass: Klass, points to the ClData metadata in the method area for the named class
//   $statics: []string, the names of the static fields of the class. These
//      can be used as keys into the JVM's statics table.'

// Makes a JLC object for a class and fills in only the name
func MakeJlcObject(classname string) *object.Object {
	o := object.MakeEmptyObject()
	o.KlassName = types.StringPoolJavaLangClassIndex
	o.FieldTable["name"] = object.Field{Ftype: types.Ref,
		Fvalue: object.StringObjectFromGoString(classname)}
	o.FieldTable["$klass"] = object.Field{Ftype: types.RawGoPointer,
		Fvalue: nil} // points to the Klass object in metadata
	o.FieldTable["$statics"] = object.Field{Ftype: types.Array,
		Fvalue: []string{}} // array of static field names for this class
	return o
}
