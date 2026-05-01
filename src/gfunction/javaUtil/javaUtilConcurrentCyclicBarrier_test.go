/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
	"testing"
	"time"
)

func newCyclicBarrierObj() *object.Object {
	className := "java/util/concurrent/CyclicBarrier"
	obj := object.MakeEmptyObjectWithClassName(&className)
	return obj
}

func TestCyclicBarrier_Basic(t *testing.T) {
	globals.InitStringPool()
	cb := newCyclicBarrierObj()

	// Init with 2 parties
	if ret := cyclicBarrierInit([]interface{}{cb, int64(2)}); ret != nil {
		t.Fatalf("cyclicBarrierInit failed: %v", ret)
	}

	if parties := cyclicBarrierGetParties([]interface{}{cb}).(int64); parties != 2 {
		t.Fatalf("expected 2 parties, got %d", parties)
	}

	if waiting := cyclicBarrierGetNumberWaiting([]interface{}{cb}).(int64); waiting != 0 {
		t.Fatalf("expected 0 waiting, got %d", waiting)
	}

	if broken := cyclicBarrierIsBroken([]interface{}{cb}).(int64); broken != types.JavaBoolFalse {
		t.Fatalf("expected not broken")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	results := make(chan int64, 2)

	go func() {
		defer wg.Done()
		res := cyclicBarrierAwait([]interface{}{cb})
		if err, ok := res.(*ghelpers.GErrBlk); ok {
			t.Errorf("Thread 1 await failed: %v", err.ErrMsg)
			return
		}
		results <- res.(int64)
	}()

	// Wait a bit to ensure Thread 1 is waiting
	// (Not foolproof but usually works for unit tests)
	for i := 0; i < 10; i++ {
		if waiting := cyclicBarrierGetNumberWaiting([]interface{}{cb}).(int64); waiting == 1 {
			break
		}
		// small sleep or yield if needed, but for unit tests simplicity:
	}

	go func() {
		defer wg.Done()
		res := cyclicBarrierAwait([]interface{}{cb})
		if err, ok := res.(*ghelpers.GErrBlk); ok {
			t.Errorf("Thread 2 await failed: %v", err.ErrMsg)
			return
		}
		results <- res.(int64)
	}()

	wg.Wait()
	close(results)

	var resSum int64
	for r := range results {
		resSum += r
	}

	// One thread returns 1, the other 0. Sum should be 1.
	if resSum != 1 {
		t.Fatalf("Expected sum of indices to be 1, got %d", resSum)
	}

	if waiting := cyclicBarrierGetNumberWaiting([]interface{}{cb}).(int64); waiting != 0 {
		t.Fatalf("expected 0 waiting after barrier trip, got %d", waiting)
	}
}

func TestCyclicBarrier_Reset(t *testing.T) {
	globals.InitStringPool()
	cb := newCyclicBarrierObj()
	cyclicBarrierInit([]interface{}{cb, int64(2)})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		res := cyclicBarrierAwait([]interface{}{cb})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected BrokenBarrierException after reset, but got success")
		}
	}()

	// Wait for thread to be waiting
	for i := 0; i < 1000; i++ {
		if cyclicBarrierGetNumberWaiting([]interface{}{cb}).(int64) > 0 {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	cyclicBarrierReset([]interface{}{cb})
	wg.Wait()

	if broken := cyclicBarrierIsBroken([]interface{}{cb}).(int64); broken != types.JavaBoolFalse {
		t.Fatalf("expected not broken after reset")
	}
}
