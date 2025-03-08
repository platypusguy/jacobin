/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/config"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/thread"
	"jacobin/trace"
	"jacobin/types"
	"jacobin/util"
	"os"
	"runtime/debug"
	"strconv"
)

var MainThread thread.ExecThread

// StartExec is where execution begins. It initializes various structures, such as
// the MTable, then using the passed-in name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string, mainThread *thread.ExecThread, globalStruct *globals.Globals) {

	MainThread = *mainThread

	me, err := classloader.FetchMethodAndCP(className, "main", "([Ljava/lang/String;)V")
	if err != nil {
		errMsg := "Class not found: " + className + ".main()"
		exceptions.ThrowEx(excNames.ClassNotFoundException, errMsg, nil)
	}

	m := me.Meth.(classloader.JmEntry)
	f := frames.CreateFrame(m.MaxStack + types.StackInflator) // experiment with stack size. See JACOBIN-494
	f.Thread = MainThread.ID
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"
	f.ClName = className
	f.CP = m.Cp                        // add its pointer to the class CP
	f.Meth = append(f.Meth, m.Code...) // copy the bytecodes over

	// allocate the local variables
	for k := 0; k < m.MaxLocals; k++ {
		f.Locals = append(f.Locals, 0)
	}

	// Create an array of string objects in locals[0].
	var objArray []*object.Object
	for _, str := range globalStruct.AppArgs {
		// sobj := object.NewStringFromGoString(str) // deprecated by JACOBIN-480
		sobj := object.StringObjectFromGoString(str)
		objArray = append(objArray, sobj)
	}
	f.Locals[0] = object.MakePrimitiveObject("[Ljava/lang/String", types.RefArray, objArray)

	// create the first thread and place its first frame on it
	MainThread.Stack = frames.CreateFrameStack()
	mainThread.Stack = MainThread.Stack

	// moved here as part of JACOBIN-554. Was previously after the InstantiateClass() call next
	if frames.PushFrame(MainThread.Stack, f) != nil {
		errMsg := "Memory error allocating frame on thread: " + strconv.Itoa(MainThread.ID)
		exceptions.ThrowEx(excNames.OutOfMemoryError, errMsg, nil)
	}

	// must first instantiate the class, so that any static initializers are run
	_, instantiateError := InstantiateClass(className, MainThread.Stack)
	if instantiateError != nil {
		errMsg := "Error instantiating: " + className + ".main()"
		exceptions.ThrowEx(excNames.InstantiationException, errMsg, nil)
	}

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("StartExec: class=%s, meth=%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, m.MaxStack, m.MaxLocals, len(m.Code))
		trace.Trace(traceInfo)
	}

	err = runThread(&MainThread)

	if globals.TraceVerbose {
		statics.DumpStatics("StartExec end", statics.SelectUser, "")
		_ = config.DumpConfig(os.Stderr)
	}
}

