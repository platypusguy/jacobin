package classloader

import (
	"jacobin/src/globals"
	"sync"
	"testing"
)

func TestJlcMapConcurrency(t *testing.T) {
	globals.InitGlobals("test")

	const numGoroutines = 100
	const numIterations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				pc := &ParsedClass{
					className: "Class" + string(rune(id)) + "_" + string(rune(j)),
				}
				// Simulate convertToPostableClass logic for JlcMap
				globals.JlcMapLock.Lock()
				globals.JlcMap[pc.className] = &Jlc{}
				globals.JlcMapLock.Unlock()
			}
		}(i)
	}

	wg.Wait()
}
