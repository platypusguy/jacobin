/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/object"
	"strings"
)

// Load_Nio_File_Paths loads MethodSignatures entries for java.nio.file.Paths
func Load_Nio_File_Paths() {

	MethodSignatures["java/nio/file/Paths.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: clinitGeneric}

	MethodSignatures["java/nio/file/Paths.<init>()V"] =
		GMeth{ParamSlots: 0, GFunction: trapProtected}

	MethodSignatures["java/nio/file/Paths.get(Ljava/lang/String;[Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 2, GFunction: pathsGet}

	MethodSignatures["java/nio/file/Paths.get(Ljava/net/URI;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
}

func pathsGet(params []interface{}) interface{} {
	pathStr := object.GoStringFromStringObject(params[0].(*object.Object))
	var more []string
	if len(params) > 1 && !object.IsNull(params[1]) {
		arrObj := params[1].(*object.Object)
		rawArr, ok := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
		if ok {
			for _, sObj := range rawArr {
				if !object.IsNull(sObj) {
					more = append(more, object.GoStringFromStringObject(sObj))
				}
			}
		}
	}

	res := pathStr
	for _, m := range more {
		if m != "" {
			if res == "" {
				res = m
			} else {
				// Join with separator, but avoid double separators if one is already present
				resHasSep := strings.HasSuffix(res, getSep())
				mHasSep := strings.HasPrefix(m, getSep())
				if resHasSep && mHasSep {
					res += m[1:]
				} else if !resHasSep && !mHasSep {
					res += getSep() + m
				} else {
					res += m
				}
			}
		}
	}
	return newPath(res)
}
