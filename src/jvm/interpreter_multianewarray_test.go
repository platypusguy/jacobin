/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

func createMultinewarrayCP(desc string) *classloader.CPool {
	CP := &classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	descCopy := desc
	nameIndex := stringPool.GetStringIndex(&descCopy)
	CP.ClassRefs = append(CP.ClassRefs, nameIndex)
	CP.Utf8Refs = append(CP.Utf8Refs, desc)
	return CP
}

func setupMultinewarrayFrame(desc string, dimCount byte, dims ...int64) *frames.Frame {
	f := newFrame(opcodes.MULTIANEWARRAY)
	f.Meth = append(f.Meth, 0x00, 0x02, dimCount)
	f.CP = createMultinewarrayCP(desc)
	for _, d := range dims {
		push(f, d)
	}
	return f
}

// TestMultianewarray1D tests 1-dimensional array creation via MULTIANEWARRAY
func TestMultianewarray1D(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Standard 1D array", func(t *testing.T) {
		f := setupMultinewarrayFrame("[I", 1, 5)
		ret := doMultinewarray(f, 0)
		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		if f.TOS != 0 {
			t.Fatalf("Expected TOS 0, got %d", f.TOS)
		}
		arrObj := pop(f).(*object.Object)
		ftype := arrObj.FieldTable["value"].Ftype
		if ftype != types.IntArray {
			t.Errorf("Expected Ftype '%s', got '%s'", types.IntArray, ftype)
		}
		val := arrObj.FieldTable["value"].Fvalue.([]int64)
		if len(val) != 5 {
			t.Errorf("Expected length 5, got %d", len(val))
		}
	})

	t.Run("Zero-sized 1D array", func(t *testing.T) {
		f := setupMultinewarrayFrame("[I", 1, 0)
		ret := doMultinewarray(f, 0)
		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		ftype := arrObj.FieldTable["value"].Ftype
		if ftype != types.IntArray {
			t.Errorf("Expected Ftype '%s', got '%s'", types.IntArray, ftype)
		}
		val := arrObj.FieldTable["value"].Fvalue.([]int64)
		if len(val) != 0 {
			t.Errorf("Expected length 0, got %d", len(val))
		}
	})
}

// TestMultianewarray2D tests 2-dimensional array creation and zero-dimension boundary cases
func TestMultianewarray2D(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Standard 2D array (3x4)", func(t *testing.T) {
		f := setupMultinewarrayFrame("[[I", 2, 3, 4)
		ret := doMultinewarray(f, 0)
		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		if arrObj.FieldTable["value"].Ftype != "[[I" {
			t.Errorf("Expected top-level Ftype '[[I', got '%s'", arrObj.FieldTable["value"].Ftype)
		}
		outer := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
		if len(outer) != 3 {
			t.Fatalf("Expected outer length 3, got %d", len(outer))
		}
		for i, innerObj := range outer {
			if innerObj.FieldTable["value"].Ftype != types.IntArray {
				t.Errorf("Sub-array %d: expected Ftype '%s', got '%s'", i, types.IntArray, innerObj.FieldTable["value"].Ftype)
			}
			inner := innerObj.FieldTable["value"].Fvalue.([]int64)
			if len(inner) != 4 {
				t.Errorf("Sub-array %d: expected length 4, got %d", i, len(inner))
			}
		}
	})

	t.Run("2D array with 0 in first dimension (0x4)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[I", 2, 0, 4)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		// Collapses to 1D array of size 0
		if arrObj.FieldTable["value"].Ftype != types.IntArray {
			t.Errorf("Expected collapsed Ftype '%s', got '%s'", types.IntArray, arrObj.FieldTable["value"].Ftype)
		}
		inner := arrObj.FieldTable["value"].Fvalue.([]int64)
		if len(inner) != 0 {
			t.Errorf("Expected length 0, got %d", len(inner))
		}
	})

	t.Run("2D array with 0 in second dimension (3x0)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[I", 2, 3, 0)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		if arrObj.FieldTable["value"].Ftype != "[[I" {
			t.Errorf("Expected top-level Ftype '[[I', got '%s'", arrObj.FieldTable["value"].Ftype)
		}
		outer := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
		if len(outer) != 3 {
			t.Fatalf("Expected outer length 3, got %d", len(outer))
		}
		for i, innerObj := range outer {
			if innerObj.FieldTable["value"].Ftype != types.IntArray {
				t.Errorf("Sub-array %d: expected Ftype '%s', got '%s'", i, types.IntArray, innerObj.FieldTable["value"].Ftype)
			}
			inner := innerObj.FieldTable["value"].Fvalue.([]int64)
			if len(inner) != 0 {
				t.Errorf("Sub-array %d: expected length 0, got %d", i, len(inner))
			}
		}
	})

	t.Run("2D array with 0 in both dimensions (0x0)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[I", 2, 0, 0)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		if arrObj.FieldTable["value"].Ftype != types.IntArray {
			t.Errorf("Expected collapsed Ftype '%s', got '%s'", types.IntArray, arrObj.FieldTable["value"].Ftype)
		}
		inner := arrObj.FieldTable["value"].Fvalue.([]int64)
		if len(inner) != 0 {
			t.Errorf("Expected length 0, got %d", len(inner))
		}
	})
}

