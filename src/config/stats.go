/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import (
	"fmt"
	"os"
	"runtime"
)

// routines to dump configuration info for debugging puproses
func DumpConfig(out *os.File) {
	fmt.Fprintf(out, "Version: %s, OS %s\n", GetJacobinVersion(), runtime.GOOS)
}
