package l

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueM_KeyExists(t *testing.T) {
	queueM := NewQueueM()

	queueM.PushBack("ab", 1)
	queueM.PushBack("abc", "2")
	queueM.PushBack("z", 3)
	queueM.PushBack("z", "4")
	queueM.PushBack("abc", 100)

	assert.True(t, queueM.KeyExists("ab"))
	assert.True(t, queueM.KeyExists("abc"))
	assert.True(t, queueM.KeyExists("z"))
	assert.False(t, queueM.KeyExists("c"))

	for is := range queueM.Iterator() {
		t.Log(is.K, is.V)
	}
}

func TestQueueM_Update(t *testing.T) {
	queueM := NewQueueM()

	queueM.PushBack("ab", 1)
	queueM.PushBack("abc", "2")
	queueM.PushBack("z", 3)
	assert.True(t, queueM.Update("abc", 2021))
	assert.False(t, queueM.Update("zzzz", "dd"))

	for is := range queueM.Iterator() {
		t.Log(is.K, is.V)
	}
}
