package ff4go

import (
	"hash/fnv"
	"math"
)

func hashStringToFloat(flagName, id string) float64 {
	hash := fnv.New64a()
	hash.Write([]byte(flagName + ":" + id))

	return float64(hash.Sum64()) / math.MaxUint64
}
