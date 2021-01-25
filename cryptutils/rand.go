package cryptutils

import (
	"crypto/rand"
	"math/big"
)

// RandInt64 random int64
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	maxBigInt := big.NewInt(max - min)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	return min + i.Int64()
}

// RandInt random int
func RandInt(min, max int) int {
	return int(RandInt64(int64(min), int64(max)))
}

// RandUInt random uint
func RandUInt(min, max uint) uint {
	return uint(RandInt64(int64(min), int64(max)))
}
