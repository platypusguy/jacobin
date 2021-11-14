/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import (
	"errors"
	"fmt"
	"os"
)

// StartExec accepts the name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string) error {
	m, cpp, err := fetchMethodAndCP(className, "main")
	if err != nil {
		return errors.New("Class not found: " + className + ".main()")
	}

	f := frame{} // create a new frame
	f.clName = className
	f.cp = cpp                                  // add its pointer to the class CP
	for i := 0; i < len(m.CodeAttr.Code); i++ { // copy the bytecodes over
		f.meth = append(f.meth, m.CodeAttr.Code[i])
	}

	// allocate the operand stack
	for j := 0; j < m.CodeAttr.MaxStack; j++ {
		f.opStack = append(f.opStack, int32(0))
	}
	f.tos = -1

	// allocate the local variables
	for k := 0; k < m.CodeAttr.MaxLocals; k++ {
		f.locals = append(f.locals, 0)
	}

	t := CreateThread(0)
	f.thread = t.id
	pushFrame(&t.stack, f)

	err = runThread(t)
	if err != nil {
		return err
	}
	return nil
}

func runThread(t execThread) error {
	currFrame := t.stack.frames[t.stack.top]
	return runFrame(currFrame)
}

func runFrame(f frame) error {
	for pc := 0; pc < len(f.meth); pc++ {
		switch f.meth[pc] {
		case 0x02: // iconst_n1    (push -1 onto opStack)
			push(&f, -1)
		case 0x03: // iconst_0     (push 0 onto opStack)
			push(&f, 0)
		case 0x04: // iconst_1     (push 1 onto opStack)
			push(&f, 1)
		case 0x05: // iconst_2     (push 2 onto opStack)
			push(&f, 2)
		case 0x06: // iconst_3     (push 3 onto opStack)
			push(&f, 3)
		case 0x07: // iconst_4     (push 4 onto opStack)
			push(&f, 4)
		case 0x08: // iconst_5     (push 5 onto opStack)
			push(&f, 5)
		case 0x10: // bipush       push the following byte as an int onto the stack
			push(&f, int32(f.meth[pc+1]))
			pc += 1
		case 0x1A: // iload_0      (push local variable 0)
			push(&f, f.locals[0])
		case 0x1B: // iload_1      (push local variable 1)
			push(&f, f.locals[1])
		case 0x1C: // iload_2      (push local variable 2)
			push(&f, f.locals[2])
		case 0x1D: // iload_3      (push local variable 3)
			push(&f, f.locals[3])
		case 0x3B: // istore_0     (store popped top of stack int into local 0)
			f.locals[0] = pop(&f)
		case 0x3C: // istore_1     (store popped top of stack int into local 1)
			f.locals[1] = pop(&f)
		case 0x3D: // istore_2     (store popped top of stack int into local 2)
			f.locals[2] = pop(&f)
		case 0x3E: // istore_3     (store popped top of stack int into local 3)
			f.locals[3] = pop(&f)
		case 0xA2: // icmpge       (jump if popped val1 >= popped val2)
			val2 := pop(&f)
			val1 := pop(&f)
			if val1 >= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int(f.meth[pc+1]) * 256) + int(f.meth[pc+2])
				pc = jumpTo - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				pc += 2
			}
		default:
			fmt.Fprintf(os.Stderr, "Invalid bytecode found: %d at location: %d in method %s\n",
				f.meth[pc], pc, f.clName)
			return errors.New("invalid bytecode encountered")
		}
	}
	return nil
}

// pop from the operand stack. TODO: need to put in checks for invalid pops
func pop(f *frame) int32 {
	value := f.opStack[f.tos]
	f.tos -= 1
	return value
}

// push onto the operand stack
func push(f *frame, i int32) {
	f.tos += 1
	f.opStack[f.tos] = i
}
