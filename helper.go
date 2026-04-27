package ff4go

import (
	"hash/fnv"
	"math"
	"time"
)

func hashStringToFloat(flagName, id string) float64 {
	hash := fnv.New64a()
	hash.Write([]byte(flagName + ":" + id))

	return float64(hash.Sum64()) / math.MaxUint64
}

func isExpired(date string) bool {
	if date == "" {
		return false
	}

	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return false
	}

	return t.Before(time.Now().UTC())
}
