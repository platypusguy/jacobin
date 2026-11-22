package pslot_test

import (
	"jacobin/src/frames"
	"jacobin/src/jvm"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/pslot"
	"math"
	"reflect"
	"testing"
)

func TestIntSlot(t *testing.T) {
	tests := []struct {
		name       string
		value      int64
		int32Value int32
		int16Value int16
		int8Value  int8
	}{
		{"zero", 0, 0, 0, 0},
		{"positive", 127, 127, 127, 127},
		{"negative", -129, -129, -129, 127},
		{"narrowing", 0x100008080, 32896, -32640, -128},
		{"max", math.MaxInt64, -1, -1, -1},
		{"min", math.MinInt64, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot := pslot.IntSlot(tt.value)
			if slot != (pslot.Slot{Int: tt.value}) {
				t.Fatalf("IntSlot(%d) = %+v", tt.value, slot)
			}
			if got := pslot.AsInt64(slot); got != tt.value {
				t.Errorf("AsInt64 = %d, want %d", got, tt.value)
			}
			if got := pslot.AsInt32(slot); got != tt.int32Value {
				t.Errorf("AsInt32 = %d, want %d", got, tt.int32Value)
			}
			if got := pslot.AsInt16(slot); got != tt.int16Value {
				t.Errorf("AsInt16 = %d, want %d", got, tt.int16Value)
			}
			if got := pslot.AsInt8(slot); got != tt.int8Value {
				t.Errorf("AsInt8 = %d, want %d", got, tt.int8Value)
			}
		})
	}
}

func TestFloatSlots(t *testing.T) {
	values := []float64{0, math.Copysign(0, -1), 1.25, -42.5, math.Inf(1), math.Inf(-1), math.NaN()}
	for _, value := range values {
		slot := pslot.Float64Slot(value)
		if slot.Int != 0 || slot.Ref != nil || math.Float64bits(pslot.AsFloat64(slot)) != math.Float64bits(value) {
			t.Errorf("Float64Slot(%v) = %+v", value, slot)
		}
		if got, want := math.Float32bits(pslot.AsFloat32(slot)), math.Float32bits(float32(value)); got != want {
			t.Errorf("AsFloat32(%v) bits = %x, want %x", value, got, want)
		}
	}
	for _, value := range []float32{0, float32(math.Copysign(0, -1)), 0.1, -42.5, float32(math.Inf(1))} {
		slot := pslot.Float32Slot(value)
		if slot.Int != 0 || slot.Ref != nil || math.Float64bits(slot.Float) != math.Float64bits(float64(value)) {
			t.Errorf("Float32Slot(%v) = %+v", value, slot)
		}
		if got := math.Float32bits(pslot.AsFloat32(slot)); got != math.Float32bits(value) {
			t.Errorf("AsFloat32(Float32Slot(%v)) bits = %x", value, got)
		}
	}
}

func TestObjSlot(t *testing.T) {
	for _, obj := range []*object.Object{nil, &object.Object{}} {
		slot := pslot.ObjSlot(obj)
		if slot != (pslot.Slot{Ref: obj}) || pslot.AsObject(slot) != obj {
			t.Errorf("ObjSlot(%p) = %+v; AsObject = %p", obj, slot, pslot.AsObject(slot))
		}
	}
}

func TestNopPreservesSlots(t *testing.T) {
	obj := &object.Object{}
	frame := &frames.Frame{
		Meth:    []byte{opcodes.NOP},
		Locals:  []interface{}{pslot.IntSlot(-42), pslot.ObjSlot(obj)},
		OpStack: []interface{}{pslot.Float64Slot(1.25), pslot.ObjSlot(obj)},
		TOS:     1,
	}
	before := *frame
	before.Locals = append([]interface{}(nil), frame.Locals...)
	before.OpStack = append([]interface{}(nil), frame.OpStack...)
	before.Meth = append([]byte(nil), frame.Meth...)
	if advance := jvm.DispatchTable[opcodes.NOP](frame, 0); advance != 1 {
		t.Errorf("NOP advance = %d, want 1", advance)
	}
	if !reflect.DeepEqual(*frame, before) {
		t.Errorf("NOP changed frame: got %+v, want %+v", *frame, before)
	}
}
