/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import "errors"

// ResolveCPmethRefs resolves the method references in the constant pool of a class
func ResolveCPmethRefs(k *Klass) error {
	if k == nil || k.Data == nil || &k.Data.CP == nil {
		return errors.New("invalid class or class data in ResolveCPmethRefs")
	}

	return nil
}
