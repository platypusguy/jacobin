package classloader

import (
	"fmt"
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
					className: fmt.Sprintf("Class%d_%d", id, j),
				}
				// Simulate convertToPostableClass logic for JlcMap
				globals.JlcMapLock.Lock()
				globals.JlcMap[pc.className] = &Jlc{
					statics: make(map[string]Field),
				}
				globals.JlcMapLock.Unlock()
			}
		}(i)
	}

	wg.Wait()
}

func TestJlcElementConcurrency(t *testing.T) {
	globals.InitGlobals("test")

	jlc := &Jlc{
		statics: make(map[string]Field),
	}

	const numGoroutines = 50
	const numIterations = 500

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // One for writers, one for readers

	// Concurrent writers to statics and initialized flag
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				fieldName := fmt.Sprintf("field_%d_%d", id, j)
				f := Field{NameStr: fieldName}

				jlc.lock.Lock()
				jlc.statics[fieldName] = f
				jlc.initialized = !jlc.initialized
				jlc.lock.Unlock()
			}
		}(i)
	}

	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				fieldName := fmt.Sprintf("field_%d_%d", id, j)

				jlc.lock.RLock()
				_ = jlc.statics[fieldName]
				_ = jlc.initialized
				_ = jlc._klass
				jlc.lock.RUnlock()
			}
		}(i)
	}

	wg.Wait()
}

func TestJlcMapAndElementInteraction(t *testing.T) {
	globals.InitGlobals("test")

	const numClasses = 16
	const numFieldsPerClass = 3
	var classNames []string
	for i := 0; i < numClasses; i++ {
		classNames = append(classNames, fmt.Sprintf("Class%d", i))
	}

	var wg sync.WaitGroup
	wg.Add(numClasses)

	// Phase 1: 16 goroutines each adding a class and filling 3 fields
	for i := 0; i < numClasses; i++ {
		go func(idx int) {
			defer wg.Done()
			className := classNames[idx]

			jlc := &Jlc{
				statics: make(map[string]Field),
			}

			// Add the class to globals.JlcMap
			globals.JlcMapLock.Lock()
			globals.JlcMap[className] = jlc
			globals.JlcMapLock.Unlock()

			// Fill in 3 fields
			for j := 0; j < numFieldsPerClass; j++ {
				fieldName := fmt.Sprintf("field_%s_%d", className, j)
				f := Field{NameStr: fieldName}

				jlc.lock.Lock()
				jlc.statics[fieldName] = f
				jlc.lock.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Phase 2: 16 more goroutines to read the fields
	wg.Add(numClasses)
	for i := 0; i < numClasses; i++ {
		go func(idx int) {
			defer wg.Done()
			className := classNames[idx]

			// Access control for globals.JlcMap
			globals.JlcMapLock.RLock()
			jlcAny, ok := globals.JlcMap[className]
			globals.JlcMapLock.RUnlock()

			if !ok {
				t.Errorf("Class %s not found in JlcMap", className)
				return
			}
			jlc := jlcAny.(*Jlc)

			// Read the 3 fields
			for j := 0; j < numFieldsPerClass; j++ {
				fieldName := fmt.Sprintf("field_%s_%d", className, j)

				jlc.lock.RLock()
				f, exists := jlc.statics[fieldName]
				jlc.lock.RUnlock()

				if !exists {
					t.Errorf("Field %s not found in class %s", fieldName, className)
				} else if f.NameStr != fieldName {
					t.Errorf("Field name mismatch: expected %s, got %s", fieldName, f.NameStr)
				}
			}
		}(i)
	}

	wg.Wait()
}
