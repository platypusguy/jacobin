package testutil

import (
	"strings"
	"testing"
)

func TestRcFromRunner(t *testing.T) {
	var rc int
	var outstr string

	rc, outstr = Runner("jacobin", "hello.class", 10, false)
	if rc != RcRunnerSuccess {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerSuccess, rc, outstr)
	}

	rc, outstr = Runner("jacobin", "hello.class", 0, false)
	if rc != RcRunnerTimeout {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerTimeout, rc, outstr)
	}

	rc, outstr = Runner("jacobin", "-NotAnOption!      hello.class", 10, false)
	if rc != RcRunnerFailure || !strings.Contains(outstr, "Usage: jacobin") {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerFailure, rc, outstr)
	}

	rc, outstr = Runner("go", "version", 10, false)
	if rc != RcRunnerSuccess {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerSuccess, rc, outstr)
	}

}
