/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"strconv"
)

var MainThread execThread

// StartExec is where execution begins. It initializes various structures, such as
// the MTable, then using the passed-in name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string, globals *globals.Globals) error {
	// initialize the MTable
	classloader.MTable = make(map[string]classloader.MTentry)
	classloader.MTableLoadNatives()

	me, err := classloader.FetchMethodAndCP(className, "main", "([Ljava/lang/String;)V")
	if err != nil {
		return errors.New("Class not found: " + className + ".main()")
	}

	m := me.Meth.(classloader.JmEntry)
	f := createFrame(m.MaxStack) // create a new frame
	f.methName = "main"
	f.clName = className
	f.cp = m.Cp                        // add its pointer to the class CP
	for i := 0; i < len(m.Code); i++ { // copy the bytecodes over
		f.meth = append(f.meth, m.Code[i])
	}

	// allocate the local variables
	for k := 0; k < m.MaxLocals; k++ {
		f.locals = append(f.locals, 0)
	}

	// create the first thread and place its first frame on it
	MainThread = CreateThread(0)
	tracing := false
	trace, exists := globals.Options["-trace"]
	if exists {
		tracing = trace.Set
	}
	MainThread.trace = tracing
	f.thread = MainThread.id

	if pushFrame(MainThread.stack, f) != nil {
		_ = log.Log("Memory error allocating frame on thread: "+strconv.Itoa(MainThread.id), log.SEVERE)
		return errors.New("outOfMemory Exception")
	}

	err = runThread(&MainThread)
	if err != nil {
		return err
	}
	return nil
}

// Point the thread to the top of the frame stack and tell it to run from there.
func runThread(t *execThread) error {
	for t.stack.Len() > 0 {
		err := runFrame(t.stack)
		if err != nil {
			return err
		}

		if t.stack.Len() == 1 { // true when the last executed frame was main()
			return nil
		}
	}
	return nil
}

