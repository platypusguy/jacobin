package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"testing"
)

var testSep = string(os.PathSeparator)

func TestFilePathCompareTo(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p2 := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p3 := newPath(fmt.Sprintf("%sa%sb%sd", testSep, testSep, testSep))

	// p1.compareTo(p2) == 0
	res := filePathCompareTo([]interface{}{p1, p2})
	if res.(int64) != 0 {
		t.Errorf("expected 0, got %d", res)
	}

	// p1.compareTo(p3) < 0
	res = filePathCompareTo([]interface{}{p1, p3})
	if res.(int64) >= 0 {
		t.Errorf("expected < 0, got %d", res)
	}

	// p3.compareTo(p1) > 0
	res = filePathCompareTo([]interface{}{p3, p1})
	if res.(int64) <= 0 {
		t.Errorf("expected > 0, got %d", res)
	}

	// null checks
	res = filePathCompareTo([]interface{}{nil, p2})
	if ge, ok := res.(*GErrBlk); !ok || ge.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NPE for this=null")
	}

	res = filePathCompareTo([]interface{}{p1, nil})
	if ge, ok := res.(*GErrBlk); !ok || ge.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NPE for other=null")
	}
}

func TestFilePathEndsWith(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	s1 := object.StringObjectFromGoString("c")
	s2 := object.StringObjectFromGoString(fmt.Sprintf("b%sc", testSep))
	s3 := object.StringObjectFromGoString("a")

	if filePathEndsWith([]interface{}{p, s1}) != types.JavaBoolTrue {
		t.Errorf("expected true for 'c'")
	}
	if filePathEndsWith([]interface{}{p, s2}) != types.JavaBoolTrue {
		t.Errorf("expected true for 'b%sc'", testSep)
	}
	if filePathEndsWith([]interface{}{p, s3}) != types.JavaBoolFalse {
		t.Errorf("expected false for 'a'")
	}
}

func TestFilePathEndsWithPath(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p1 := newPath("c")
	p2 := newPath(fmt.Sprintf("b%sc", testSep))
	p3 := newPath("a")

	if filePathEndsWithPath([]interface{}{p, p1}) != types.JavaBoolTrue {
		t.Errorf("expected true for path 'c'")
	}
	if filePathEndsWithPath([]interface{}{p, p2}) != types.JavaBoolTrue {
		t.Errorf("expected true for path 'b%sc'", testSep)
	}
	if filePathEndsWithPath([]interface{}{p, p3}) != types.JavaBoolFalse {
		t.Errorf("expected false for path 'a'")
	}
}

func TestFilePathEquals(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p2 := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p3 := newPath(fmt.Sprintf("%sa%sb%sd", testSep, testSep, testSep))

	if filePathEquals([]interface{}{p1, p2}) != types.JavaBoolTrue {
		t.Errorf("expected true for equal paths")
	}
	if filePathEquals([]interface{}{p1, p3}) != types.JavaBoolFalse {
		t.Errorf("expected false for different paths")
	}
	if filePathEquals([]interface{}{p1, object.StringObjectFromGoString(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))}) != types.JavaBoolFalse {
		t.Errorf("expected false for different object type")
	}
}

func TestFilePathGetFileName(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathGetFileName([]interface{}{p1}).(*object.Object)
	if object.GoStringFromStringObject(res) != "c" {
		t.Errorf("expected 'c', got %s", object.GoStringFromStringObject(res))
	}

	p2 := newPath(testSep)
	res2 := filePathGetFileName([]interface{}{p2})
	if !object.IsNull(res2) {
		t.Errorf("expected null for root, got %v", res2)
	}
}

func TestFilePathGetName(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathGetName([]interface{}{p, int64(0)}).(*object.Object)
	if object.GoStringFromStringObject(res) != "a" {
		t.Errorf("expected 'a', got %s", object.GoStringFromStringObject(res))
	}
}

func TestFilePathGetNameCount(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathGetNameCount([]interface{}{p})
	// getPathParts("/a/b/c") -> ["a", "b", "c"] -> length 3
	if res.(int64) != 3 {
		t.Errorf("expected 3, got %d", res)
	}
}

