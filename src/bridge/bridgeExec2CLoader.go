/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package bridge

import "jacobin/classloader"

func LoadClassFromName(name string) error {
	return classloader.LoadClassFromNameOnly(name)
}
