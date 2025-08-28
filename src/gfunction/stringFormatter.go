package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"runtime"
	"strconv"
	"strings"
)

// helper wrapper to keep Java integral bit-width for formatting
// bits=32 for int/short/byte, bits=64 for long
type intWithBits struct {
	v    int64
	bits int
}

// String formatting given a format string and a slice of arguments.
// Called by sprintf, javaIoConsole.go, and javaIoPrintStream.go.
func StringFormatter(params []interface{}) interface{} {
	// params[0]: format string
	// params[1]: argument slice (array of object pointers)

	// Check the parameter length. It should be 2.
	lenParams := len(params)
	if lenParams < 1 || lenParams > 2 {
		errMsg := fmt.Sprintf("StringFormatter: Invalid parameter count: %d", lenParams)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if lenParams == 1 { // No parameters beyond the format string
		formatStringObj := params[0].(*object.Object)
		return formatStringObj
	}

	// Check the format string.
	var formatString string
	switch params[0].(type) {
	case *object.Object:
		formatStringObj := params[0].(*object.Object) // the format string is passed as a pointer to a string object
		switch formatStringObj.FieldTable["value"].Fvalue.(type) {
		case []byte:
			formatString = object.GoStringFromStringObject(formatStringObj)
		case []types.JavaByte:
			formatString =
				object.GoStringFromJavaByteArray(formatStringObj.FieldTable["value"].Fvalue.([]types.JavaByte))
		default:
			errMsg := fmt.Sprintf("StringFormatter: In the format string object, expected Ftype=%s but observed: %s",
				types.ByteArray, formatStringObj.FieldTable["value"].Ftype)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	default:
		errMsg := fmt.Sprintf("StringFormatter: Expected a string object for the format string but observed: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Make sure that the argument slice is a reference array.
	field := params[1].(*object.Object).FieldTable["value"]
	if !strings.HasPrefix(field.Ftype, types.RefArray) {
		errMsg := fmt.Sprintf("StringFormatter: Expected Ftype=%s for params[1]: fld.Ftype=%s, fld.Fvalue=%v",
			types.RefArray, field.Ftype, field.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// valuesIn = the reference array
	valuesIn := field.Fvalue.([]*object.Object)

	// Convert input arguments but keep unknown refs for later handling
	rawArgs := make([]interface{}, 0, len(valuesIn))
	for ii := 0; ii < len(valuesIn); ii++ {
		obj := valuesIn[ii]
		if obj == nil || object.IsNull(obj) {
			rawArgs = append(rawArgs, nil)
			continue
		}
		fld := obj.FieldTable["value"]
		if fld.Ftype == types.StringClassRef {
			var str string
			switch fld.Fvalue.(type) {
			case []byte:
				str = string(fld.Fvalue.([]byte))
			case []types.JavaByte:
				str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
			}
			rawArgs = append(rawArgs, str)
			continue
		}
		switch fld.Ftype {
		case types.ByteArray:
			var str string
			switch fld.Fvalue.(type) {
			case []byte:
				str = string(fld.Fvalue.([]byte))
			case []types.JavaByte:
				str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
			}
			rawArgs = append(rawArgs, str)
		case types.Bool:
			var zz bool
			if fld.Fvalue.(int64) == 0 {
				zz = false
			} else {
				zz = true
			}
			rawArgs = append(rawArgs, zz)
		case types.Char:
			rooney := rune(fld.Fvalue.(int64))
			rawArgs = append(rawArgs, rooney)
		case types.Double:
			rawArgs = append(rawArgs, fld.Fvalue.(float64))
		case types.Float:
			rawArgs = append(rawArgs, fld.Fvalue.(float64))
		case types.Int:
			rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 32})
		case types.Long:
			rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 64})
		case types.Short:
			rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 16})
		case types.Byte:
			rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 8})
		default:
			// keep the full object for later processing (e.g., %s/%b/%h)
			rawArgs = append(rawArgs, obj)
		}
	}

	// Transform Java format string and arguments to Go compatible ones
	newFmt, newArgs := translateJavaFormat(formatString, rawArgs)

	str := fmt.Sprintf(newFmt, newArgs...)
	return object.StringObjectFromGoString(str)
}

