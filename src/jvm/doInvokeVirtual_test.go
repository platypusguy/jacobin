/*
 * Jacobin VM - A Java virtual machine
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"testing"
)

// This regression test currently fails: resetting globals also resets the string
// pool, but does not clear dispatchCache, so unrelated classes can reuse its keys.
func TestInvokeVirtualDispatchCacheStaleAfterStringPoolReset(t *testing.T) {
	var firstKey dispatchKey

	t.Run("first run", func(t *testing.T) {
		globals.InitGlobals("test")

		className := "ClassA"
		methodName := "run"
		methodType := "()V"
		firstKey = dispatchKey{
			recv: stringPool.GetStringIndex(&className),
			name: stringPool.GetStringIndex(&methodName),
			desc: stringPool.GetStringIndex(&methodType),
		}

		// Represent the resolved dispatch to ClassA.run()V.
		dispatchCache.Store(firstKey, &dispatchVal{className: "ClassA"})
	})

	t.Run("second run after reset", func(t *testing.T) {
		globals.InitGlobals("test")

		className := "ClassB"
		methodName := "run"
		methodType := "()V"
		secondKey := dispatchKey{
			recv: stringPool.GetStringIndex(&className),
			name: stringPool.GetStringIndex(&methodName),
			desc: stringPool.GetStringIndex(&methodType),
		}

		if secondKey != firstKey {
			t.Fatalf("expected string-pool reset to reproduce dispatch key %v, got %v", firstKey, secondKey)
		}

		cached, ok := dispatchCache.Load(secondKey)
		if !ok {
			t.Fatal("expected dispatch cache entry from first run to remain after reset")
		}
		if got := cached.(*dispatchVal).className; got != "ClassB" {
			t.Errorf("stale dispatch cache selected %s after reset; want ClassB", got)
		}
	})
}
