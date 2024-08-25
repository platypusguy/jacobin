/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

// routines to dump configuration info for debugging puproses. Can be redirected to any file.
func DumpConfig(out *os.File) error {
	versionAndOs := fmt.Sprintf("Version: %s, OS: %s", GetJacobinVersion(), runtime.GOOS)
	n, err := fmt.Fprintln(out, versionAndOs)
	if err != nil {
		return errors.New(fmt.Sprintf("Error occurred %s, output %d bytes", err.Error(), n))
	} else {
		return nil
	}
}
