package cryptutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandInt64(t *testing.T) {
	for i := 0; i < 100000; i++ {
		v1 := RandInt64(-10000, 100000)
		v2 := RandInt64(-10000, 100000)
		if v1 > v2 {
			v1, v2 = v2, v1
		}
		v := RandInt64(v1, v2)
		assert.True(t, v >= v1 && v <= v2)
	}

	for i := 0; i < 100000; i++ {
		v1 := RandInt64(10, 15)
		v2 := RandInt64(10, 15)
		if v1 > v2 {
			v1, v2 = v2, v1
		}
		v := RandInt64(v1, v2)
		assert.True(t, v >= v1 && v <= v2)
	}
}