func TestFilePathGetParent(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathGetParent([]interface{}{p}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sb", testSep, testSep) {
		t.Errorf("expected '%sa%sb', got %s", testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathGetRoot(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathGetRoot([]interface{}{p}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != testSep {
		t.Errorf("expected '%s', got %s", testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathHashCode(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathHashCode([]interface{}{p})
	if res.(int64) == 0 {
		t.Errorf("expected non-zero hash code")
	}
}

func TestFilePathIsAbsolute(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p2 := newPath(fmt.Sprintf("a%sb%sc", testSep, testSep))
	if filePathIsAbsolute([]interface{}{p1}) != types.JavaBoolTrue {
		t.Errorf("expected true for absolute path")
	}
	if filePathIsAbsolute([]interface{}{p2}) != types.JavaBoolFalse {
		t.Errorf("expected false for relative path")
	}
}

func TestFilePathNormalize(t *testing.T) {
	p2 := newPath(fmt.Sprintf("%sa%s%sb", testSep, testSep, testSep))
	res := filePathNormalize([]interface{}{p2}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sb", testSep, testSep) {
		t.Errorf("expected '%sa%sb', got %s", testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathRelativize(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb", testSep, testSep))
	p2 := newPath(fmt.Sprintf("%sa%sb%sc%sd", testSep, testSep, testSep, testSep))
	res := filePathRelativize([]interface{}{p1, p2}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("c%sd", testSep) {
		t.Errorf("expected 'c%sd', got %s", testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathResolve(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb", testSep, testSep))
	s := object.StringObjectFromGoString("c")
	res := filePathResolve([]interface{}{p, s}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep) {
		t.Errorf("expected '%sa%sb%sc', got %s", testSep, testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathResolvePath(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb", testSep, testSep))
	p2 := newPath("c")
	res := filePathResolvePath([]interface{}{p1, p2}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep) {
		t.Errorf("expected '%sa%sb%sc', got %s", testSep, testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathResolveSibling(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb", testSep, testSep))
	s := object.StringObjectFromGoString("c")
	res := filePathResolveSibling([]interface{}{p, s}).(*object.Object)
	// parent of /a/b is /a. /a + / + c = /a/c
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sc", testSep, testSep) {
		t.Errorf("expected '%sa%sc', got %s", testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathResolveSiblingPath(t *testing.T) {
	p1 := newPath(fmt.Sprintf("%sa%sb", testSep, testSep))
	p2 := newPath("c")
	res := filePathResolveSiblingPath([]interface{}{p1, p2}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sc", testSep, testSep) {
		t.Errorf("expected '%sa%sc', got %s", testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathStartsWith(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	s1 := object.StringObjectFromGoString(fmt.Sprintf("%sa", testSep))
	s2 := object.StringObjectFromGoString(fmt.Sprintf("%sa%sb", testSep, testSep))
	s3 := object.StringObjectFromGoString("b")

	if filePathStartsWith([]interface{}{p, s1}) != types.JavaBoolTrue {
		t.Errorf("expected true for prefix1")
	}
	if filePathStartsWith([]interface{}{p, s2}) != types.JavaBoolTrue {
		t.Errorf("expected true for prefix2")
	}
	if filePathStartsWith([]interface{}{p, s3}) != types.JavaBoolFalse {
		t.Errorf("expected false for prefix3")
	}
}

func TestFilePathStartsWithPath(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	p1 := newPath(fmt.Sprintf("%sa", testSep))
	p2 := newPath(fmt.Sprintf("%sa%sb", testSep, testSep))
	p3 := newPath("b")

	if filePathStartsWithPath([]interface{}{p, p1}) != types.JavaBoolTrue {
		t.Errorf("expected true for path prefix1")
	}
	if filePathStartsWithPath([]interface{}{p, p2}) != types.JavaBoolTrue {
		t.Errorf("expected true for path prefix2")
	}
	if filePathStartsWithPath([]interface{}{p, p3}) != types.JavaBoolFalse {
		t.Errorf("expected false for path prefix3")
	}
}

func TestFilePathSubpath(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc%sd", testSep, testSep, testSep, testSep))
	// getPathParts("/a/b/c/d") -> ["a", "b", "c", "d"]
	res := filePathSubpath([]interface{}{p, int64(0), int64(2)}).(*object.Object)
	// parts[0:2] -> ["a", "b"] -> "a/b"
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("a%sb", testSep) {
		t.Errorf("expected 'a%sb', got %s", testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathToAbsolutePath(t *testing.T) {
	globals.InitGlobals("test")
	globals.SetSystemProperty("user.dir", testSep)
	p := newPath(fmt.Sprintf("a%sb", testSep))
	res := filePathToAbsolutePath([]interface{}{p}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sb", testSep, testSep) {
		t.Errorf("expected '%sa%sb', got %s", testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathToRealPath(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%s%sb", testSep, testSep, testSep))
	res := filePathToRealPath([]interface{}{p}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != fmt.Sprintf("%sa%sb", testSep, testSep) {
		t.Errorf("expected '%sa%sb', got %s", testSep, testSep, object.GoStringFromStringObject(val))
	}
}

func TestFilePathIterator(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathIterator([]interface{}{p}).(*object.Object)
	// res is a java/util/LinkedList containing String objects
	// We just check if it's not nil
	if object.IsNull(res) {
		t.Errorf("expected non-null java/util/LinkedList")
	}
}

func TestFilePathToString(t *testing.T) {
	p := newPath(fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep))
	res := filePathToString([]interface{}{p}).(*object.Object)
	if object.GoStringFromStringObject(res) != fmt.Sprintf("%sa%sb%sc", testSep, testSep, testSep) {
		t.Errorf("expected '%sa%sb%sc', got %s", testSep, testSep, testSep, object.GoStringFromStringObject(res))
	}
}
