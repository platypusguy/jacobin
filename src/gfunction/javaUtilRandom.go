package gfunction

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/types"
	"math"
	"math/rand"
	"time"
)

func Load_Util_Random() map[string]GMeth {

	MethodSignatures["java/util/Random.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/Random.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomInitVoid,
		}

	MethodSignatures["java/util/Random.<init>(J)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  randomInitLong,
		}

	MethodSignatures["java/util/Random.nextBoolean()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomNextBoolean,
		}

	MethodSignatures["java/util/Random.nextBytes([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  randomNextBytes,
		}

	MethodSignatures["java/util/Random.nextDouble()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomNextDouble,
		}

	MethodSignatures["java/util/Random.nextFloat()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomNextFloat,
		}

	MethodSignatures["java/util/Random.nextGaussian()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomNextGaussian,
		}

	MethodSignatures["java/util/Random.nextInt()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomNextInt,
		}

	MethodSignatures["java/util/Random.nextInt(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  randomNextIntBound,
		}

	MethodSignatures["java/util/Random.nextLong()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  randomNextLong,
		}

	MethodSignatures["java/util/Random.nextLong(J)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  randomNextLongBound,
		}

	MethodSignatures["java/util/Random.setSeed(J)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  randomSetSeed,
		}

	return MethodSignatures
}

/*
ChatGPT key Points:
* Initialization: The NewRandom() function initializes a Random instance using the current time as a seed.
* Seed Setting: The SetSeed(seed int64) method allows setting a specific seed for the random number generator.
* Random Number Generation: Methods like NextInt(), NextIntBound(bound int), NextFloat32(), NextFloat64(),
                                         NextBoolean(), and NextGaussian()
                            mimic the behavior of their Java counterparts using Go's math/rand package.
* Concurrency: A sync.Mutex is used to ensure thread safety when accessing shared state like haveNextNextGaussian
               and nextNextGaussian in NextGaussian().

* object.Object Ftype = types.Struct
*/

type Random struct {
	rand                 *rand.Rand
	nextNextGaussian     float64
	haveNextNextGaussian bool
}

// Primitive to update a Random object with a Random struct.
func UpdateRandomObjectFromStruct(objPtr *object.Object, argStruct Random) {
	fld := object.Field{Ftype: types.Struct, Fvalue: argStruct}
	objPtr.FieldTable["value"] = fld
}

// Primitive to fetch a Random struct from a Random object.
func GetStructFromRandomObject(objPtr *object.Object) Random {
	randStruct := objPtr.FieldTable["value"].Fvalue.(Random)
	return randStruct
}

// "java/util/Random.<init>()V"
// NewRandom creates a new Random instance initialized with the current time as seed.
// chatGPT generated: func NewRandom() *Random
func randomInitVoid(params []interface{}) interface{} {
	source := rand.NewSource(time.Now().UnixNano())
	randStruct := Random{
		rand:                 rand.New(source),
		nextNextGaussian:     0.0,
		haveNextNextGaussian: false,
	}
	obj := params[0].(*object.Object)
	UpdateRandomObjectFromStruct(obj, randStruct)
	return nil
}

// "java/util/Random.<init>(J)V"
// Same as randomInitVoid except a seed is supplied.
func randomInitLong(params []interface{}) interface{} {
	seed := params[1].(int64)
	source := rand.NewSource(seed)
	randStruct := Random{
		rand:                 rand.New(source),
		nextNextGaussian:     0.0,
		haveNextNextGaussian: false,
	}
	obj := params[0].(*object.Object)
	UpdateRandomObjectFromStruct(obj, randStruct)
	return nil
}

// randomSetSeed sets the seed of the random number generator.
// ChatGPT: func (r *Random) SetSeed(seed int64)
func randomSetSeed(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.RandomLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	seed := params[1].(int64)
	r.rand.Seed(seed)
	r.haveNextNextGaussian = false
	UpdateRandomObjectFromStruct(obj, r)
	global.RandomLock.Unlock() // <-------------------
	return nil
}

// randomNextInt returns the next pseudorandom, uniformly distributed int64 value.
// ChatGPT: func (r *Random) NextInt() int {
func randomNextInt(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	output := r.rand.Int63()
	return output
}

// randomNextIntBound returns a pseudorandom, uniformly distributed int value between 0 (inclusive) and bound (exclusive).
// ChatGPT: func (r *Random) NextIntBound(bound int) int, error {
func randomNextIntBound(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	bound := params[1].(int64)
	if bound < 1 {
		errMsg := fmt.Sprintf("Random.NextIntBound: bound must be positive, observed: %d", bound)
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}
	output := r.rand.Int63n(bound)
	return output
}

// randomNextLong is identical to randomNextInt.
// ChatGPT: func (r *Random) NextLong() int64 {
func randomNextLong(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	output := r.rand.Int63()
	return output
}

// randomNextLongBound returns a pseudorandom, uniformly distributed long value between 0 (inclusive) and bound (exclusive).
func randomNextLongBound(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	bound := params[1].(int64)
	if bound < 1 {
		errMsg := fmt.Sprintf("Random.NextLongBound: bound must be positive, observed: %d", bound)
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}
	output := r.rand.Int63n(bound)
	return output
}

// randomNextBoolean returns the next pseudorandom, uniformly distributed boolean value.
// ChatGPT: func (r *Random) NextBoolean() bool
func randomNextBoolean(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	if r.rand.Int63n(2) == 0 {
		return int64(1)
	} else {
		return int64(0)
	}
}

// Given an array of bytes, fill each element with a random number between 0 and 255.
func randomNextBytes(params []interface{}) interface{} {
	robj := params[0].(*object.Object)
	r := GetStructFromRandomObject(robj)
	bobj := params[1].(*object.Object)
	bytes := bobj.FieldTable["value"].Fvalue.([]byte)
	r.rand.Read(bytes)
	bobj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: bytes}
	return nil
}

// randomNextFloat returns the next pseudorandom, uniformly distributed float64 value between 0.0 (inclusive) and 1.0 (exclusive).
// ChatGPT: func (r *Random) NextFloat32() float32
func randomNextFloat(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	return r.rand.Float64()
}

// randomNextDouble is identical to randomNextFloat.
// ChatGPT: func (r *Random) NextFloat64() float64
func randomNextDouble(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)
	return r.rand.Float64()
}

// NextGaussian returns the next pseudorandom, Gaussian ("normally") distributed float64 value with mean 0.0 and standard deviation 1.0.
// ChatGPT: func (r *Random) NextGaussian() float64
func randomNextGaussian(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.RandomLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	r := GetStructFromRandomObject(obj)

	if r.haveNextNextGaussian {
		r.haveNextNextGaussian = false
		UpdateRandomObjectFromStruct(obj, r)
		global.RandomLock.Unlock() // <-------------------
		return r.nextNextGaussian
	}

	var v1, v2, s float64
	for {
		v1 = 2*r.rand.Float64() - 1 // between -1.0 and 1.0
		v2 = 2*r.rand.Float64() - 1 // between -1.0 and 1.0
		s = v1*v1 + v2*v2
		if s < 1.0 && s != 0.0 {
			break
		}
	}

	multiplier := math.Sqrt(-2 * math.Log(s) / s)
	r.nextNextGaussian = v2 * multiplier
	r.haveNextNextGaussian = true
	UpdateRandomObjectFromStruct(obj, r)

	global.RandomLock.Unlock() // <-------------------

	return v1 * multiplier
}
