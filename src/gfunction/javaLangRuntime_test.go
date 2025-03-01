/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import "testing"

func TestMaxMemory(t *testing.T) {
	mem := maxMemory(nil)
	if mem.(int64) <= 1_000_000 {
		t.Errorf("maxMemory() = %d bytes; expected > 1,000,000", mem)
	}
}

func TestTotalMemory(t *testing.T) {
	mem := totalMemory(nil)
	if mem.(int64) <= 1_000_000 {
		t.Errorf("totalMemory() = %d bytes; expected > 1,000,000", mem)
	}
}

func TestRuntimeCPUs(t *testing.T) {
	cpus := runtimeAvailableProcessors(nil)
	if cpus.(int64) <= 1 {
		t.Errorf("runtimeCPUs() = %d; expected > 1", cpus)
	}
}
