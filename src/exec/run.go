/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import (
	"errors"
	"fmt"
	"jacobin/log"
	"strconv"
)

var MainThread execThread

// StartExec is where execution begins. It initializes various structures, such as
// the VTable, then using the passed-in name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string) error {
	// initialize the VTable
	VTable = make(map[string]Ventry)
	VTableLoad()

	m, cpp, err := fetchMethodAndCP(className, "main")
	if err != nil {
		return errors.New("Class not found: " + className + ".main()")
	}

	f := createFrame(m.CodeAttr.MaxStack) // create a new frame
	f.methName = "main"
	f.clName = className
	f.cp = cpp                                  // add its pointer to the class CP
	for i := 0; i < len(m.CodeAttr.Code); i++ { // copy the bytecodes over
		f.meth = append(f.meth, m.CodeAttr.Code[i])
	}

	// allocate the local variables
	for k := 0; k < m.CodeAttr.MaxLocals; k++ {
		f.locals = append(f.locals, 0)
	}

	// create the first thread and place its first frame on it
	MainThread = CreateThread(0)
	f.thread = MainThread.id
	if pushFrame(&MainThread.stack, f) != nil {
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
	for t.stack.top > 0 {
		currFrame := t.stack.frames[t.stack.top]
		_ = runFrame(&currFrame)
		_ = popFrame(&t.stack)
	}
	return nil
}

func runFrame(f *frame) error {
	// f := *fr
	if f.ftype == 'G' { // if the frame contains a Golang method ('G')
		return runGframe(f) // run it differently
	}

	for f.pc < len(f.meth) {
		switch f.meth[f.pc] { // cases listed in numerical value of opcode
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
		case LDC: // 		0x12   	(push constant from CP indexed by next byte)
			push(f, int64(f.meth[f.pc+1]))
			f.pc += 1
		case ILOAD_0: // 	0x1A    (push local variable 0)
			push(f, f.locals[0])
		case ILOAD_1: //    OX1B    (push local variable 1)
			push(f, f.locals[1])
		case ILOAD_2: //    0X1C    (push local variable 2)
			push(f, f.locals[2])
		case ILOAD_3: //    0x1D    (push local variable 3)
			push(f, f.locals[3])
		case ISTORE_0: //   0x3B    (store popped top of stack int into local 0)
			f.locals[0] = pop(f)
		case ISTORE_1: //   0x3C    (store popped top of stack int into local 1)
			f.locals[1] = pop(f)
		case ISTORE_2: //   0x3D    (store popped top of stack int into local 2)
			f.locals[2] = pop(f)
		case ISTORE_3: //   0x3E    (store popped top of stack int into local 3)
			f.locals[3] = pop(f)
		case ISUB: //   0x64	(subtract top 2 items on operand stack, push result)
			i2 := pop(f)
			i1 := pop(f)
			push(f, i1-i2)
		case IINC: //   0x84    (increment local variable by a constant)
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
		case GOTO: // 0xA7     (goto an instruction)
			jumpTo := (int16(f.meth[f.pc+1]) * 256) + int16(f.meth[f.pc+2])
			f.pc = f.pc + int(jumpTo) - 1 // -1 because this loop will increment f.pc by 1
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
			if CPentry.Type != FieldRef { // the pointed-to CP entry must be a field reference
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
			fieldName := FetchUTF8stringFromCPEntryNumber(f.cp, fieldNameIndex)
			fieldName = className + "." + fieldName

			// was this static field previously loaded? Is so, get its location and move on.
			prevLoaded, ok := Statics[fieldName]
			if ok { // if preloaded, then push the index into the array of constant fields
				push(f, prevLoaded)
				break
			}

			fieldTypeIndex := nAndT.DescIndex
			fieldType := FetchUTF8stringFromCPEntryNumber(f.cp, fieldTypeIndex)
			// println("full field name: " + fieldName + ", type: " + fieldType)
			newStatic := Static{
				Class:     'L',
				Type:      fieldType,
				ValueRef:  "",
				ValueInt:  0,
				ValueFP:   0,
				ValueStr:  "",
				ValueFunc: nil,
				CP:        f.cp,
			}
			StaticsArray = append(StaticsArray, newStatic)
			Statics[fieldName] = int64(len(StaticsArray) - 1)

			// push the pointer to the stack of the frame
			push(f, int64(len(StaticsArray)-1))

		case INVOKEVIRTUAL: // 	0xB6 invokevirtual (create new frame, invoke function)
			CPslot := (int(f.meth[f.pc+1]) * 256) + int(f.meth[f.pc+2]) // next 2 bytes point to CP entry
			f.pc += 2
			CPentry := f.cp.CpIndex[CPslot]
			if CPentry.Type != MethodRef { // the pointed-to CP entry must be a field reference
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
			methodName := FetchUTF8stringFromCPEntryNumber(f.cp, methodNameIndex)
			methodName = className + "." + methodName
			// println("Method name for invokevirtual: " + methodName)

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := FetchUTF8stringFromCPEntryNumber(f.cp, methodSigIndex)
			// println("Method signature for invokevirtual: " + methodName + methodType)

			v := VTable[methodName+methodType]
			if v.Fu != nil && v.MethType == 'G' { // so we have a golang function in the queue
				gf := createFrame(v.ParamSlots)
				gf.thread = f.thread
				gf.methName = methodName + methodType
				gf.clName = className
				gf.meth = nil
				gf.cp = nil
				gf.locals = nil
				gf.ftype = 'G' // a golang function

				var argList []int64
				for i := 0; i < v.ParamSlots; i++ {
					arg := pop(f)
					argList = append(argList, arg)
				}
				for j := len(argList) - 1; j >= 0; j-- {
					push(&gf, argList[j])
				}
				gf.tos = len(gf.opStack) - 1
				pushFrame(&MainThread.stack, gf)
				runGframe(&gf)
				popFrame(&MainThread.stack)
				break
			}
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

// runs a frame whose method is a golang (so, native) method. It copies the parameters
// from the operand stack and passes them to the go function, here called Fu.
// TODO: Handle how return values are placed back on the stack.
func runGframe(fr *frame) error {
	ve := VTable[fr.methName]
	if ve.Fu == nil {
		return errors.New("go method not found: " + fr.methName)
	}

	var params = new([]interface{})
	for _, v := range fr.opStack {
		*params = append(*params, v)
	}

	ve.Fu(*params)

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