// translateJavaFormat parses a Java String.format-style format and maps it to a Go fmt format,
// returning the rewritten format string and adjusted argument slice.
func translateJavaFormat(fmtJava string, rawArgs []interface{}) (string, []interface{}) {
	var b strings.Builder
	outArgs := make([]interface{}, 0)

	nextIndex := 0
	lastIndex := -1
	newline := "\n"
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}

	for i := 0; i < len(fmtJava); i++ {
		ch := fmtJava[i]
		if ch != '%' {
			b.WriteByte(ch)
			continue
		}
		// handle %% literal
		if i+1 < len(fmtJava) && fmtJava[i+1] == '%' {
			b.WriteString("%%")
			i++
			continue
		}
		j := i + 1
		// parse argument_index (digits+$) but only accept if a trailing '$' is present
		argIndex := -1
		// tentatively scan digits
		tmp := j
		for tmp < len(fmtJava) && fmtJava[tmp] >= '0' && fmtJava[tmp] <= '9' {
			tmp++
		}
		if tmp < len(fmtJava) && fmtJava[tmp] == '$' && tmp > i+1 {
			idxStr := fmtJava[i+1 : tmp]
			if v, err := strconv.Atoi(idxStr); err == nil && v > 0 {
				argIndex = v - 1
			}
			j = tmp + 1
		} else {
			// no argument index; keep j right after '%'
			j = i + 1
		}
		// parse flags (we will collect but drop '<')
		flagsStart := j
		reusePrev := false
		for j < len(fmtJava) {
			c := fmtJava[j]
			if strings.ContainsRune("-#+ 0,(<", rune(c)) {
				if c == '<' {
					reusePrev = true
				}
				j++
			} else {
				break
			}
		}
		flags := fmtJava[flagsStart:j]
		flags = strings.ReplaceAll(flags, "<", "") // we'll resolve it ourselves

		// width
		widthStart := j
		for j < len(fmtJava) && fmtJava[j] >= '0' && fmtJava[j] <= '9' {
			j++
		}
		width := fmtJava[widthStart:j]
		// precision
		precision := ""
		if j < len(fmtJava) && fmtJava[j] == '.' {
			k := j + 1
			for k < len(fmtJava) && fmtJava[k] >= '0' && fmtJava[k] <= '9' {
				k++
			}
			precision = fmtJava[j:k]
			j = k
		}
		if j >= len(fmtJava) {
			// malformed, copy as-is
			b.WriteString(fmtJava[i:])
			break
		}
		conv := fmtJava[j]

		// %n newline (no arg consumed)
		if conv == 'n' {
			b.WriteString(newline)
			i = j
			continue
		}

		// Resolve which argument index to use
		useIndex := -1
		if argIndex >= 0 {
			useIndex = argIndex
		} else if reusePrev {
			useIndex = lastIndex
		} else {
			useIndex = nextIndex
			nextIndex++
		}
		lastIndex = useIndex

		// Prepare Go conv and value
		goConv := string(conv)
		switch conv {
		case 'b', 'B':
			goConv = "t"
			val := coerceJavaBoolean(rawArgs, useIndex)
			outArgs = append(outArgs, val)
		case 's', 'S':
			goConv = "s"
			val := coerceJavaString(rawArgs, useIndex)
			if conv == 'S' {
				val = strings.ToUpper(val)
			}
			outArgs = append(outArgs, val)
		case 'h', 'H':
			// Java %h/%H produces hex hash of arg (null -> "null")
			if rawArgs == nil || useIndex < 0 || useIndex >= len(rawArgs) || rawArgs[useIndex] == nil {
				goConv = "s"
				outArgs = append(outArgs, "null")
			} else {
				goConv = string(conv) // x or X
				if conv == 'h' {
					goConv = "x"
				} else {
					goConv = "X"
				}
				outArgs = append(outArgs, javaHashValue(rawArgs[useIndex]))
			}
		case 't', 'T':
			// Not fully supported: degrade to %s on placeholder
			goConv = "s"
			val := coerceJavaString(rawArgs, useIndex)
			outArgs = append(outArgs, val)
		case 'x', 'X', 'o':
			// For hex/octal, Java uses two's complement unsigned representation of the primitive width
			v := rawArgs[useIndex]
			switch iv := v.(type) {
			case intWithBits:
				var u uint64
				switch iv.bits {
				case 64:
					u = uint64(iv.v)
				case 32:
					u = uint64(uint32(iv.v))
				case 16:
					u = uint64(uint16(iv.v))
				case 8:
					u = uint64(uint8(iv.v))
				default:
					u = uint64(uint32(iv.v))
				}
				outArgs = append(outArgs, u)
			default:
				outArgs = append(outArgs, normalizeForGo(rawArgs, useIndex))
			}
		case 'f':
			// Let Go handle width/precision/flags for fixed-point
			goConv = "f"
			outArgs = append(outArgs, normalizeForGo(rawArgs, useIndex))
		default:
			// Pass-through, but ensure we supply the raw argument in Go-native type
			outArgs = append(outArgs, normalizeForGo(rawArgs, useIndex))
		}

		// Rebuild the format specifier without index or '<'
		b.WriteByte('%')
		b.WriteString(flags)
		b.WriteString(width)
		b.WriteString(precision)
		b.WriteString(goConv)

		i = j
	}
	return b.String(), outArgs
}

