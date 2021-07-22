/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"testing"
)

func TestLogLevelTooLow(t *testing.T) {
	Global = initGlobals("test")
	err := SetLogLevel(0)
	if err == nil {
		t.Error("setting logging level to 0 did not generate an error")
	}
}

func TestLogLevelTooHigh(t *testing.T) {
	Global = initGlobals("test")
	err := SetLogLevel(99)
	if err == nil {
		t.Error("setting logging level to 99 did not generate an error")
	}
}
