/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package execdata

import (
	"jacobin/globals"
	"runtime/debug"
)

// This package extracts data about the Jacobin executable
// and makes it available to the JVM

// ReadBuildInfo gets the complete set of available info
// of the currently executing Jacobin instance and load it
// into a map in the globals. The fetch of the data occurs
// only once due to the initial test.
func GetExecBuildInfo(g *globals.Globals) {
	if g.JacobinBuildData == nil {
		g.JacobinBuildData = make(map[string]string)
		info, _ := debug.ReadBuildInfo()
		for i := 0; i < len(info.Settings); i++ {
			k := info.Settings[i].Key
			v := info.Settings[i].Value
			g.JacobinBuildData[k] = v
		}
	}
}
