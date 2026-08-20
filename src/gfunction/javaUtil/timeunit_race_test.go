package javaUtil

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
	"testing"
)

func TestTimeUnit_Race(t *testing.T) {
	globals.InitStringPool()
	
	// Create a TimeUnit object (which in this VM seems to be represented by its name string as an object)
	// The implementation calls object.GoStringFromStringObject(unit)
	unit := object.StringObjectFromGoString(SECONDS)
	
	var wg sync.WaitGroup
	const goroutines = 10
	const iterations = 100
	
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				params := []interface{}{unit, int64(10)}
				_ = toMillis(params)
				_ = toSeconds(params)
				_ = toMinutes(params)
				_ = toHours(params)
				_ = toDays(params)
			}
		}()
	}
	
	// Also have a goroutine that "updates" the object to trigger a race if not protected
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < iterations; j++ {
			unit.ThMutex.Lock()
			unit.FieldTable["value"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.JavaByteArrayFromGoString(SECONDS)}
			unit.ThMutex.Unlock()
		}
	}()

	wg.Wait()
}
