// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	t.Parallel()

	t.Run("Push and Pop items in priority order", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		pq := NewPriorityQueue[string]()
		pq.Push("low", 10.0)
		pq.Push("medium", 5.0)
		pq.Push("high", 1.0)

		item, err := pq.Pop()
		is.NoError(err)
		is.Equal("high", item, "Expected item with highest priority (lowest value)")

		item, err = pq.Pop()
		is.NoError(err)
		is.Equal("medium", item, "Expected next highest priority item")

		item, err = pq.Pop()
		is.NoError(err)
		is.Equal("low", item, "Expected last item with lowest priority")

		_, err = pq.Pop()
		is.Error(err)
		is.ErrorIs(err, ErrPriorityQueueEmpty)
	})

	t.Run("SetPriority updates item priority", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		pq := NewPriorityQueue[string]()
		pq.Push("task1", 3.0)
		pq.Push("task2", 2.0)
		pq.SetPriority("task1", 1.0)

		item, err := pq.Pop()
		is.NoError(err)
		is.Equal("task1", item, "task1 priority was updated to highest")
	})
}

func TestStack(t *testing.T) {
	t.Parallel()

	t.Run("Push and Pop items", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		s := NewStack[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)

		item, ok := s.Pop()
		is.True(ok)
		is.Equal(3, item, "Expected last pushed item")

		item, ok = s.Pop()
		is.True(ok)
		is.Equal(2, item, "Expected second-to-last pushed item")

		item, ok = s.Pop()
		is.True(ok)
		is.Equal(1, item, "Expected first pushed item")

		_, ok = s.Pop()
		is.False(ok, "Stack should be empty")
	})

	t.Run("Contains checks item existence", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		s := NewStack[int]()
		s.Push(1)
		s.Push(2)

		is.True(s.Contains(1))
		is.True(s.Contains(2))
		is.False(s.Contains(3))
	})

	t.Run("IsEmpty detects empty Stack", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		s := NewStack[int]()
		is.True(s.IsEmpty())

		s.Push(1)
		is.False(s.IsEmpty())
	})
}

func TestStackOfStacks(t *testing.T) {
	t.Parallel()

	t.Run("Push and Pop stacks", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		sos := NewStackOfStacks[int]()
		stack1 := NewStack[int]()
		stack2 := NewStack[int]()

		stack1.Push(1)
		stack2.Push(2)

		sos.Push(stack1)
		sos.Push(stack2)

		topStack, err := sos.Pop()
		is.NoError(err)
		topItem, ok := topStack.Pop()
		is.True(ok)
		is.Equal(2, topItem, "Expected item from the Top Stack")

		topStack, err = sos.Pop()
		is.NoError(err)
		topItem, ok = topStack.Pop()
		is.True(ok)
		is.Equal(1, topItem, "Expected item from the second Stack")

		_, err = sos.Pop()
		is.Error(err)
		is.ErrorIs(err, ErrStackEmpty)
	})

	t.Run("Top returns the Top Stack without removing it", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		sos := NewStackOfStacks[int]()
		stack := NewStack[int]()
		stack.Push(42)

		sos.Push(stack)

		topStack, err := sos.Top()
		is.NoError(err)

		topItem, ok := topStack.Top()
		is.True(ok)
		is.Equal(42, topItem, "Expected item from the Top Stack")
	})
}
