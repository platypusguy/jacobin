/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
)

type cyclicBarrierState struct {
	parties       int
	count         int
	generation    int
	broken        bool
	lastBroken    bool
	barrierCond   *sync.Cond
	barrierAction *object.Object
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
			ParamSlots: 0,
			GFunction:  cyclicBarrierAwait,
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

func getCyclicBarrierState(self *object.Object) (*cyclicBarrierState, interface{}) {
	field, exists := self.FieldTable["state"]
	if !exists {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "getCyclicBarrierState: CyclicBarrier not initialized")
	}
	state, ok := field.Fvalue.(*cyclicBarrierState)
	if !ok {
		return nil, ghelpers.GetGErrBlk(excNames.VirtualMachineError, "getCyclicBarrierState: Invalid CyclicBarrier storage")
	}
	return state, nil
}

func cyclicBarrierInit(params []interface{}) interface{} {
	return cyclicBarrierInitAction([]interface{}{params[0], params[1], object.Null})
}

func cyclicBarrierInitAction(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	parties, ok := params[1].(int64)
	if !ok || parties <= 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "CyclicBarrier parties must be positive")
	}
	barrierAction, _ := params[2].(*object.Object)

	mu := &sync.Mutex{}
	state := &cyclicBarrierState{
		parties:       int(parties),
		count:         int(parties),
		generation:    0,
		broken:        false,
		barrierCond:   sync.NewCond(mu),
		barrierAction: barrierAction,
	}

	self.FieldTable["state"] = object.Field{Ftype: types.ArrayList, Fvalue: state} // Using ArrayList type as a placeholder for pointer
	return nil
}

func cyclicBarrierAwait(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}

	state.barrierCond.L.Lock()
	defer state.barrierCond.L.Unlock()

	generation := state.generation

	if state.broken {
		return ghelpers.GetGErrBlk(excNames.BrokenBarrierException, "CyclicBarrier is broken")
	}

	// In a real JVM, we should also check for thread interruption

	index := state.count - 1
	state.count = index

	if index == 0 {
		// Last thread arrived
		ranAction := false
		if state.barrierAction != nil && state.barrierAction != object.Null {
			// Action execution not fully implemented (requires Java call)
		}

		if !ranAction {
			// Success
			state.lastBroken = false
			state.generation++
			state.count = state.parties
			state.barrierCond.Broadcast()
			return int64(0)
		}
		// If action failed, break barrier
		state.broken = true
		state.lastBroken = true
		state.barrierCond.Broadcast()
		return ghelpers.GetGErrBlk(excNames.BrokenBarrierException, "Barrier action failed")
	}

	// Wait for others
	for generation == state.generation {
		state.barrierCond.Wait()
		if generation != state.generation {
			if state.lastBroken {
				return ghelpers.GetGErrBlk(excNames.BrokenBarrierException, "CyclicBarrier broken or reset during wait")
			}
			return int64(index)
		}
	}

	return int64(index)
}

func cyclicBarrierGetParties(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}
	return int64(state.parties)
}

func cyclicBarrierIsBroken(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}
	return object.JavaBooleanFromGoBoolean(state.broken)
}

func cyclicBarrierReset(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}

	state.barrierCond.L.Lock()
	defer state.barrierCond.L.Unlock()

	state.broken = true
	state.lastBroken = true
	state.generation++ // Mark this generation as done/broken
	state.barrierCond.Broadcast()

	state.count = state.parties
	state.broken = false // Start new generation fresh

	return nil
}

func cyclicBarrierGetNumberWaiting(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getCyclicBarrierState(self)
	if err != nil {
		return err
	}
	state.barrierCond.L.Lock()
	defer state.barrierCond.L.Unlock()
	return int64(state.parties - state.count)
}
