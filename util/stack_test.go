package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackPush(t *testing.T) {
	s := NewStack()
	for i := 0; i < 3; i++ {
		s.Push(i)
	}
	assert.Equal(t, 3, s.Size())
}

func TestStackPeek(t *testing.T) {
	s := Stack{0, 1, 2}
	assert.Equal(t, 2, s.Peek())
	assert.Equal(t, 3, s.Size())
}

func TestStackPop(t *testing.T) {
	s := Stack{2, 1, 0}
	for i := 0; i < 3; i++ {
		assert.Equal(t, i, s.Pop())
	}
	assert.Equal(t, 0, s.Size())
}

func TestStackPopEmpty(t *testing.T) {
	s := NewStack()
	assert.Panics(t, func() {
		s.Pop()
	})
}