// runFrame() is the principal execution function in Jacobin. It first tests for a
// golang function in the present frame. If it is a golang function, it's sent to
// a different function for execution. Otherwise, bytecode interpretation takes
// place through a giant switch statement.
func runFrame(fs *list.List) error {
	// the current frame is always the head of the linked list of frames.
	// the next statement converts the address of that frame to the more readable 'f'
	f := fs.Front().Value.(*frame)

	// if the frame contains a golang method, execute it using runGframe(),
	// which returns a value (possibly nil) and an error code. Presuming no error,
	// if the return value (here, retval) is not nil, it is placed on the stack
	// of the calling frame.
	if f.ftype == 'G' {
		retval, err := runGframe(f)

		if retval != nil {
			f = fs.Front().Next().Value.(*frame)
			push(f, retval.(int64))
		}
		return err
	}

	// the frame's method is not a golang method, so it's Java bytecode, which
	// is interpreted in the rest of this function.
	for f.pc < len(f.meth) {
		if MainThread.trace {
			_ = log.Log("class: "+f.clName+
				", meth: "+f.methName+
				", pc: "+strconv.Itoa(f.pc)+
				", inst: "+BytecodeNames[int(f.meth[f.pc])]+
				", tos: "+strconv.Itoa(f.tos),
				log.TRACE_INST)
		}
		switch f.meth[f.pc] { // cases listed in numerical value of opcode
		case NOP:
			break
		case ICONST_N1: //	0x02	(push -1 onto opStack)
			push(f, -1)
		case ICONST_0: // 	0x03	(push 0 onto opStack)
			push(f, 0)
		case ICONST_1: //  	0x04	(push 1 onto opStack)
			push(f, 1)
		case ICONST_2: //   0x05	(push 2 onto opStack)
			push(f, 2)
		case ICONST_3: //   0x06	(push 3 onto opStack)
			push(f, 3)
		case ICONST_4: //   0x07	(push 4 onto opStack)
			push(f, 4)
		case ICONST_5: //   0x08	(push 5 onto opStack)
			push(f, 5)
		case BIPUSH: //     0x10	(push the following byte as an int onto the stack)
			push(f, int64(f.meth[f.pc+1]))
			f.pc += 1
		case LDC: // 	0x12   	(push constant from CP indexed by next byte)
			push(f, int64(f.meth[f.pc+1]))
			f.pc += 1
		case ILOAD_0: // 	0x1A    (push local variable 0)
			push(f, f.locals[0])
		case ILOAD_1: //    OX1B    (push local variable 1)
			push(f, f.locals[1])
		case ILOAD_2: //    0X1C    (push local variable 2)
			push(f, f.locals[2])
		case ILOAD_3: //  	0x1D   	(push local variable 3)
			push(f, f.locals[3])
		case LLOAD_0: //	0x1E	(push local variable 0, as long)
			push(f, f.locals[0])
		case LLOAD_1: //	0x1F	(push local variable 1, as long)
			push(f, f.locals[1])
		case LLOAD_2: //	0x20	(push local variable 2, as long)
			push(f, f.locals[2])
		case LLOAD_3: //	0x21	(push local variable 3, as long)
			push(f, f.locals[3])
		case ALOAD_0: //	0x2A	(push reference stored in local variable 0)
			push(f, f.locals[0])
		case ALOAD_1: //	0x2B	(push reference stored in local variable 1)
			push(f, f.locals[1])
		case ALOAD_2: //	0x2C    (push reference stored in local variable 2)
			push(f, f.locals[2])
		case ALOAD_3: //	0x2D	(push reference stored in local variable 3)
			push(f, f.locals[3])
		case ISTORE_0: //   0x3B    (store popped top of stack int into local 0)
			f.locals[0] = pop(f)
		case ISTORE_1: //   0x3C   	(store popped top of stack int into local 1)
			f.locals[1] = pop(f)
		case ISTORE_2: //   0x3D   	(store popped top of stack int into local 2)
			f.locals[2] = pop(f)
		case ISTORE_3: //   0x3E    (store popped top of stack int into local 3)
			f.locals[3] = pop(f)
		case LSTORE_0: //   0x3F    (store long from top of stack into locals 0 and 1)
			f.locals[0] = pop(f)
			f.locals[1] = f.locals[0]
		case LSTORE_1: //   0x40    (store long from top of stack into locals 1 and 2)
			f.locals[1] = pop(f)
			f.locals[2] = f.locals[1]
		case LSTORE_2: //   0x41    (store long from top of stack into locals 2 and 3)
			f.locals[2] = pop(f)
			f.locals[3] = f.locals[2]
		case LSTORE_3: //   0x42    (store long from top of stack into locals 3 and 4)
			f.locals[3] = pop(f)
			f.locals[4] = f.locals[3]
		case ASTORE_0: //	0x4B	(pop reference into local variable 0)
			f.locals[0] = pop(f)
		case ASTORE_1: //   0x4C	(pop reference into local variable 1)
			f.locals[1] = pop(f)
		case ASTORE_2: // 	0x4D	(pop reference into local variable 2)
			f.locals[2] = pop(f)
		case ASTORE_3: //	0x4E	(pop reference into local variable 3)
			f.locals[3] = pop(f)
		case IADD: //   0x60	(add top 2 items on operand stack, push result)
			i2 := pop(f)
			i1 := pop(f)
			push(f, i1+i2)
		case IMUL: //  0x68  	(multiply 2 items on operand stack, push result)
			i2 := pop(f)
			i1 := pop(f)
			push(f, i1*i2)
		case ISUB: //  0x64	(subtract top 2 items on operand stack, push result)
			i2 := pop(f)
			i1 := pop(f)
			push(f, i1-i2)
		case IINC: // 	0x84    (increment local variable by a constant)
			localVarIndex := int(f.meth[f.pc+1])
			constAmount := int(f.meth[f.pc+2])
			f.pc += 2
			f.locals[localVarIndex] += int64(constAmount)
		case IF_ICMPLT: //  0xA1    (jump if popped val1 < popped val2)
			val2 := pop(f)
			val1 := pop(f)
			if val1 < val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.meth[f.pc+1]) * 256) + int16(f.meth[f.pc+2])
				f.pc = f.pc + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.pc += 2
			}
		case IF_ICMPGE: //  0xA2    (jump if popped val1 >= popped val2)
			val2 := pop(f)
			val1 := pop(f)
			if val1 >= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.meth[f.pc+1]) * 256) + int16(f.meth[f.pc+2])
				f.pc = f.pc + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.pc += 2
			}
		case IF_ICMPLE: //	0xA4	(jump if popped val1 <= popped val2)
			val2 := pop(f)
			val1 := pop(f)
			if val1 <= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.meth[f.pc+1]) * 256) + int16(f.meth[f.pc+2])
				f.pc = f.pc + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.pc += 2
			}
		case GOTO: // 0xA7     (goto an instruction)
			jumpTo := (int16(f.meth[f.pc+1]) * 256) + int16(f.meth[f.pc+2])
			f.pc = f.pc + int(jumpTo) - 1 // -1 because this loop will increment f.pc by 1
		case IRETURN: // 0xAC (return an int and exit current frame)
			valToReturn := pop(f)
			f = fs.Front().Next().Value.(*frame)
			push(f, valToReturn) // TODO: check what happens when main() ends on IRETURN
			return nil
		case RETURN: // 0xB1    (return from void function)
			f.tos = -1 // empty the stack
			return nil
		case GETSTATIC: // 0xB2		(get static field)
			// TODO: getstatic will instantiate a static class if it's not already instantiated
			// that logic has not yet been implemented and the code here is simply a reasonable
			// placeholder, which consists of creating a struct that holds most of the needed info
			// puts it into a slice of such static fields and pushes the index of this item in the slice
			// onto the stack of the frame.
			CPslot := (int(f.meth[f.pc+1]) * 256) + int(f.meth[f.pc+2]) // next 2 bytes point to CP entry
			f.pc += 2
			CPentry := f.cp.CpIndex[CPslot]
			if CPentry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
				return fmt.Errorf("Expected a field ref on getstatic, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.pc, f.methName, f.clName)
			}

			// get the field entry
			field := f.cp.FieldRefs[CPentry.Slot]

			// get the class entry from the field entry for this field. It's the class name.
			classRef := field.ClassIndex
			classNameIndex := f.cp.ClassRefs[f.cp.CpIndex[classRef].Slot]
			classNameEntry := f.cp.CpIndex[classNameIndex]
			className := f.cp.Utf8Refs[classNameEntry.Slot]
			// println("Field name: " + className)

			// process the name and type entry for this field
			nAndTindex := field.NameAndType
			nAndTentry := f.cp.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := f.cp.NameAndTypes[nAndTslot]
			fieldNameIndex := nAndT.NameIndex
			fieldName := classloader.FetchUTF8stringFromCPEntryNumber(f.cp, fieldNameIndex)
			fieldName = className + "." + fieldName

			// was this static field previously loaded? Is so, get its location and move on.
			prevLoaded, ok := classloader.Statics[fieldName]
			if ok { // if preloaded, then push the index into the array of constant fields
				push(f, prevLoaded)
				break
			}

			fieldTypeIndex := nAndT.DescIndex
			fieldType := classloader.FetchUTF8stringFromCPEntryNumber(f.cp, fieldTypeIndex)
			// println("full field name: " + fieldName + ", type: " + fieldType)
			newStatic := classloader.Static{
				Class:     'L',
				Type:      fieldType,
				ValueRef:  "",
				ValueInt:  0,
				ValueFP:   0,
				ValueStr:  "",
				ValueFunc: nil,
				CP:        f.cp,
			}
			classloader.StaticsArray = append(classloader.StaticsArray, newStatic)
			classloader.Statics[fieldName] = int64(len(classloader.StaticsArray) - 1)

			// push the pointer to the stack of the frame
			push(f, int64(len(classloader.StaticsArray)-1))

		case INVOKEVIRTUAL: // 	0xB6 invokevirtual (create new frame, invoke function)
			CPslot := (int(f.meth[f.pc+1]) * 256) + int(f.meth[f.pc+2]) // next 2 bytes point to CP entry
			f.pc += 2
			CPentry := f.cp.CpIndex[CPslot]
			if CPentry.Type != classloader.MethodRef { // the pointed-to CP entry must be a method reference
				return fmt.Errorf("Expected a method ref for invokevirtual, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.pc, f.methName, f.clName)
			}

			// get the methodRef entry
			method := f.cp.MethodRefs[CPentry.Slot]

			// get the class entry from this method
			classRef := method.ClassIndex
			classNameIndex := f.cp.ClassRefs[f.cp.CpIndex[classRef].Slot]
			classNameEntry := f.cp.CpIndex[classNameIndex]
			className := f.cp.Utf8Refs[classNameEntry.Slot]

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := f.cp.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := f.cp.NameAndTypes[nAndTslot]
			methodNameIndex := nAndT.NameIndex
			methodName := classloader.FetchUTF8stringFromCPEntryNumber(f.cp, methodNameIndex)
			methodName = className + "." + methodName

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := classloader.FetchUTF8stringFromCPEntryNumber(f.cp, methodSigIndex)
			// println("Method signature for invokevirtual: " + methodName + methodType)

			v := classloader.MTable[methodName+methodType]
			if v.Meth != nil && v.MType == 'G' { // so we have a golang function
				// gFunc := v.meth.(GmEntry).Fu
				paramSlots := v.Meth.(classloader.GmEntry).ParamSlots
				gf := createFrame(paramSlots)
				gf.thread = f.thread
				gf.methName = methodName + methodType
				gf.clName = className
				gf.meth = nil
				gf.cp = nil
				gf.locals = nil
				gf.ftype = 'G' // a golang function

				var argList []int64
				for i := 0; i < paramSlots; i++ {
					arg := pop(f)
					argList = append(argList, arg)
				}
				for j := len(argList) - 1; j >= 0; j-- {
					push(gf, argList[j])
				}
				gf.tos = len(gf.opStack) - 1

				fs.PushFront(gf)              // push the new frame
				f = fs.Front().Value.(*frame) // point f to the new head

				err := runFrame(fs)
				if err != nil {
					return err
				}

				fs.Remove(fs.Front())         // pop the frame off
				f = fs.Front().Value.(*frame) // point f the head again
				break
			}
		case INVOKESTATIC: // 	0xB8 invokestatic (create new frame, invoke static function)
			CPslot := (int(f.meth[f.pc+1]) * 256) + int(f.meth[f.pc+2]) // next 2 bytes point to CP entry
			f.pc += 2
			CPentry := f.cp.CpIndex[CPslot]
			// get the methodRef entry
			method := f.cp.MethodRefs[CPentry.Slot]

			// get the class entry from this method
			classRef := method.ClassIndex
			classNameIndex := f.cp.ClassRefs[f.cp.CpIndex[classRef].Slot]
			classNameEntry := f.cp.CpIndex[classNameIndex]
			className := f.cp.Utf8Refs[classNameEntry.Slot]

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := f.cp.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := f.cp.NameAndTypes[nAndTslot]
			methodNameIndex := nAndT.NameIndex
			methodName := classloader.FetchUTF8stringFromCPEntryNumber(f.cp, methodNameIndex)
			// println("Method name for invokestatic: " + className + "." + methodName)

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := classloader.FetchUTF8stringFromCPEntryNumber(f.cp, methodSigIndex)
			// println("Method signature for invokestatic: " + methodName + methodType)

			// m, cpp, err := fetchMethodAndCP(className, methodName, methodType)
			mtEntry, err := classloader.FetchMethodAndCP(className, methodName, methodType)
			if err != nil {
				return errors.New("Class not found: " + className + methodName)
			}

			if mtEntry.MType == 'G' {
				runGmethod(mtEntry, fs, className, methodName, methodType)
			} else if mtEntry.MType == 'J' {
				m := mtEntry.Meth.(classloader.JmEntry)
				maxStack := m.MaxStack
				fram := createFrame(maxStack)

				fram.clName = className
				fram.methName = methodName
				fram.cp = m.Cp                     // add its pointer to the class CP
				for i := 0; i < len(m.Code); i++ { // copy the bytecodes over
					fram.meth = append(fram.meth, m.Code[i])
				}

				// allocate the local variables
				for k := 0; k < m.MaxLocals; k++ {
					fram.locals = append(fram.locals, 0)
				}

				// pop the parameters off the present stack and put them in the new frame's locals
				var argList []int64
				paramsToPass := ParseIncomingParamsFromMethTypeString(methodType)
				if len(paramsToPass) > 0 {
					for i := 0; i < len(paramsToPass); i++ {
						arg := pop(f)
						argList = append(argList, arg)
						if paramsToPass[i] == 'D' || paramsToPass[i] == 'J' {
							pop(f) // doubles and longs occupy two slots on the operand stack
						}
					}
				}

				destLocal := 0
				for j := len(argList) - 1; j >= 0; j-- {
					fram.locals[destLocal] = argList[j]
					destLocal += 1
				}
				fram.tos = -1

				fs.PushFront(fram)            // push the new frame
				f = fs.Front().Value.(*frame) // point f to the new head
				err = runFrame(fs)
				if err != nil {
					return err
				}

				// if the static method is main(), when we get here the
				// frame stack will be empty to exit from here, otherwise
				// there's still a frame on the stack, pop it off and continue.
				if fs.Len() == 0 {
					return nil
				}
				fs.Remove(fs.Front()) // pop the frame off

				// the previous frame pop might have been main()
				// if so, then we can't reset f to a non-existent frame
				// so we test for this before resetting f.
				if fs.Len() != 0 {
					f = fs.Front().Value.(*frame)
				} else {
					return nil
				}
			}
		case NEW: // 0xBB 	new: create and instantiate a new object
			CPslot := (int(f.meth[f.pc+1]) * 256) + int(f.meth[f.pc+2]) // next 2 bytes point to CP entry
			f.pc += 2
			CPentry := f.cp.CpIndex[CPslot]
			if CPentry.Type != classloader.ClassRef && CPentry.Type != classloader.Interface {
				msg := fmt.Sprintf("Invalid type for new object")
				_ = log.Log(msg, log.SEVERE)
			}

			// the classref points to a UTF8 record with the name of the class to instantiate
			var className string
			if CPentry.Type == classloader.ClassRef {
				utf8Index := f.cp.ClassRefs[CPentry.Slot]
				className = classloader.FetchUTF8stringFromCPEntryNumber(f.cp, utf8Index)
			}

			ref, err := instantiateClass(className)
			if err != nil {
				_ = log.Log("Error instantiating class: "+className, log.SEVERE)
				return errors.New("Error instnatiating class")
			}
			push(f, ref.(int64))

		default:
			msg := fmt.Sprintf("Invalid bytecode found: %d at location %d in method %s() of class %s\n",
				f.meth[f.pc], f.pc, f.methName, f.clName)
			_ = log.Log(msg, log.SEVERE)
			return errors.New("invalid bytecode encountered")
		}
		f.pc += 1
	}
	return nil
}

