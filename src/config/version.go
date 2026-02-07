/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import "strconv"

// v 0.6.000 begun 31 Aug 24 at 3,012 GitHub commits (end of year 3 of the project)
// v 0.6.100 made   5 Nov 24 at 3,280 GitHub commits - last version with the original intepreter
// v 0.6.200 made   6 Nov 24 - first version with the new interpreter
// v 0.7.000 made  28 Feb 25 at 3,733 GitHub commits (technically, this file updated on 2 Mar 24)
// v 0.8.000 made  31 Aug 25 at 4,245 GitHub commits (end of year 4 of the project)

var JacobinVersion = "0.8.113"

// GetJacobinVersion returns a manually updated version number and an
// automatically updated build #. The latter being updated by bumpbuildno.go
// on the build cycles of @alb's build system.
func GetJacobinVersion() string {
	return JacobinVersion + " Build " + strconv.Itoa(BuildNo)
}
