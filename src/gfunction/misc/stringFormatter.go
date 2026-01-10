package misc

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaMath"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// helper wrapper to keep Java integral bit-width for formatting
// bits=32 for int/short/byte, bits=64 for long
type intWithBits struct {
	v    int64
	bits int
}

// wrapper to retain BigDecimal identity for hashing while providing float for %f
type bigDecimalArg struct {
	obj *object.Object
	f64 float64
}

// wrapper for Float to compute Java Float.hashCode correctly and still print as number
type float32Arg struct {
	f32 float32
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
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
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	default:
		errMsg := fmt.Sprintf("StringFormatter: Expected a string object for the format string but observed: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Make sure that the argument slice is a reference array.
	field := params[1].(*object.Object).FieldTable["value"]
	if !strings.HasPrefix(field.Ftype, types.RefArray) {
		errMsg := fmt.Sprintf("StringFormatter: Expected Ftype=%s for params[1]: fld.Ftype=%s, fld.Fvalue=%v",
			types.RefArray, field.Ftype, field.Fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
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
		// Prefer handling by presence of a "value" field; fall back to BigInteger/BigDecimal detection.
		fld, hasValue := obj.FieldTable["value"]
		if hasValue && fld.Ftype == types.StringClassRef {
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
		if hasValue {
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
				// Preserve float32 identity for correct %h hashing; store as wrapper
				f32 := float32(fld.Fvalue.(float64))
				rawArgs = append(rawArgs, float32Arg{f32: f32})
			case types.Int:
				rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 32})
			case types.Long:
				rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 64})
			case types.Short:
				rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 16})
			case types.Byte:
				rawArgs = append(rawArgs, intWithBits{v: fld.Fvalue.(int64), bits: 8})
			case types.BigInteger:
				// Underlying Go value is *big.Int in the value field
				if bi, ok := fld.Fvalue.(*big.Int); ok {
					rawArgs = append(rawArgs, bi)
				} else {
					rawArgs = append(rawArgs, obj)
				}
			default:
				// keep the full object for later processing (e.g., %s/%b/%h)
				rawArgs = append(rawArgs, obj)
			}
			continue
		}
		// No "value" field: detect BigDecimal by fields intVal/scale
		if ivFld, ok := obj.FieldTable["intVal"]; ok {
			if scFld, ok2 := obj.FieldTable["scale"]; ok2 {
				intValObj, ok3 := ivFld.Fvalue.(*object.Object)
				if ok3 {
					if vFld, ok4 := intValObj.FieldTable["value"]; ok4 {
						if unscaled, ok5 := vFld.Fvalue.(*big.Int); ok5 {
							scale, _ := scFld.Fvalue.(int64)
							bf := new(big.Float).SetInt(unscaled)
							if scale != 0 {
								den := new(big.Float).SetFloat64(math.Pow10(int(scale)))
								bf.Quo(bf, den)
							}
							f64, _ := bf.Float64()
							// Keep identity for hashing; provide f64 for %f via wrapper
							rawArgs = append(rawArgs, bigDecimalArg{obj: obj, f64: f64})
							continue
						}
					}
				}
			}
		}
		// Fallback: keep object reference
		rawArgs = append(rawArgs, obj)
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
			// Java %b/%B: boolean formatted as true/false; %B must be uppercased
			val := coerceJavaBoolean(rawArgs, useIndex)
			str := "false"
			if val {
				str = "true"
			}
			if conv == 'B' {
				str = strings.ToUpper(str)
			}
			goConv = "s"
			outArgs = append(outArgs, str)
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
		case 'c', 'C':
			// Java %c/%C: character; %C is uppercased variant
			// Determine rune code point from various input types
			var r rune
			if rawArgs == nil || useIndex < 0 || useIndex >= len(rawArgs) || rawArgs[useIndex] == nil {
				// Java would throw on null for %c; degrade to 0 rune
				r = 0
			} else {
				v := rawArgs[useIndex]
				switch vv := v.(type) {
				case rune:
					r = vv
				case intWithBits:
					// Treat as Unicode code point (like Java Formatter)
					r = rune(uint32(vv.v))
				case int64:
					r = rune(uint32(vv))
				case string:
					// If a string sneaks in, take first rune
					for _, rr := range vv {
						r = rr
						break
					}
				default:
					// Fallback via fmt to string then first rune
					s := coerceJavaString(rawArgs, useIndex)
					for _, rr := range s {
						r = rr
						break
					}
				}
			}
			if conv == 'C' {
				// Uppercase the code point
				r = unicode.ToUpper(r)
			}
			goConv = "c"
			outArgs = append(outArgs, r)
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
			if bd, ok := rawArgs[useIndex].(bigDecimalArg); ok {
				outArgs = append(outArgs, bd.f64)
			} else if f32a, ok := rawArgs[useIndex].(float32Arg); ok {
				outArgs = append(outArgs, float64(f32a.f32))
			} else {
				outArgs = append(outArgs, normalizeForGo(rawArgs, useIndex))
			}
		case 'e', 'E':
			// Scientific notation; support BigDecimal and Float by passing float64 value
			if conv == 'E' {
				goConv = "E"
			} else {
				goConv = "e"
			}
			if bd, ok := rawArgs[useIndex].(bigDecimalArg); ok {
				outArgs = append(outArgs, bd.f64)
			} else if f32a, ok := rawArgs[useIndex].(float32Arg); ok {
				outArgs = append(outArgs, float64(f32a.f32))
			} else {
				outArgs = append(outArgs, normalizeForGo(rawArgs, useIndex))
			}
		case 'g', 'G':
			// General format; support BigDecimal and Float similarly
			if conv == 'G' {
				goConv = "G"
			} else {
				goConv = "g"
			}
			if bd, ok := rawArgs[useIndex].(bigDecimalArg); ok {
				outArgs = append(outArgs, bd.f64)
			} else if f32a, ok := rawArgs[useIndex].(float32Arg); ok {
				outArgs = append(outArgs, float64(f32a.f32))
			} else {
				outArgs = append(outArgs, normalizeForGo(rawArgs, useIndex))
			}
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
	if v == nil {
		return false
	}
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
		// Special-case BigInteger and BigDecimal to return their numeric string
		if valFld, ok := vv.FieldTable["value"]; ok {
			if valFld.Ftype == types.BigInteger {
				if bi, ok2 := valFld.Fvalue.(*big.Int); ok2 {
					return bi.String()
				}
			}
		}
		if iv, ok := vv.FieldTable["intVal"]; ok {
			if sc, ok2 := vv.FieldTable["scale"]; ok2 {
				if biObj, ok3 := iv.Fvalue.(*object.Object); ok3 {
					if vFld, ok4 := biObj.FieldTable["value"]; ok4 {
						if unscaled, ok5 := vFld.Fvalue.(*big.Int); ok5 {
							scale, _ := sc.Fvalue.(int64)
							return javaMath.FormatDecimalString(unscaled, scale)
						}
					}
				}
			}
		}
		// For other refs, mimic minimal Object.toString(): ClassName@hex
		klass := object.GetClassNameSuffix(vv, true)
		addr := fmt.Sprintf("%p", vv) // e.g., 0x1234abcd
		if strings.HasPrefix(addr, "0x") {
			addr = addr[2:]
		}
		return fmt.Sprintf("%s@%s", klass, addr)
	default:
		// Unwrap BigDecimal proxy for string coercion
		if bd, ok := vv.(bigDecimalArg); ok {
			// Delegate to object-based path by crafting decimal string
			if bd.obj != nil {
				if iv, ok := bd.obj.FieldTable["intVal"]; ok {
					if sc, ok2 := bd.obj.FieldTable["scale"]; ok2 {
						if biObj, ok3 := iv.Fvalue.(*object.Object); ok3 {
							if vFld, ok4 := biObj.FieldTable["value"]; ok4 {
								if unscaled, ok5 := vFld.Fvalue.(*big.Int); ok5 {
									scale, _ := sc.Fvalue.(int64)
									return javaMath.FormatDecimalString(unscaled, scale)
								}
							}
						}
					}
				}
			}
		}
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
	// Return a 32-bit Java-style hash value (as uint64 for Go's fmt),
	// matching Java Formatter's %h semantics which uses Wrapper.hashCode().
	switch vv := v.(type) {
	case nil:
		return 0
	case bool:
		if vv {
			return 1231
		}
		return 1237
	case intWithBits:
		// Emulate Java primitive wrapper hashCode()
		if vv.bits == 64 {
			u := uint64(vv.v)
			h := uint32((u >> 32) ^ (u & 0xffffffff))
			return uint64(h)
		}
		// 32/16/8-bit signed values, sign-extended to 32-bit
		var i32 int32
		switch vv.bits {
		case 32:
			i32 = int32(vv.v)
		case 16:
			i32 = int32(int16(vv.v))
		case 8:
			i32 = int32(int8(vv.v))
		default:
			i32 = int32(vv.v)
		}
		return uint64(uint32(i32))
	case int64:
		// Treat as long
		u := uint64(vv)
		h := uint32((u >> 32) ^ (u & 0xffffffff))
		return uint64(h)
	case float64:
		// Java Double.hashCode: int(bits ^ (bits >>> 32)) with canonical NaN
		bits := mathFloat64bits(vv)
		if math.IsNaN(vv) {
			bits = 0x7ff8000000000000
		}
		h := uint32(uint32(bits>>32) ^ uint32(bits))
		return uint64(h)
	case float32Arg:
		// Java Float.hashCode: int(Float.floatToIntBits(f)) with canonical NaN 0x7fc00000
		f := float64(vv.f32)
		bits := uint32(math.Float32bits(vv.f32))
		if math.IsNaN(f) {
			bits = 0x7fc00000
		}
		return uint64(bits)
	case rune:
		return uint64(uint32(vv))
	case string:
		return uint64(uint32(javaStringHashCode(vv)))
	case *big.Int:
		if vv == nil {
			return 0
		}
		// Exact JDK BigInteger.hashCode:
		// int hash = 0; for each 32-bit word of magnitude (big-endian unsigned): hash = 31*hash + (word & 0xffffffff);
		// then: hash *= signum;
		abs := new(big.Int).Abs(vv)
		bytes := abs.Bytes() // big-endian magnitude without sign
		if len(bytes) == 0 {
			return 0
		}
		// left-pad to multiple of 4 bytes
		pad := (4 - (len(bytes) % 4)) % 4
		if pad != 0 {
			padded := make([]byte, pad+len(bytes))
			copy(padded[pad:], bytes)
			bytes = padded
		}
		var h int32 = 0
		for i := 0; i < len(bytes); i += 4 {
			word := uint32(bytes[i])<<24 | uint32(bytes[i+1])<<16 | uint32(bytes[i+2])<<8 | uint32(bytes[i+3])
			h = int32(31)*h + int32(word)
		}
		// apply signum
		if vv.Sign() < 0 {
			h = -h
		}
		return uint64(uint32(h))
	case bigDecimalArg:
		// BigDecimal.hashCode() per JDK: int h = 31*intVal.hashCode() + scale
		if vv.obj == nil {
			return 0
		}
		iv, ok1 := vv.obj.FieldTable["intVal"]
		sc, ok2 := vv.obj.FieldTable["scale"]
		if !ok1 || !ok2 {
			return 0
		}
		biObj, ok3 := iv.Fvalue.(*object.Object)
		if !ok3 || biObj == nil {
			return 0
		}
		vFld, ok4 := biObj.FieldTable["value"]
		if !ok4 {
			return 0
		}
		bi, ok5 := vFld.Fvalue.(*big.Int)
		if !ok5 {
			return 0
		}
		biHash := int32(javaHashValue(bi))
		scale32 := int32(sc.Fvalue.(int64))
		h := int32(31)*biHash + scale32
		return uint64(uint32(h))
	case *object.Object:
		if vv == nil || object.IsNull(vv) {
			return 0
		}
		// Per Java spec, %h uses arg.hashCode(); for default Object, our implementation is ptr^(ptr>>32)
		ptr := uintptr(unsafe.Pointer(vv))
		h := uint32(ptr ^ (ptr >> 32))
		return uint64(h)
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
