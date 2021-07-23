/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// test the main() function. As in some of the other tests in this file,
// stdout is rerouted so that the test can capture and test the output.
// It is then restored to its usual settings
func TestMainFunc(t *testing.T) {
	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "All rights reserved.") ||
		!strings.Contains(msg, "2021") {
		t.Error("Copyright notice in main() does not appear or appears incorrectly")
	}
}

func TestShowCopyright(t *testing.T) {
	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showCopyright()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "All rights reserved.") ||
		!strings.Contains(msg, "2021") {
		t.Error("Copyright does not contain expected terms")
	}
}
