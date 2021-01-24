package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "b", "c"}
	ret := reflect.DeepEqual(s1, s2)
	assert.True(t, ret)
	s3 := []string{"a", "c", "b"}
	ret = reflect.DeepEqual(s1, s3)
	assert.False(t, ret)
	s4 := make([]string, 0, 100)
	s4 = append(s4, "a")
	s4 = append(s4, "b")
	s4 = append(s4, "c")
	ret = reflect.DeepEqual(s1, s4)
	assert.True(t, ret)
}

func TestMap(t *testing.T) {
	m1 := make(map[string]string)
	m2 := make(map[string]string)
	m1["a"] = "aa"
	m1["b"] = "bb"
	m1["c"] = "cc"
	m2["a"] = "aa"
	m2["b"] = "bb"
	m2["c"] = "cc"
	ret := reflect.DeepEqual(m1, m2)
	assert.True(t, ret)

	m3 := make(map[string]string)
	m3["c"] = "cc"
	m3["a"] = "aa"
	m3["b"] = "bb"
	ret = reflect.DeepEqual(m1, m3)
	assert.True(t, ret)
}
