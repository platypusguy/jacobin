/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import "strconv"

// v 0.6.0 begun 8/31/24 at 3,012 GitHub commits

var JacobinVersion = "0.6.017"

// GetJacobinVersion returns a manually updated version number and an
// automatically updated build #. The latter being updated by bumpbuildno.go
// on the build cycles of @alb's build system.
func GetJacobinVersion() string {
	return JacobinVersion + " Build " + strconv.Itoa(BuildNo)
}