// Point the thread to the top of the frame stack and tell it to run from there.
func runThread(t *thread.ExecThread) error {

	defer func() int {
		// only an untrapped panic gets us here
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack
			exceptions.ShowPanicCause(r)
			exceptions.ShowFrameStack(t)
			exceptions.ShowGoStackTrace(nil)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	for t.Stack.Len() > 0 {
		interpret(t.Stack)
	}

	if t.Stack.Len() == 0 { // true when the last executed frame was main()
		return nil
	}
	return nil
}

/*
// runFrame() is the principal execution function in Jacobin. It first tests for a
// golang function in the present frame. If it is a golang function, it's sent to
// a different function for execution. Otherwise, bytecode interpretation takes
// place through a giant switch statement.
func runFrame(fs *list.List) error {
	glob := globals.GetGlobalRef()

frameInterpreter:
	// the current frame is always the head of the linked list of frames.
	// the next statement converts the address of that frame to the more readable 'f'
	f := fs.Front().Value.(*frames.Frame)
	f.WideInEffect = false

	// the frame's method is not a golang method, so it's Java bytecode, which
	// is interpreted in the rest of this function.
	for f.PC < len(f.Meth) {
		if globals.TraceInst {
			traceInfo := emitTraceData(f)
			trace.Trace(traceInfo)
		}

		opcode := f.Meth[f.PC]
		f.ExceptionPC = f.PC // in the event of an exception, here's where we were
		switch opcode {      // cases listed in numerical value of opcode

		case opcodes.INVOKEINTERFACE: // 0xB9 invoke an interface
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			count := f.Meth[f.PC+3]
			zeroByte := f.Meth[f.PC+4]
			f.PC += 4

			CP := f.CP.(*classloader.CPool)
			if count < 1 || CPslot >= len(CP.CpIndex) || zeroByte != 0x00 {
				errMsg := fmt.Sprintf("Invalid values for INVOKEINTERFACE bytecode")
				status := exceptions.ThrowEx(excNames.IllegalClassFormatException, errMsg, f)
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type != classloader.Interface {
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("INVOKEINTERFACE: CP entry type (%d) did not point to an interface method type (%d)",
					CPentry.Type, classloader.Interface)
				status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, f) // this is the error thrown by JDK
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			method := CP.InterfaceRefs[CPentry.Slot]

			// get the class entry from this method
			interfaceRef := method.ClassIndex
			interfaceNameIndex := CP.ClassRefs[CP.CpIndex[interfaceRef].Slot]
			interfaceNamePtr := stringPool.GetStringPointer(interfaceNameIndex)
			interfaceName := *interfaceNamePtr

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := CP.NameAndTypes[nAndTslot]
			interfaceMethodNameIndex := nAndT.NameIndex
			interfaceMethodName := classloader.FetchUTF8stringFromCPEntryNumber(CP, interfaceMethodNameIndex)

			// get the signature for this method
			interfaceMethodSigIndex := nAndT.DescIndex
			interfaceMethodType := classloader.FetchUTF8stringFromCPEntryNumber(
				CP, interfaceMethodSigIndex)

			// now get the objRef pointing to the class containing the call to the method
			// described just previously. It is located on the f.OpStack below the args to
			// be passed to the method.
			// The objRef object has previously been instantiated and its constructor called.
			objRef := f.OpStack[f.TOS-int(count)+1]
			if objRef == nil {
				errMsg := fmt.Sprintf("INVOKEINTERFACE: object whose method, %s, is invoked is null",
					interfaceName+interfaceMethodName+interfaceMethodType)
				status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, f)
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			// get the name of the objectRef's class, and make sure it's loaded
			objRefClassName := *(stringPool.GetStringPointer(objRef.(*object.Object).KlassName))
			if err := classloader.LoadClassFromNameOnly(objRefClassName); err != nil {
				// in this case, LoadClassFromNameOnly() will have already thrown the exception
				if globals.JacobinHome() == "test" {
					return err // applies only if in test
				}
			}

			class := classloader.MethAreaFetch(objRefClassName)
			if class == nil {
				// in theory, this can't happen due to immediately previous loading, but making sure
				errMsg := fmt.Sprintf("INVOKEINTERFACE: class %s not found", objRefClassName)
				status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, f)
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			var mtEntry classloader.MTentry
			var err error
			mtEntry, err = locateInterfaceMeth(class, f, objRefClassName, interfaceName,
				interfaceMethodName, interfaceMethodType)
			if err != nil { // any error will already have been handled
				continue
			}

			clData := *class.Data
			if mtEntry.MType == 'J' {
				entry := mtEntry.Meth.(classloader.JmEntry)
				fram, err := createAndInitNewFrame(
					clData.Name, interfaceMethodName, interfaceMethodType, &entry, true, f)
				if err != nil {
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := "INVOKEINTERFACE: Error creating frame in: " + clData.Name + "." +
						interfaceMethodName + interfaceMethodType
					status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, f)
					if status != exceptions.Caught {
						return errors.New(errMsg) // applies only if in test
					}
				}

				f.PC += 1                            // to point to the next bytecode before exiting
				fs.PushFront(fram)                   // push the new frame
				f = fs.Front().Value.(*frames.Frame) // point f to the new head
				goto frameInterpreter
			} else if mtEntry.MType == 'G' { // it's a gfunction (i.e., a native function implemented in golang)
				gmethData := mtEntry.Meth.(gfunction.GMeth)
				paramCount := gmethData.ParamSlots
				var params []interface{}
				for i := 0; i < paramCount; i++ {
					params = append(params, pop(f))
				}

				if globals.TraceInst {
					infoMsg := fmt.Sprintf("G-function: interface=%s, meth=%s%s", interfaceName, interfaceName, interfaceMethodType)
					trace.Trace(infoMsg)
				}
				ret := gfunction.RunGfunction(mtEntry, fs, interfaceName, interfaceMethodName, interfaceMethodType, &params, true, globals.TraceVerbose)
				if ret != nil {
					switch ret.(type) {
					case error:
						if glob.JacobinName == "test" {
							errRet := ret.(error)
							return errRet
						} else if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
							f.PC += 1
							goto frameInterpreter
						}
					default: // if it's not an error, then it's a legitimate return value, which we simply push
						push(f, ret)
						if strings.HasSuffix(interfaceMethodType, "D") || strings.HasSuffix(interfaceMethodType, "J") {
							push(f, ret) // push twice if long or double
						}
					}
					// any exception will already have been handled.
				}
			}


		case opcodes.CHECKCAST: // 0xC0 same as INSTANCEOF but does nothing on null,
			// and doesn't change the stack if the cast is legal.
			// Because this uses the same logic as INSTANCEOF, any change here should
			// be made to INSTANCEOF

			ref := peek(f)  // peek b/c the objectRef is *not* removed from the op stack
			if ref == nil { // if ref is nil, just carry on
				f.PC += 2 // move past two bytes pointing to comp object
				f.PC += 1
				continue // cannot goto checkcastOK, b/c golang doesn't allow a jump over variable initialization
			}

			var obj *object.Object
			var objName string
			switch ref.(type) {
			case *object.Object:
				if object.IsNull(ref) { // if ref is null, just carry on
					f.PC += 2 // move past two bytes pointing to comp object
					f.PC += 1
					continue
				} else {
					obj = (ref).(*object.Object)
					objName = *(stringPool.GetStringPointer(obj.KlassName))
				}
			default: // objectRef must be a reference to an object
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("CHECKCAST: Invalid class reference, type=%T", ref)
				status := exceptions.ThrowEx(excNames.ClassCastException, errMsg, f)
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			// at this point, we know we have a non-nil, non-null pointer to an object;
			// now, get the class we're casting the object to.
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			// CPentry := CP.CpIndex[CPslot]
			classNamePtr := classloader.FetchCPentry(CP, CPslot)

			var objClassType = types.Error
			if strings.HasPrefix(objName, "[") {
				objClassType = types.Array
			} else {
				objData := classloader.MethAreaFetch(objName)
				if objData == nil || objData.Data == nil {
					_ = classloader.LoadClassFromNameOnly(objName)
					objData = classloader.MethAreaFetch(objName)
				}
				if objData.Data.Access.ClassIsInterface {
					objClassType = types.Interface
				} else {
					objClassType = types.NonArrayObject
				}
			}

			var checkcastStatus bool
			switch objClassType {
			case types.NonArrayObject:
				checkcastStatus = checkcastNonArrayObject(obj, *(classNamePtr.StringVal))
			case types.Array:
				checkcastStatus = checkcastArray(obj, *(classNamePtr.StringVal))
			case types.Interface:
				checkcastStatus = checkcastInterface(obj, *(classNamePtr.StringVal))
			default:
				errMsg := fmt.Sprintf("CHECKCAST: expected to verify class or interface, but got none")
				status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, f)
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			if checkcastStatus == false {
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("CHECKCAST: %s is not castable with respect to %s",
					*(stringPool.GetStringPointer(obj.KlassName)), *(classNamePtr.StringVal))
				status := exceptions.ThrowEx(excNames.ClassCastException, errMsg, f)
				if status != exceptions.Caught {
					return errors.New(errMsg) // applies only if in test
				}
			}

			// if it is castable, do nothing.
			/* // TODO CODE to review for use in runUtils.go
			if CPentry.Type == classloader.ClassRef {
				// slot of ClassRef points to a CP entry for a UTF8 record w/ name of class
				var className string
				classNamePtr = classloader.FetchCPentry(CP, CPslot)
				if classNamePtr.RetType != classloader.IS_STRING_ADDR {
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("CHECKCAST: Invalid classRef found, classNamePtr.RetType=%d", classNamePtr.RetType)
					trace.Error(errMsg)
					return errors.New(errMsg)
				} else {
					errMsg := fmt.Sprintf("CHECKCAST: expected to verify class or interface, but got none")
					status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, f)
					if status != exceptions.Caught {
						return errors.New(errMsg) // applies only if in test
					}
				}

				// we now know we point to a valid class, array, or interface. We handle classes and arrays here.
				className = *(classNamePtr.StringVal)
				if globals.TraceVerbose {
					var traceInfo string
					if strings.HasPrefix(className, "[") {
						traceInfo = fmt.Sprintf("CHECKCAST: class is an array = %s", className)
					} else {
						traceInfo = fmt.Sprintf("CHECKCAST: className = %s", className)
					}
					trace.Trace(traceInfo)
				}

			e now have the resolved class (className) and the objectref (obj)
			    The rules for identifying obj can be cast to classname are (from the JVM 17 spec):

				If objectref can be cast to the resolved class, array, or interface type, the operand stack is
			    unchanged; otherwise, the checkcast instruction throws a ClassCastException.

				S = obj
				T = className

				If S is the type of the object referred to by objectref, and T is the resolved class, array, or
				interface type, then checkcast determines whether objectref can be cast to type T as follows:

				If S is a class type, then:
				* If T is a class type, then S must be the same class as T, or S must be a subclass of T;
				* If T is an interface type, then S must implement interface T.

				If S is an array type SC[], that is, an array of components of type SC, then:
				* If T is a class type, then T must be Object.
				* If T is an interface type, then T must be one of the interfaces implemented by arrays (JLS §4.10.3).
				* If T is an array type TC[], that is, an array of components of type TC, then one of the following
				  must be true:
					> TC and SC are the same primitive type.
					> TC and SC are reference types, and type SC can be cast to TC by
				      recursive application of these rules.


		case opcodes.INSTANCEOF: // 0xC1 validate the type of object (if not nil or null)
			// because this uses similar logic to CHECKCAST, any change here should
			// likely be made to CHECKCAST as well
			ref := pop(f)
			if ref == nil || ref == object.Null {
				push(f, int64(0))
				f.PC += 2 // move past index bytes to comp object
				break
			}

			switch ref.(type) {
			case *object.Object:
				if ref == object.Null {
					push(f, int64(0))
					f.PC += 2 // move past two bytes pointing to comp object
					break
				} else {
					obj := *ref.(*object.Object)
					CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
					f.PC += 2
					CP := f.CP.(*classloader.CPool)
					CPentry := CP.CpIndex[CPslot]
					if CPentry.Type == classloader.ClassRef { // slot of ClassRef points to
						// a CP entry for a stringPool entry for name of class
						var className string
						classNamePtr := classloader.FetchCPentry(CP, CPslot)
						if classNamePtr.RetType != classloader.IS_STRING_ADDR {
							glob.ErrorGoStack = string(debug.Stack())
							errMsg := "INSTANCEOF: Invalid classRef found"
							trace.Error(errMsg)
							return errors.New(errMsg)
						} else {
							className = *(classNamePtr.StringVal)
							if globals.TraceVerbose {
								traceInfo := fmt.Sprintf("INSTANCEOF: className = %s", className)
								trace.Trace(traceInfo)
							}
						}
						classPtr := classloader.MethAreaFetch(className)
						if classPtr == nil { // class wasn't loaded, so load it now
							if classloader.LoadClassFromNameOnly(className) != nil {
								glob.ErrorGoStack = string(debug.Stack())
								errMsg := "INSTANCEOF: Could not load class: " + className
								trace.Error(errMsg)
								return errors.New(errMsg)
							}
							classPtr = classloader.MethAreaFetch(className)
						}
						if classPtr == classloader.MethAreaFetch(*(stringPool.GetStringPointer(obj.KlassName))) {
							push(f, int64(1))
						} else {
							push(f, int64(0))
						}
					}
				}
			}


		f.PC += 1
	}
	return nil
}
*/

// multiply two numbers
func multiply[N frames.Number](num1, num2 N) N {
	return num1 * num2
}

func subtract[N frames.Number](num1, num2 N) N {
	return num1 - num2
}

// create a new frame and load up the local variables with the passed
// arguments, set up the stack, and all the remaining items to begin execution
// Note: the includeObjectRef parameter is a boolean. When true, it indicates
// that in addition to the method parameter, an object reference is also on
// the stack and needs to be popped off the caller's opStack and passed in.
// (This would be the case for invokevirtual, among others.) When false, no
// object pointer is needed (for invokestatic, among others).
func createAndInitNewFrame(
	className string, methodName string, methodType string,
	m *classloader.JmEntry,
	includeObjectRef bool,
	currFrame *frames.Frame) (*frames.Frame, error) {

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("createAndInitNewFrame: class=%s, meth=%s%s, includeObjectRef=%v, maxStack=%d, maxLocals=%d",
			className, methodName, methodType, includeObjectRef, m.MaxStack, m.MaxLocals)
		trace.Trace(traceInfo)
	}

	f := currFrame

	stackSize := m.MaxStack + types.StackInflator // Experimental addition, see JACOBIN-494
	if stackSize < 1 {
		stackSize = 2
	}

	fram := frames.CreateFrame(stackSize)
	fram.Thread = currFrame.Thread
	fram.FrameStack = currFrame.FrameStack
	fram.ClName = className
	fram.MethName = methodName
	fram.MethType = methodType
	fram.CP = m.Cp                           // add its pointer to the class CP
	fram.Meth = append(fram.Meth, m.Code...) // copy the method's bytecodes over

	// pop the parameters off the present stack and put them in
	// the new frame's locals. This is done in reverse order so
	// that the parameters are pushed in the right order to be
	// popped off by the receiving function
	var argList []interface{}
	paramsToPass :=
		util.ParseIncomingParamsFromMethTypeString(methodType)

	// primitives use a single byte/letter, but arrays can be many bytes:
	// a minimum of two (e.g., [I for array of ints). If the array
	// is multidimensional, the bytes will be [[I with one instance
	// of [ for every dimension. In the case of multidimensional
	// arrays, the arrays are always pushed as arrays of references,
	// and we simply mark off the number of [. For single-dimensional
	// arrays, we pass the kind of pointer that applies and mark off
	// a single instance of [
	for j := len(paramsToPass) - 1; j > -1; j-- {
		param := paramsToPass[j]
		primitive := param[0]

		arrayDimensions := 0
		if primitive == '[' {
			i := 0
			for i = 0; i < len(param); i++ {
				if param[i] == '[' {
					arrayDimensions += 1
				} else {
					break
				}
			}
			// param[i] now holds the primitive of the array
			primitive = param[i]
		}

		if arrayDimensions > 1 { // a multidimensional array
			// if the array is multidimensional, then we are
			// passing in a pointer to an array of references
			// to objects (lower arrays) regardless of the
			// lowest level of primitive in the array
			arg := pop(f).(*object.Object)
			argList = append(argList, arg)
			continue
		}

		if arrayDimensions == 1 { // a single-dimension array
			// a bunch of Java functions return raw arrays (like String.toCharArray()), which
			// are not really viewed by the JVM as objects in the full sense of the term. These
			// almost invariably are single-dimension arrays. So we test for these here and
			// return the corresponding object entity.
			value := pop(f)
			arg := object.MakeArrayFromRawArray(value)
			argList = append(argList, arg)
			continue
		}

		switch primitive { // it's not an array
		case 'D': // double
			arg := pop(f).(float64)
			argList = append(argList, arg)
		case 'F': // float
			arg := pop(f).(float64)
			argList = append(argList, arg)
		case 'B', 'C', 'I', 'S': // byte, char, integer, short
			arg := pop(f)
			switch arg.(type) {
			case int: // the arg should be int64, but is occasionally int. Tracking this down.
				arg = int64(arg.(int))
			}
			argList = append(argList, arg)
		case 'J': // long
			arg := pop(f).(int64)
			argList = append(argList, arg)
		case 'L': // pointer/reference
			arg := pop(f) // can't be *Object b/c the arg could be nil, which would panic
			argList = append(argList, arg)
		default:
			arg := pop(f)
			argList = append(argList, arg)
		}
	}

	// Initialize lenLocals = max (m.MaxLocals, len(argList)) but at least 1
	lenArgList := len(argList)
	lenLocals := m.MaxLocals
	if lenArgList > m.MaxLocals {
		lenLocals = lenArgList
	}
	if lenLocals < 1 {
		lenLocals = 1
	}

	// allocate the local variables
	for k := 0; k < lenLocals; k++ {
		fram.Locals = append(fram.Locals, int64(0))
	}

	// if includeObjectRef is true then objectRef != nil.
	// Insert it in the local[0]
	// This is used in invokevirtual, invokespecial, and invokeinterface.
	destLocal := 0
	if includeObjectRef {
		fram.Locals[0] = pop(f)
		fram.Locals = append(fram.Locals, int64(0)) // add the slot taken up by objectRef
		destLocal = 1                               // The first parameter starts at index 1
		lenLocals++                                 // There is 1 more local needed
	}

	if globals.TraceVerbose {
		traceInfo := fmt.Sprintf("\tcreateAndInitNewFrame: lenArgList=%d, lenLocals=%d, stackSize=%d",
			lenArgList, lenLocals, stackSize)
		trace.Trace(traceInfo)
	}

	ptpx := 0
	for j := lenArgList - 1; j >= 0; j-- {
		fram.Locals[destLocal] = argList[j]
		switch paramsToPass[ptpx] {
		case "D", "J":
			destLocal += 2
		default:
			destLocal += 1
		}
		ptpx++
	}

	fram.TOS = -1

	return fram, nil
}