// runs a frame whose method is a golang method. It copies the parameters
// from the operand stack and passes them to the go function, here called Fu.
// Any return value from the method is returned to the call from run(), where
// it is placed on the stack of the calling function.
func runGframe(fr *frame) (interface{}, error) {
	// get the go method from the MTable
	me := classloader.MTable[fr.methName]
	if me.Meth == nil {
		return nil, errors.New("go method not found: " + fr.methName)
	}

	// pull arguments for the function off the frame's operand stack and put them in a slice
	var params = new([]interface{})
	for _, v := range fr.opStack {
		*params = append(*params, v)
	}

	// call the function passing a pointer to the slice of arguments
	ret := me.Meth.(classloader.GmEntry).Fu(*params)
	return ret, nil
}

func runGmethod(mt classloader.MTentry, fs *list.List, className, methodName, methodType string) error {
	f := fs.Front().Value.(*frame)

	paramSlots := mt.Meth.(classloader.GmEntry).ParamSlots
	gf := createFrame(paramSlots)
	gf.thread = f.thread
	gf.methName = className + "." + methodName + methodType
	gf.clName = className
	gf.meth = nil
	gf.cp = nil
	gf.locals = nil
	gf.ftype = 'G' // a golang function

	var argList []int64
	for i := 0; i < paramSlots; i++ {
		arg := pop(f)
		argList = append(argList, arg)
	}
	for j := len(argList) - 1; j >= 0; j-- {
		push(gf, argList[j])
	}
	gf.tos = len(gf.opStack) - 1

	fs.PushFront(gf)              // push the new frame
	f = fs.Front().Value.(*frame) // point f to the new head

	err := runFrame(fs) // this will eventually find its way to runGFrame()
	if err != nil {
		return err
	}

	fs.Remove(fs.Front())         // pop the frame off
	f = fs.Front().Value.(*frame) // point f the head again
	return nil
}

// pop from the operand stack. TODO: need to put in checks for invalid pops
func pop(f *frame) int64 {
	value := f.opStack[f.tos]
	f.tos -= 1
	return value
}

// push onto the operand stack
func push(f *frame, i int64) {
	f.tos += 1
	f.opStack[f.tos] = i
}
