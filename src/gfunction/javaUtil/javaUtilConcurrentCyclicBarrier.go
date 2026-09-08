/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"container/list"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
)

type generation struct {
	broken bool
}

type cyclicBarrierState struct {
	mu            sync.RWMutex
	parties       int64
	count         int64
	gen           *generation
	barrierCond   *sync.Cond
	barrierAction *object.Object
}

func (s *cyclicBarrierState) nextGeneration() {
	s.barrierCond.Broadcast()
	s.count = s.parties
	s.gen = &generation{broken: false}
}

func (s *cyclicBarrierState) breakBarrier() {
	s.gen.broken = true
	s.count = s.parties
	s.barrierCond.Broadcast()
}

func Load_Util_Concurrent_CyclicBarrier() {
	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cyclicBarrierInit,
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.<init>(ILjava/lang/Runnable;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  cyclicBarrierInitAction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.await()I"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			NeedsContext: true,
			GFunction:    cyclicBarrierAwait,
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.await(JLjava/util/concurrent/TimeUnit;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction, // Timeout await not yet implemented
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.getParties()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cyclicBarrierGetParties,
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.isBroken()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cyclicBarrierIsBroken,
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cyclicBarrierReset,
		}

	ghelpers.MethodSignatures["java/util/concurrent/CyclicBarrier.getNumberWaiting()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cyclicBarrierGetNumberWaiting,
		}
}

func isThreadInterrupted(th *object.Object) bool {
	if th == nil || th == object.Null {
		return false
	}
	th.ThMutex.RLock()
	defer th.ThMutex.RUnlock()
	fld, ok := th.FieldTable["interrupted"]
	if !ok {
		return false
	}
	switch v := fld.Fvalue.(type) {
	case int64:
		return v == types.JavaBoolTrue
	case int:
		return v != 0
	default:
		return false
	}
}

func clearThreadInterrupted(th *object.Object) {
	if th == nil || th == object.Null {
		return
	}
	th.ThMutex.Lock()
	defer th.ThMutex.Unlock()
	fld, ok := th.FieldTable["interrupted"]
	if ok {
		fld.Fvalue = types.JavaBoolFalse
		th.FieldTable["interrupted"] = fld
	}
}

func getCyclicBarrierState(self *object.Object) (*cyclicBarrierState, interface{}) {
	if self == nil || object.IsNull(self) {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "getCyclicBarrierState: CyclicBarrier is null")
	}
	self.ThMutex.RLock()
	defer self.ThMutex.RUnlock()
	field, exists := self.FieldTable["state"]
	if !exists {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "getCyclicBarrierState: CyclicBarrier not initialized")
	}
	state, ok := field.Fvalue.(*cyclicBarrierState)
	if !ok || state == nil {
		return nil, ghelpers.GetGErrBlk(excNames.VirtualMachineError, "getCyclicBarrierState: Invalid CyclicBarrier storage")
	}
	return state, nil
}

func cyclicBarrierInit(params []interface{}) interface{} {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "CyclicBarrier.<init>: null parameters")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "CyclicBarrier.<init>: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "CyclicBarrier.<init>: missing parties parameter")
	}
	return cyclicBarrierInitAction([]interface{}{params[0], params[1], object.Null})
}

func cyclicBarrierInitAction(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "CyclicBarrier.<init>: insufficient parameters")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "CyclicBarrier.<init>: self is null")
	}
	parties, ok := params[1].(int64)
	if !ok || parties <= 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "CyclicBarrier parties must be positive")
	}
	var barrierAction *object.Object
	if len(params) > 2 {
		barrierAction, _ = params[2].(*object.Object)
	}

	state := &cyclicBarrierState{
		parties:       parties,
		count:         parties,
		gen:           &generation{broken: false},
		barrierAction: barrierAction,
	}
	state.barrierCond = sync.NewCond(&state.mu)

	self.ThMutex.Lock()
	defer self.ThMutex.Unlock()
	self.FieldTable["state"] = object.Field{Ftype: types.Ref, Fvalue: state}
	return nil
}

func cyclicBarrierAwait(params []interface{}) interface{} {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierAwait: null parameters")
	}

	var self *object.Object
	var currentThread *object.Object

	if fs, ok := params[0].(*list.List); ok {
		// NeedsContext=true: params = [fs, self]
		if len(params) < 2 {
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierAwait: missing self object")
		}
		var isObj bool
		self, isObj = params[1].(*object.Object)
		if !isObj || object.IsNull(self) {
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierAwait: null self object")
		}
		if fs != nil && fs.Front() != nil {
			if fr, ok := fs.Front().Value.(*frames.Frame); ok {
				gr := globals.GetGlobalRef()
				gr.ThreadLock.RLock()
				if th, exists := gr.Threads[fr.Thread]; exists && th != nil {
					currentThread, _ = th.(*object.Object)
				}
				gr.ThreadLock.RUnlock()
			}
		}
	} else if obj, ok := params[0].(*object.Object); ok {
		self = obj
		if len(params) > 1 {
			currentThread, _ = params[1].(*object.Object)
		}
	} else {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "cyclicBarrierAwait: invalid parameters")
	}

	if object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierAwait: null self object")
	}

	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	g := state.gen

	if g.broken {
		return ghelpers.GetGErrBlk(excNames.BrokenBarrierException, "CyclicBarrier is broken")
	}

	if isThreadInterrupted(currentThread) {
		clearThreadInterrupted(currentThread)
		state.breakBarrier()
		return ghelpers.GetGErrBlk(excNames.InterruptedException, "Thread interrupted before wait")
	}

	index := state.count - 1
	state.count = index

	if index == 0 {
		state.nextGeneration()
		return int64(0)
	}

	for {
		state.barrierCond.Wait()

		if isThreadInterrupted(currentThread) {
			clearThreadInterrupted(currentThread)
			if g == state.gen && !g.broken {
				state.breakBarrier()
				return ghelpers.GetGErrBlk(excNames.InterruptedException, "Thread interrupted during wait")
			}
			return ghelpers.GetGErrBlk(excNames.InterruptedException, "Thread interrupted during wait")
		}

		if g.broken {
			return ghelpers.GetGErrBlk(excNames.BrokenBarrierException, "CyclicBarrier broken or reset during wait")
		}

		if g != state.gen {
			return index
		}
	}
}

func cyclicBarrierGetParties(params []interface{}) interface{} {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierGetParties: null parameters")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierGetParties: null self object")
	}
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.parties
}

func cyclicBarrierIsBroken(params []interface{}) interface{} {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierIsBroken: null parameters")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierIsBroken: null self object")
	}
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}
	state.mu.RLock()
	defer state.mu.RUnlock()
	return object.JavaBooleanFromGoBoolean(state.gen.broken)
}

func cyclicBarrierReset(params []interface{}) interface{} {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierReset: null parameters")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierReset: null self object")
	}
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	state.breakBarrier()
	state.nextGeneration()

	return nil
}

func cyclicBarrierGetNumberWaiting(params []interface{}) interface{} {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierGetNumberWaiting: null parameters")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "cyclicBarrierGetNumberWaiting: null self object")
	}
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.parties - state.count
}
