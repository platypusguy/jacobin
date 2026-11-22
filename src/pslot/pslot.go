package pslot

import (
	"jacobin/src/object"
)

// Slot is a single operand-stack or local-variable slot.
// Primitives are stored unboxed in Value; only references are pointers.
type Slot struct {
	Int   int64          // booleans, integers (8-, 16-, 32-, and 64-bits)
	Float float64        // Java float and double
	Ref   *object.Object // object pointers
}

func IntSlot(arg int64) Slot            { return Slot{Int: arg, Float: 0, Ref: nil} }
func Float64Slot(arg float64) Slot      { return Slot{Int: 0, Float: arg, Ref: nil} }
func Float32Slot(arg float32) Slot      { return Slot{Int: 0, Float: float64(arg), Ref: nil} }
func ObjSlot(obj *object.Object) Slot   { return Slot{Int: 0, Float: 0, Ref: obj} }
func AsInt64(slot Slot) int64           { return slot.Int }
func AsInt32(slot Slot) int32           { return int32(slot.Int) }
func AsInt16(slot Slot) int16           { return int16(slot.Int) }
func AsInt8(slot Slot) int8             { return int8(slot.Int) }
func AsFloat64(slot Slot) float64       { return slot.Float }
func AsFloat32(slot Slot) float32       { return float32(slot.Float) }
func AsObject(slot Slot) *object.Object { return slot.Ref }

