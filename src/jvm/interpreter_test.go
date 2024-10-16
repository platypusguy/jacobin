package jvm

import (
	"jacobin/frames"
	"jacobin/opcodes"
	"testing"
)

func execOpCode(t *testing.T, opCode int, fptr *frames.Frame) {
	err := interpretBytecodes(opCode, fptr)
	if err != nil {
		t.Error(err.Error())
		return
	}
	str := emitTraceData(fptr)
	t.Logf(str)
}

func TestInterpreter(t *testing.T) {
	f := newFrame(opcodes.NOP)
	f.ClName = "KitchenActivity"
	f.MethName = "airFry"
	f.MethType = "([B)[B"
	f.Meth = []byte{opcodes.ACONST_NULL, opcodes.ICONST_M1, opcodes.ICONST_5, opcodes.LCONST_0, opcodes.FCONST_2, opcodes.DCONST_1}
	for _, opCode := range f.Meth {
		execOpCode(t, int(opCode), &f)
	}
}
