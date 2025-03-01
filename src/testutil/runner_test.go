package testutil

import (
	"strings"
	"testing"
)

func TestRcFromRunner(t *testing.T) {
	var rc int
	var outstr string

	rc, outstr = Runner("java", "hello", 10, false)
	if rc != RcRunnerSuccess {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerSuccess, rc, outstr)
	}

	rc, outstr = Runner("java", "hello", 0, false)
	if rc != RcRunnerTimeout {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerTimeout, rc, outstr)
	}

	rc, outstr = Runner("go", "version", 10, false)
	if rc != RcRunnerSuccess {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerSuccess, rc, outstr)
	}

	rc, outstr = Runner("go", "-NotAnOption!      hello.class", 10, false)
	if rc != RcRunnerFailure || !strings.Contains(outstr, "Usage:") {
		t.Errorf("TestRcFromRunner: expected rc=%d, observed rc=%d, outstr=%s", RcRunnerFailure, rc, outstr)
	}

}