// TestMultianewarray3D tests 3-dimensional array creation and zero-dimension boundary cases
func TestMultianewarray3D(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Standard 3D array (2x3x4)", func(t *testing.T) {
		f := setupMultinewarrayFrame("[[[I", 3, 2, 3, 4)
		ret := doMultinewarray(f, 0)
		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		if arrObj.FieldTable["value"].Ftype != types.RefArray {
			t.Errorf("Expected 1st dim Ftype '%s', got '%s'", types.RefArray, arrObj.FieldTable["value"].Ftype)
		}
		dim1 := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
		if len(dim1) != 2 {
			t.Fatalf("Expected 1st dim length 2, got %d", len(dim1))
		}
		for i, d2Obj := range dim1 {
			if d2Obj.FieldTable["value"].Ftype != "[[I" {
				t.Errorf("Dim1[%d]: expected 2nd dim Ftype '[[I', got '%s'", i, d2Obj.FieldTable["value"].Ftype)
			}
			dim2 := d2Obj.FieldTable["value"].Fvalue.([]*object.Object)
			if len(dim2) != 3 {
				t.Fatalf("Dim1[%d]: expected 2nd dim length 3, got %d", i, len(dim2))
			}
			for j, leafObj := range dim2 {
				if leafObj.FieldTable["value"].Ftype != types.IntArray {
					t.Errorf("Dim1[%d] Dim2[%d]: expected leaf Ftype '%s', got '%s'", i, j, types.IntArray, leafObj.FieldTable["value"].Ftype)
				}
				leaf := leafObj.FieldTable["value"].Fvalue.([]int64)
				if len(leaf) != 4 {
					t.Errorf("Dim1[%d] Dim2[%d]: expected leaf length 4, got %d", i, j, len(leaf))
				}
			}
		}
	})

	t.Run("3D array with 0 in first dimension (0x3x4)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[[I", 3, 0, 3, 4)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		// Collapses to 1D array of size 0
		if arrObj.FieldTable["value"].Ftype != types.IntArray {
			t.Errorf("Expected collapsed Ftype '%s', got '%s'", types.IntArray, arrObj.FieldTable["value"].Ftype)
		}
		val := arrObj.FieldTable["value"].Fvalue.([]int64)
		if len(val) != 0 {
			t.Errorf("Expected length 0, got %d", len(val))
		}
	})

	t.Run("3D array with 0 in second dimension (2x0x4)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[[I", 3, 2, 0, 4)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		// Collapses to 2D array (2 elements of 0-length arrays)
		if arrObj.FieldTable["value"].Ftype != "[[I" {
			t.Errorf("Expected Ftype '[[I', got '%s'", arrObj.FieldTable["value"].Ftype)
		}
		dim1 := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
		if len(dim1) != 2 {
			t.Fatalf("Expected length 2, got %d", len(dim1))
		}
		for i, innerObj := range dim1 {
			if innerObj.FieldTable["value"].Ftype != types.IntArray {
				t.Errorf("Dim1[%d]: expected sub-array Ftype '%s', got '%s'", i, types.IntArray, innerObj.FieldTable["value"].Ftype)
			}
			inner := innerObj.FieldTable["value"].Fvalue.([]int64)
			if len(inner) != 0 {
				t.Errorf("Dim1[%d]: expected length 0, got %d", i, len(inner))
			}
		}
	})

	t.Run("3D array with 0 in third dimension (2x3x0)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[[I", 3, 2, 3, 0)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		if arrObj.FieldTable["value"].Ftype != types.RefArray {
			t.Errorf("Expected 1st dim Ftype '%s', got '%s'", types.RefArray, arrObj.FieldTable["value"].Ftype)
		}
		dim1 := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
		if len(dim1) != 2 {
			t.Fatalf("Expected 1st dim length 2, got %d", len(dim1))
		}
		for i, d2Obj := range dim1 {
			if d2Obj.FieldTable["value"].Ftype != "[[I" {
				t.Errorf("Dim1[%d]: expected 2nd dim Ftype '[[I', got '%s'", i, d2Obj.FieldTable["value"].Ftype)
			}
			dim2 := d2Obj.FieldTable["value"].Fvalue.([]*object.Object)
			if len(dim2) != 3 {
				t.Fatalf("Dim1[%d]: expected 2nd dim length 3, got %d", i, len(dim2))
			}
			for j, leafObj := range dim2 {
				if leafObj.FieldTable["value"].Ftype != types.IntArray {
					t.Errorf("Dim1[%d] Dim2[%d]: expected leaf Ftype '%s', got '%s'", i, j, types.IntArray, leafObj.FieldTable["value"].Ftype)
				}
				leaf := leafObj.FieldTable["value"].Fvalue.([]int64)
				if len(leaf) != 0 {
					t.Errorf("Dim1[%d] Dim2[%d]: expected leaf length 0, got %d", i, j, len(leaf))
				}
			}
		}
	})

	t.Run("3D array with 0 in all dimensions (0x0x0)", func(t *testing.T) {
		normalStderr := os.Stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[[I", 3, 0, 0, 0)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		os.Stderr = normalStderr

		if ret != 4 {
			t.Fatalf("Expected return 4, got %d", ret)
		}
		arrObj := pop(f).(*object.Object)
		if arrObj.FieldTable["value"].Ftype != types.IntArray {
			t.Errorf("Expected collapsed Ftype '%s', got '%s'", types.IntArray, arrObj.FieldTable["value"].Ftype)
		}
		val := arrObj.FieldTable["value"].Fvalue.([]int64)
		if len(val) != 0 {
			t.Errorf("Expected length 0, got %d", len(val))
		}
	})
}

// TestMultianewarrayAllPrimitiveTypesAndRef tests MULTIANEWARRAY with all supported element types
func TestMultianewarrayAllPrimitiveTypesAndRef(t *testing.T) {
	globals.InitGlobals("test")

	testCases := []struct {
		name         string
		desc3D       string
		desc2D       string
		desc1D       string
		expectedLeaf string
		checkLeaf    func(t *testing.T, leafObj *object.Object)
	}{
		{
			name:         "Boolean (Z)",
			desc3D:       "[[[Z",
			desc2D:       "[[Z",
			desc1D:       "[Z",
			expectedLeaf: types.BoolArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]types.JavaByte)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []types.JavaByte of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Byte (B)",
			desc3D:       "[[[B",
			desc2D:       "[[B",
			desc1D:       "[B",
			expectedLeaf: types.JavaByteArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]types.JavaByte)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []types.JavaByte of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Char (C)",
			desc3D:       "[[[C",
			desc2D:       "[[C",
			desc1D:       "[C",
			expectedLeaf: types.CharArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]int64)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []int64 of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Double (D)",
			desc3D:       "[[[D",
			desc2D:       "[[D",
			desc1D:       "[D",
			expectedLeaf: types.DoubleArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]float64)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []float64 of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Float (F)",
			desc3D:       "[[[F",
			desc2D:       "[[F",
			desc1D:       "[F",
			expectedLeaf: types.FloatArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]float64)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []float64 of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Int (I)",
			desc3D:       "[[[I",
			desc2D:       "[[I",
			desc1D:       "[I",
			expectedLeaf: types.IntArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]int64)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []int64 of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Long (J)",
			desc3D:       "[[[J",
			desc2D:       "[[J",
			desc1D:       "[J",
			expectedLeaf: types.LongArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]int64)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []int64 of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Short (S)",
			desc3D:       "[[[S",
			desc2D:       "[[S",
			desc1D:       "[S",
			expectedLeaf: types.ShortArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]int64)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []int64 of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
		{
			name:         "Reference (L)",
			desc3D:       "[[[Ljava/lang/String;",
			desc2D:       "[[Ljava/lang/String;",
			desc1D:       "[Ljava/lang/String;",
			expectedLeaf: types.RefArray,
			checkLeaf: func(t *testing.T, leafObj *object.Object) {
				val, ok := leafObj.FieldTable["value"].Fvalue.([]*object.Object)
				if !ok || len(val) != 2 {
					t.Errorf("Expected []*object.Object of length 2, got %T %v", leafObj.FieldTable["value"].Fvalue, val)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+" 3D", func(t *testing.T) {
			f := setupMultinewarrayFrame(tc.desc3D, 3, 2, 2, 2)
			ret := doMultinewarray(f, 0)
			if ret != 4 {
				t.Fatalf("Expected return 4, got %d", ret)
			}
			top := pop(f).(*object.Object)
			dim1 := top.FieldTable["value"].Fvalue.([]*object.Object)
			dim2 := dim1[0].FieldTable["value"].Fvalue.([]*object.Object)
			leafObj := dim2[0]
			if leafObj.FieldTable["value"].Ftype != tc.expectedLeaf {
				t.Errorf("Expected leaf Ftype '%s', got '%s'", tc.expectedLeaf, leafObj.FieldTable["value"].Ftype)
			}
			tc.checkLeaf(t, leafObj)
		})

		t.Run(tc.name+" 2D", func(t *testing.T) {
			f := setupMultinewarrayFrame(tc.desc2D, 2, 2, 2)
			ret := doMultinewarray(f, 0)
			if ret != 4 {
				t.Fatalf("Expected return 4, got %d", ret)
			}
			top := pop(f).(*object.Object)
			dim1 := top.FieldTable["value"].Fvalue.([]*object.Object)
			leafObj := dim1[0]
			if leafObj.FieldTable["value"].Ftype != tc.expectedLeaf {
				t.Errorf("Expected leaf Ftype '%s', got '%s'", tc.expectedLeaf, leafObj.FieldTable["value"].Ftype)
			}
			tc.checkLeaf(t, leafObj)
		})

		t.Run(tc.name+" 1D", func(t *testing.T) {
			f := setupMultinewarrayFrame(tc.desc1D, 1, 2)
			ret := doMultinewarray(f, 0)
			if ret != 4 {
				t.Fatalf("Expected return 4, got %d", ret)
			}
			leafObj := pop(f).(*object.Object)
			if leafObj.FieldTable["value"].Ftype != tc.expectedLeaf {
				t.Errorf("Expected leaf Ftype '%s', got '%s'", tc.expectedLeaf, leafObj.FieldTable["value"].Ftype)
			}
			tc.checkLeaf(t, leafObj)
		})
	}
}

