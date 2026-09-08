/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/frames"
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
		res := cyclicBarrierAwait([]interface{}{cb, object.Null})
		if err, ok := res.(*ghelpers.GErrBlk); ok {
			t.Errorf("Thread 1 await failed: %v", err.ErrMsg)
			return
		}
		results <- res.(int64)
	}()

	for i := 0; i < 100; i++ {
		if waiting := cyclicBarrierGetNumberWaiting([]interface{}{cb}).(int64); waiting == 1 {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	go func() {
		defer wg.Done()
		res := cyclicBarrierAwait([]interface{}{cb, object.Null})
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

func TestCyclicBarrier_MultipleCycles(t *testing.T) {
	globals.InitStringPool()
	cb := newCyclicBarrierObj()
	cyclicBarrierInit([]interface{}{cb, int64(3)})

	for cycle := 0; cycle < 3; cycle++ {
		var wg sync.WaitGroup
		wg.Add(3)
		results := make(chan int64, 3)

		for i := 0; i < 3; i++ {
			go func() {
				defer wg.Done()
				res := cyclicBarrierAwait([]interface{}{cb})
				if err, ok := res.(*ghelpers.GErrBlk); ok {
					t.Errorf("Await failed in cycle %d: %v", cycle, err.ErrMsg)
					return
				}
				results <- res.(int64)
			}()
		}

		wg.Wait()
		close(results)

		var sum int64
		for r := range results {
			sum += r
		}
		// Indices should be 2 + 1 + 0 = 3
		if sum != 3 {
			t.Fatalf("Cycle %d expected index sum 3, got %d", cycle, sum)
		}
		if waiting := cyclicBarrierGetNumberWaiting([]interface{}{cb}).(int64); waiting != 0 {
			t.Fatalf("Cycle %d expected 0 waiting, got %d", cycle, waiting)
		}
	}
}

func TestCyclicBarrier_Interrupt(t *testing.T) {
	globals.InitStringPool()
	cb := newCyclicBarrierObj()
	cyclicBarrierInit([]interface{}{cb, int64(2)})

	th := object.MakeEmptyObject()
	th.FieldTable["interrupted"] = object.Field{Ftype: types.Int, Fvalue: types.JavaBoolTrue}

	// Should fail immediately if already interrupted
	res := cyclicBarrierAwait([]interface{}{cb, th})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.InterruptedException {
		t.Fatalf("Expected InterruptedException, got %v", res)
	}

	if broken := cyclicBarrierIsBroken([]interface{}{cb}).(int64); broken != types.JavaBoolTrue {
		t.Fatalf("expected broken barrier after interrupt, got %v", broken)
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
		res := cyclicBarrierAwait([]interface{}{cb, object.Null})
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

func TestCyclicBarrier_NeedsContext_Await(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	th := object.MakeEmptyObject()
	th.FieldTable["interrupted"] = object.Field{Ftype: types.Int, Fvalue: types.JavaBoolFalse}

	thID := 42
	gr := globals.GetGlobalRef()
	gr.ThreadLock.Lock()
	gr.Threads[thID] = th
	gr.ThreadLock.Unlock()

	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(0)
	f.Thread = thID
	_ = frames.PushFrame(fs, f)

	cb := newCyclicBarrierObj()
	cyclicBarrierInit([]interface{}{cb, int64(1)})

	// Await with [fs, cb]
	res := cyclicBarrierAwait([]interface{}{fs, cb})
	if idx, ok := res.(int64); !ok || idx != 0 {
		t.Fatalf("expected return 0 for single party barrier, got %v (%T)", res, res)
	}
}

func TestCyclicBarrier_MethodSignatures(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Util_Concurrent_CyclicBarrier()

	awaitSig := "java/util/concurrent/CyclicBarrier.await()I"
	gm, ok := ghelpers.MethodSignatures[awaitSig]
	if !ok {
		t.Fatalf("missing signature %s", awaitSig)
	}
	if !gm.NeedsContext {
		t.Fatalf("expected %s NeedsContext to be true", awaitSig)
	}

	initSig := "java/util/concurrent/CyclicBarrier.<init>(ILjava/lang/Runnable;)V"
	if gmInit, ok := ghelpers.MethodSignatures[initSig]; !ok || gmInit.ParamSlots != 2 {
		t.Fatalf("missing or invalid %s", initSig)
	}
}

func TestCyclicBarrier_DefensiveChecks(t *testing.T) {
	globals.InitStringPool()

	// Null / invalid params to Init
	if err := cyclicBarrierInit(nil); err == nil {
		t.Fatalf("expected error on nil params to init")
	}
	if err := cyclicBarrierInit([]interface{}{object.Null, int64(2)}); err == nil {
		t.Fatalf("expected error on null obj to init")
	}
	cb := newCyclicBarrierObj()
	if err := cyclicBarrierInit([]interface{}{cb, int64(-1)}); err == nil {
		t.Fatalf("expected error on negative parties to init")
	}

	// Uninitialized await / parties / isBroken / reset
	uninitCb := newCyclicBarrierObj()
	if err, ok := cyclicBarrierAwait([]interface{}{uninitCb}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NPE on uninitialized await, got %v", err)
	}
	if err, ok := cyclicBarrierGetParties([]interface{}{uninitCb}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NPE on uninitialized getParties, got %v", err)
	}
	if err, ok := cyclicBarrierIsBroken([]interface{}{uninitCb}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NPE on uninitialized isBroken, got %v", err)
	}
	if err, ok := cyclicBarrierReset([]interface{}{uninitCb}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NPE on uninitialized reset, got %v", err)
	}
	if err, ok := cyclicBarrierGetNumberWaiting([]interface{}{uninitCb}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NPE on uninitialized getNumberWaiting, got %v", err)
	}
}
