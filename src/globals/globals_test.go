/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"os"
	"testing"
)

func TestJavaHomeFormat(t *testing.T) {
	origJavaHome := os.Getenv("JAVA_HOME")
	os.Setenv("JAVA_HOME", "foo/bar")
	InitJavaHome()
	ret := JavaHome()
	if ret != "foo\\bar\\" {
		t.Errorf("Expecting a JAVA_HOME of 'foo\\bar\\', got: %s", ret)
	}
	os.Setenv("JAVA_HOME", origJavaHome)
}

func TestJacobinHomeFormat(t *testing.T) {
	origJavaHome := os.Getenv("JAVA_HOME")
	os.Setenv("JACOBIN_HOME", "foo/bar")
	InitJacobinHome()
	ret := GetGlobalRef().JacobinHome
	if ret != "foo\\bar\\" {
		t.Errorf("Expecting a JACOBIN_HOME of 'foo\\bar\\', got: %s", ret)
	}
	os.Setenv("JACOBIN_HOME", origJavaHome)
}
