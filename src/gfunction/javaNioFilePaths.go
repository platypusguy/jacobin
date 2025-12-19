/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/object"
	"path/filepath"
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
	if len(params) > 1 && params[1] != nil {
		arrObj := params[1].(*object.Object)
		if arrObj != nil {
			rawArr, ok := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
			if ok {
				more = make([]string, len(rawArr))
				for i, sObj := range rawArr {
					more[i] = object.GoStringFromStringObject(sObj)
				}
			}
		}
	}

	allParts := append([]string{pathStr}, more...)
	joined := filepath.Join(allParts...)
	return newPath(joined)
}