// TestMultianewarrayBoundaryDimCounts tests boundary dimension counts (0, >3)
func TestMultianewarrayBoundaryDimCounts(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Dimension count 0", func(t *testing.T) {
		f := setupMultinewarrayFrame("[[[I", 0)
		ret := doMultinewarray(f, 0)
		if ret != 4 {
			t.Errorf("Expected return 4 for dimCount 0, got %d", ret)
		}
		if f.TOS != -1 {
			t.Errorf("Expected stack unchanged (TOS -1), got %d", f.TOS)
		}
	})

	t.Run("Dimension count 4 (exceeds limit 3)", func(t *testing.T) {
		normalStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[[[I", 4, 1, 1, 1, 1)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		msg, _ := io.ReadAll(r)
		os.Stderr = normalStderr

		if ret != ERROR_OCCURRED {
			t.Errorf("Expected ERROR_OCCURRED (-1), got %d", ret)
		}
		errMsg := string(msg)
		if !strings.Contains(errMsg, "Jacobin supports arrays only up to three dimensions") {
			t.Errorf("Expected error message about dimension limit, got: %s", errMsg)
		}
	})

	t.Run("Dimension count 255 (exceeds limit 3)", func(t *testing.T) {
		normalStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		f := setupMultinewarrayFrame("[[[[I", 255, 1)
		ret := doMultinewarray(f, 0)

		_ = w.Close()
		msg, _ := io.ReadAll(r)
		os.Stderr = normalStderr

		if ret != ERROR_OCCURRED {
			t.Errorf("Expected ERROR_OCCURRED (-1), got %d", ret)
		}
		errMsg := string(msg)
		if !strings.Contains(errMsg, "Jacobin supports arrays only up to three dimensions") {
			t.Errorf("Expected error message about dimension limit, got: %s", errMsg)
		}
	})
}

// TestMultianewarrayInvalidTypes tests unexpected/invalid raw array types
func TestMultianewarrayInvalidTypes(t *testing.T) {
	globals.InitGlobals("test")

	invalidTypes := []string{
		"[[[X",
		"[[[?",
		"[[V",
		"[[1",
	}

	for _, invalidDesc := range invalidTypes {
		t.Run("Invalid type "+invalidDesc, func(t *testing.T) {
			normalStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			f := setupMultinewarrayFrame(invalidDesc, 2, 2, 2)
			ret := doMultinewarray(f, 0)

			_ = w.Close()
			msg, _ := io.ReadAll(r)
			os.Stderr = normalStderr

			if ret != ERROR_OCCURRED {
				t.Errorf("Expected ERROR_OCCURRED (-1) for %s, got %d", invalidDesc, ret)
			}
			errMsg := string(msg)
			if !strings.Contains(errMsg, "MULTIANEWARRAY: Unexpected raw array type") {
				t.Errorf("Expected error message about unexpected raw array type, got: %s", errMsg)
			}
		})
	}
}

// TestMultianewarrayInterpreterLoop tests MULTIANEWARRAY executed via interpret() loop
func TestMultianewarrayInterpreterLoop(t *testing.T) {
	globals.InitGlobals("test")

	f := setupMultinewarrayFrame("[[[I", 3, 2, 3, 4)
	fs := frames.CreateFrameStack()
	fs.PushFront(f)

	interpret(fs)

	if f.TOS != 0 {
		t.Fatalf("Expected TOS 0 after interpret, got %d", f.TOS)
	}

	arrObj := pop(f).(*object.Object)
	if arrObj == nil {
		t.Fatal("Expected non-nil array object from interpret")
	}
	dim1 := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(dim1) != 2 {
		t.Errorf("Expected outer dimension 2, got %d", len(dim1))
	}
}
