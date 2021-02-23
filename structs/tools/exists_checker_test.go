package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistsCheckerWithMaxSize(t *testing.T) {
	ec := NewExistsCheckerWithMaxSize(3)
	ec.Add("1")
	assert.True(t, ec.Exists("1"))
	assert.False(t, ec.Exists("2"))
	ec.Add("2")
	assert.True(t, ec.Exists("1"))
	assert.True(t, ec.Exists("2"))
	ec.Add("3")
	ec.Add("4")

	assert.False(t, ec.Exists("1"))
	assert.True(t, ec.Exists("2"))
	assert.True(t, ec.Exists("3"))
	assert.True(t, ec.Exists("4"))
}

func TestExistsCheckerWithMaxSize2(t *testing.T) {
	ec := NewExistsCheckerWithMaxSize(3)
	ec.Add("1")
	ec.Add("2")
	ec.Add("3")
	ec.Add("1")
	ec.Add("4")
	assert.True(t, ec.Exists("1"))
	assert.False(t, ec.Exists("2"))
	assert.True(t, ec.Exists("3"))
	assert.True(t, ec.Exists("4"))
}
