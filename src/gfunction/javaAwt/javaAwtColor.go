package javaAwt

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
)

func Load_Awt_Color() {

	ghelpers.MethodSignatures["java/awt/Color.<init>(FFF)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/Color.<init>(FFFF)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/Color.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/Color.<init>(IZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/Color.<init>(II)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/Color.<init>(III)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  initColorRGB,
		}

	ghelpers.MethodSignatures["java/awt/Color.<init>(Ljava/awt/color/ColorSpace;[FF)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/awt/Color.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/awt/Color.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  colorHashCode,
		}

}

// <init>(III)V: alpha defaults to 255
func initColorRGB(params []any) any {
	objectRef, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("initColorRGB: expected an object, observed: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	components, errBlk := getIntArguments("initColorRGB", params, 3)
	if errBlk != nil {
		return errBlk
	}
	return setColorValue(objectRef, components[0], components[1], components[2], 255)
}

func colorHashCode(params []any) any {
	objectRef, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("colorHashCode: expected an object, observed: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	valueField, found := objectRef.FieldTable["value"]
	if !found {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException,
			"colorHashCode: Color object has no value field")
	}

	argb, ok := valueField.Fvalue.(int64)
	if !ok {
		return ghelpers.WrongArgType("colorHashCode", "an int value field", valueField.Fvalue)
	}

	// A Java int is 32 bits signed. Any color with alpha >= 0x80 (including all opaque colors)
	// has the sign bit set, so narrow and re-widen to get the same hash as HotSpot.
	return int64(int32(argb))
}

// getIntArguments reads count int arguments starting at params[1].
func getIntArguments(funcName string, params []any, count int) ([]int64, any) {
	arguments := make([]int64, count)
	for index := 0; index < count; index++ {
		intValue, ok := params[index+1].(int64)
		if !ok {
			return nil, ghelpers.WrongArgType(funcName, "an int argument", params[index+1])
		}
		arguments[index] = intValue
	}
	return arguments, nil
}

func setColorValue(objectRef *object.Object, red, green, blue, alpha int64) any {
	badComponents := ""
	if alpha < 0 || alpha > 255 {
		badComponents += " Alpha"
	}
	if red < 0 || red > 255 {
		badComponents += " Red"
	}
	if green < 0 || green > 255 {
		badComponents += " Green"
	}
	if blue < 0 || blue > 255 {
		badComponents += " Blue"
	}
	if badComponents != "" {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"Color parameter outside of expected range:"+badComponents)
	}
	storeColorValue(objectRef, (alpha&0xFF)<<24|(red&0xFF)<<16|(green&0xFF)<<8|(blue&0xFF))
	return nil
}

// storeColorValue keeps the field's existing metadata and stores the packed ARGB as a
// sign-extended Java int, so 0xFF000002 reads back as -16777214, as in HotSpot.
func storeColorValue(objectRef *object.Object, packed int64) {
	valueField := objectRef.FieldTable["value"]
	valueField.Fvalue = int64(int32(packed))
	objectRef.FieldTable["value"] = valueField
}
