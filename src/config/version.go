/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import "strconv"

// v 0.6.000 begun 31 Aug 24 at 3,012 GitHub commits
// v 0.6.100 made   5 Nov 24 at 3,280 GitHub commits - last version with the original intepreter
// v 0.6.200 made   6 Nov 24 - first version with the original interpreter removed

var JacobinVersion = "0.6.201"

// GetJacobinVersion returns a manually updated version number and an
// automatically updated build #. The latter being updated by bumpbuildno.go
// on the build cycles of @alb's build system.
func GetJacobinVersion() string {
	return JacobinVersion + " Build " + strconv.Itoa(BuildNo)
}
