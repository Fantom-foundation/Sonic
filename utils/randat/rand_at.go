package randat

import (
	"math/rand/v2"
)

type cached struct {
	seed uint64
	r    uint64
}

var (
	gSeed = rand.Uint64()
	cache = cached{}
)

// RandAt returns random number with seed
// Not safe for concurrent use
func RandAt(seed uint64) uint64 {
	if seed != 0 && cache.seed == seed {
		return cache.r
	}
	cache.seed = seed
	cache.r = rand.New(rand.NewPCG(gSeed, seed)).Uint64()
	return cache.r
}
