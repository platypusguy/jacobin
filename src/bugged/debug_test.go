package bugged

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"testing"
)

func TestTVO(t *testing.T) {
	globals.InitGlobals("test")
	obj := object.MakeEmptyObject()
	_ = TVO(obj)

	clName1 := "apple/beet/carrot"
	clName2 := "dandelion/daisy/bluebell"
	obj1 := object.MakeEmptyObjectWithClassName(&clName1)
	obj2 := object.MakeEmptyObjectWithClassName(&clName2)
	obj1.FieldTable["field6b"] = object.Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}
	obj1.FieldTable["field7"] = object.Field{Ftype: types.ByteArray, Fvalue: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	obj1.FieldTable["field8"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{61, 62, 63, 64, 65, 66, 67, 68, 69, 70}}
	obj1.FieldTable["field9"] = object.Field{Ftype: types.Int, Fvalue: uint32(math.Pow(2, 27) - 1)}
	obj1.FieldTable["field10"] = object.Field{Ftype: types.Int, Fvalue: uint16(32767)}
	obj1.FieldTable["field1"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString("value1")}
	obj1.FieldTable["field2"] = object.Field{Ftype: types.Int, Fvalue: 42}
	obj1.FieldTable["field3"] = object.Field{Ftype: types.Double, Fvalue: 3.14}
	obj1.FieldTable["field4"] = object.Field{Ftype: types.StringClassName, Fvalue: object.JavaByteArrayFromGoString("value4")}
	obj1.FieldTable["field5"] = object.Field{Ftype: types.Ref, Fvalue: obj2}
	obj1.FieldTable["field6a"] = object.Field{Ftype: types.Bool, Fvalue: types.JavaBoolTrue}
	obj3 := object.StringObjectFromGoString("Hey diddle diddle .....")
	obj1.FieldTable["field11"] = object.Field{Ftype: types.StringClassName, Fvalue: obj3}

	t.Log(TVO(obj1))
}
