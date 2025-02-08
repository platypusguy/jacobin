package gfunction

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
	"math/big"
	"reflect"
	"testing"
)

/***

TGrunner
========

Parameters:
	t *testing.T - Go's standard unit test block
    className, methodName, methodType string - FQN components
    expected interface{} - Expected result
    args []interface{} - Parameters for the method

Return value: none

***/

var FlagTGinit = false

func TGrunner(t *testing.T, className, methodName, methodType string,
	expected interface{},
	obj *object.Object, // can be nil
	args []interface{}) {

	// String form of FQN.
	fqn := fmt.Sprintf("%s.%s%s", className, methodName, methodType)

	// Initialize Jacobin classloader infrastructure.
	if !FlagTGinit {
		if !TGinit(t) {
			return
		}
	}
	// Create empty frame stack (fs).
	fs := frames.CreateFrameStack()

	// Create frame (fr).
	fr := frames.CreateFrame(3)
	fr.Thread = 0 // Mainthread
	fr.FrameStack = fs
	fr.ClName = className
	fr.MethName = methodName
	fr.MethType = methodType

	// Add CP to the frame.
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0} // should be a method ref
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		ClName:  "testClass",
		FldName: "testField",
		FldType: "I",
	}
	fr.CP = &CP

	// Push fr to front of fs.
	_ = frames.PushFrame(fs, fr)

	// Load the G functions.
	classloader.MTable = make(map[string]classloader.MTentry)
	MTableLoadGFunctions(&classloader.MTable)

	// Create mtEntry.
	mtEntry := classloader.MTable[fqn]
	if mtEntry.Meth == nil { // if the method is not yet in the method table, find it
		t.Errorf("TGrunner ERROR): classloader.MTable[%s] not found", fqn)
		return
	}

	paramCount := len(args)

	// params = args in reverse order (expected by RunGfunction).
	params := make([]interface{}, paramCount)
	for ix := 0; ix < paramCount; ix++ {
		params[ix] = args[paramCount-1-ix]
	}

	// Add the object reference (Java class or file I/O).
	if obj != nil {
		params = append(params, obj)
	}

	// Run the G function.
	observed := RunGfunction(mtEntry, fs, className, methodName, methodType, &params, true, false)

	// Check for nil result.
	if observed == nil {
		if expected != nil {
			t.Errorf("TGrunner ERROR: FQN %s returned nil but caller expected %T", fqn, expected)
			return
		}
	} else {
		// Not nil. Check for two types of errors.
		switch observed.(type) {
		case error:
			t.Errorf("TGrunner ERROR: FQN %s returned an error, text: %s", fqn, observed.(error).Error())
			return
		case *GErrBlk:
			t.Errorf("TGrunner ERROR: FQN %s returned a GErrBlk, text: %s", fqn, observed.(*GErrBlk).ErrMsg)
			return
		default:
			// Not any kind of error.
			// Make sure that the observed value is the same type as the expected.
			if reflect.TypeOf(observed) != reflect.TypeOf(expected) {
				t.Errorf("TGrunner ERROR: FQN %s expected return type %T, observed type %T", fqn, expected, observed)
				return
			}
			// Go data types agree.
			// Check the Jacobin field types and values.
			switch observed.(type) {
			case *object.Object:
				// Some type of object. Check the value field.
				fobs, okobs := observed.(*object.Object).FieldTable["value"]
				if okobs {
					fexp, okexp := expected.(*object.Object).FieldTable["value"]
					if okexp {
						// Both have a value field. So far, so good.
						// Compare the field types.
						if fobs.Ftype != fexp.Ftype {
							t.Errorf("TGrunner ERROR: FQN %s observed field type %s != expected field type %s",
								fqn, fobs.Ftype, fexp.Ftype)
							return
						}
						// Special checking for BigInteger values.
						if fexp.Ftype == types.BigInteger {
							biexp := fexp.Fvalue.(*big.Int)
							biobs := fexp.Fvalue.(*big.Int)
							if biexp.Cmp(biobs) != 0 {
								t.Errorf("TGrunner ERROR: FQN %s observed field value %s != expected field value %s",
									fqn, biobs.String(), biexp.String())
								return
							}
						} else {
							// For every other Field Ftype .....
							// Compare the 2 Fvalues, byte for byte, regardless of type.
							// TODO: This does not work for an array of objects!
							if !CompareBlobs(fobs, fexp) {
								t.Errorf("TGrunner ERROR: FQN %s observed return value %v != expected value %v",
									fqn, expected, observed)
								return
							}
						}
					} else {
						t.Errorf("TGrunner ERROR: FQN %s expected object is missing field \"value\"", fqn)
						return
					}
				} else {
					t.Errorf("TGrunner ERROR: FQN %s observed object is missing field \"value\"", fqn)
				}
			default:
				// Not an object. Assume it is a simple variable like an int64.
				if observed != expected {
					t.Errorf("TGrunner ERROR: FQN %s observed return value %v != expected value %v",
						fqn, expected, observed)
				}
			}
		}
	}

	t.Logf("TGrunner SUCCESS: FQN %s", fqn)

}

// Initialize Jacobin classloader infrastructure.
func TGinit(t *testing.T) bool {
	globals.InitGlobals("test")
	statics.Statics = make(map[string]statics.Static)
	err := classloader.Init()
	if err != nil {
		t.Errorf("TGinit ERROR): classloader.Init() failed")
		return false
	}
	classloader.LoadBaseClasses() // must follow classloader.Init()
	FlagTGinit = true
	return true
}

// CompareBlobs compares two anythings byte-for-byte
func AnyToBytes(blob interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(blob)
	return buf.Bytes()
}
func CompareBlobs(tweedleDee, tweedleDum interface{}) bool {
	return bytes.Equal(AnyToBytes(tweedleDee), AnyToBytes(tweedleDum))
}