func coerceJavaBoolean(args []interface{}, idx int) bool {
	if args == nil || idx < 0 || idx >= len(args) {
		return false
	}
	v := args[idx]
	switch vv := v.(type) {
	case bool:
		return vv
	case *object.Object:
		if vv == nil || object.IsNull(vv) {
			return false
		}
		return true
	default:
		// primitives and strings are considered non-null in Java formatting
		return true
	}
}

func coerceJavaString(args []interface{}, idx int) string {
	if args == nil || idx < 0 || idx >= len(args) {
		return "null"
	}
	v := args[idx]
	switch vv := v.(type) {
	case string:
		return vv
	case *object.Object:
		if vv == nil || object.IsNull(vv) {
			return "null"
		}
		// For non-String ref, mimic minimal Object.toString(): ClassName@hex
		klass := object.GetClassNameSuffix(vv, true)
		addr := fmt.Sprintf("%p", vv) // e.g., 0x1234abcd
		if strings.HasPrefix(addr, "0x") {
			addr = addr[2:]
		}
		return fmt.Sprintf("%s@%s", klass, addr)
	default:
		return fmt.Sprintf("%v", vv)
	}
}

func normalizeForGo(args []interface{}, idx int) interface{} {
	if args == nil || idx < 0 || idx >= len(args) {
		return nil
	}
	v := args[idx]
	switch vv := v.(type) {
	case *object.Object:
		// For unknown object refs, format as their toString-like string
		return coerceJavaString(args, idx)
	case intWithBits:
		return vv.v
	default:
		return vv
	}
}

func javaHashValue(v interface{}) uint64 {
	switch vv := v.(type) {
	case nil:
		return 0
	case bool:
		if vv {
			return 1231
		}
		return 1237
	case int64:
		return uint64(vv)
	case float64:
		// use IEEE bits as basis
		return uint64(mathFloat64bits(vv))
	case rune:
		return uint64(vv)
	case string:
		return uint64(javaStringHashCode(vv))
	case *object.Object:
		if vv == nil || object.IsNull(vv) {
			return 0
		}
		// use pointer address as identity hash surrogate
		addrStr := fmt.Sprintf("%p", vv)
		if strings.HasPrefix(addrStr, "0x") {
			addrStr = addrStr[2:]
		}
		ui, _ := strconv.ParseUint(addrStr, 16, 64)
		return ui
	default:
		return 0
	}
}

func javaStringHashCode(s string) int32 {
	var h int32 = 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int32(s[i])
	}
	return h
}

// mathFloat64bits avoids importing math for a single call
func mathFloat64bits(f float64) uint64 {
	return math.Float64bits(f)
}
