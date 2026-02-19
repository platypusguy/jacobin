/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"sync"
	"testing"
)

func TestObjLock_HighContentionStress(t *testing.T) {
	obj := MakeEmptyObject()
	const numThreads = 64
	const iterations = 2000
	var wg sync.WaitGroup
	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		tid := int32(i)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if err := obj.ObjLock(tid); err != nil {
					t.Errorf("Thread %d failed to lock at iteration %d: %v", tid, j, err)
					return
				}
				// Simulate some work
				if j%10 == 0 {
					// occasional recursive lock
					if err := obj.ObjLock(tid); err == nil {
						_ = obj.ObjUnlock(tid)
					}
				}
				if err := obj.ObjUnlock(tid); err != nil {
					t.Errorf("Thread %d failed to unlock at iteration %d: %v", tid, j, err)
					return
				}
			}
		}()
	}

	wg.Wait()
}
