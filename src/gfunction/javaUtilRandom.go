package gfunction

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

/*
ChatGPT

Key Points:

Initialization: The NewRandom() function initializes a Random instance using the current time as a seed.

Seed Setting: The SetSeed(seed int64) method allows setting a specific seed for the random number generator.

Random Number Generation: Methods like NextInt(), NextIntBound(bound int), NextFloat32(), NextFloat64(),
                                       NextBoolean(), and NextGaussian()
mimic the behavior of their Java counterparts using Go's math/rand package.

Concurrency: A sync.Mutex (r.mu) is used to ensure thread safety when accessing shared state like haveNextNextGaussian
             and nextNextGaussian in NextGaussian().
*/

// Random implements a random number generator.
type Random struct {
	rand                 *rand.Rand
	nextNextGaussian     float64
	haveNextNextGaussian bool
	mu                   sync.Mutex
}

// NewRandom creates a new Random instance initialized with the current time as seed.
func NewRandom() *Random {
	source := rand.NewSource(time.Now().UnixNano())
	return &Random{
		rand: rand.New(source),
	}
}

// SetSeed sets the seed of the random number generator.
func (r *Random) SetSeed(seed int64) {
	r.rand.Seed(seed)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.haveNextNextGaussian = false
}

// NextInt returns the next pseudorandom, uniformly distributed int value.
func (r *Random) NextInt() int {
	return r.rand.Int()
}

// NextIntBound returns a pseudorandom, uniformly distributed int value between 0 (inclusive) and bound (exclusive).
func (r *Random) NextIntBound(bound int) (int, error) {
	if bound <= 0 {
		return 0, errors.New("bound must be positive")
	}
	return r.rand.Intn(bound), nil
}

// NextLong returns the next pseudorandom, uniformly distributed int64 value.
func (r *Random) NextLong() int64 {
	return int64(r.rand.Uint64())
}

// NextBoolean returns the next pseudorandom, uniformly distributed boolean value.
func (r *Random) NextBoolean() bool {
	return r.rand.Intn(2) == 0
}

// NextFloat32 returns the next pseudorandom, uniformly distributed float32 value between 0.0 (inclusive) and 1.0 (exclusive).
func (r *Random) NextFloat32() float32 {
	return r.rand.Float32()
}

// NextFloat64 returns the next pseudorandom, uniformly distributed float64 value between 0.0 (inclusive) and 1.0 (exclusive).
func (r *Random) NextFloat64() float64 {
	return r.rand.Float64()
}

// NextGaussian returns the next pseudorandom, Gaussian ("normally") distributed float64 value with mean 0.0 and standard deviation 1.0.
func (r *Random) NextGaussian() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.haveNextNextGaussian {
		r.haveNextNextGaussian = false
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

	return v1 * multiplier
}
