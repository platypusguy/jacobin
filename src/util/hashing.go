package util

import (
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/cespare/xxhash/v2"
)

// For any argument, compute and return a Base64-encoded hash (string).
func HashAnything(arg interface{}) (uint64, error) {
	digest := xxhash.New()
	switch v := arg.(type) {
	case int64:
		barr := make([]byte, 8)
		binary.LittleEndian.PutUint64(barr, uint64(v))
		digest.Write(barr)
	case float64:
		barr := make([]byte, 8)
		binary.LittleEndian.PutUint64(barr, math.Float64bits(v))
		digest.Write(barr)
	case []byte:
		digest.Write(v)
	default:
		// Handle structs or other unknown types
		jsonBytes, err := json.Marshal(arg)
		if err != nil {
			return 0, err
		}
		digest.Write(jsonBytes)
	}
	return digest.Sum64(), nil
}
