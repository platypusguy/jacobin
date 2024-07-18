/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_StringBuilder() {

	MethodSignatures["java/lang/StringBuilder.isLatin1()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/Object;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderAppendObject,
		}

}

// "java/lang/StringBuilder.append(Ljava/lang/Object;)Ljava/lang/StringBuilder;"
// Appends the string representation of the Object argument.
// The overall effect is exactly as if the argument were converted to
// a string by the method String.valueOf(Object), and the characters of that string
// were then appended to this character sequence.
func stringBuilderAppendObject(params []interface{}) interface{} {
	// params[0]: input StringBuilder Object
	inObj := params[0].(*object.Object)
	str := object.ObjectFieldToString(inObj, "value")
	errMsg := fmt.Sprintf("Not working yet. Object: %s", str)
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)

	//return nil
}
