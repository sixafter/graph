// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueuePushPop(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	pq := NewPriorityQueue[string]()
	pq.Enqueue("low", 10.0)
	pq.Enqueue("medium", 5.0)
	pq.Enqueue("high", 1.0)

	item, err := pq.Dequeue()
	is.NoError(err)
	is.Equal("high", item, "Expected item with highest priority (lowest value)")

	item, err = pq.Dequeue()
	is.NoError(err)
	is.Equal("medium", item, "Expected next highest priority item")

	item, err = pq.Dequeue()
	is.NoError(err)
	is.Equal("low", item, "Expected last item with lowest priority")

	_, err = pq.Dequeue()
	is.Error(err)
	is.ErrorIs(err, ErrPriorityQueueEmpty)
}

func TestPriorityQueueSetPriority(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	pq := NewPriorityQueue[string]()
	pq.Enqueue("task1", 3.0)
	pq.Enqueue("task2", 2.0)
	pq.SetPriority("task1", 1.0)

	item, err := pq.Dequeue()
	is.NoError(err)
	is.Equal("task1", item, "task1 priority was updated to highest")
}

func TestStackPushPop(t *testing.T) {
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
}

func TestStackContains(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s := NewStack[int]()
	s.Push(1)
	s.Push(2)

	is.True(s.Contains(1))
	is.True(s.Contains(2))
	is.False(s.Contains(3))
}

func TestStackIsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s := NewStack[int]()
	is.True(s.IsEmpty())

	s.Push(1)
	is.False(s.IsEmpty())
}

func TestStackOfStacksPushPop(t *testing.T) {
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
	is.Equal(2, topItem, "Expected item from the top stack")

	topStack, err = sos.Pop()
	is.NoError(err)
	topItem, ok = topStack.Pop()
	is.True(ok)
	is.Equal(1, topItem, "Expected item from the second stack")

	_, err = sos.Pop()
	is.Error(err)
	is.ErrorIs(err, ErrStackEmpty)
}

func TestStackOfStacksTop(t *testing.T) {
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
	is.Equal(42, topItem, "Expected item from the top stack")
}

func TestPriorityQueueUpdateNonexistentItem(t *testing.T) {
	t.Parallel()

	pq := NewPriorityQueue[string]()
	pq.Enqueue("task1", 5.0)

	pq.SetPriority("task2", 1.0)
}

func TestStackUnderflow(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s := NewStack[int]()
	_, ok := s.Pop()
	is.False(ok, "Popping from an empty stack should fail")
}

func TestStackOfStacksUnderflow(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	sos := NewStackOfStacks[int]()
	_, err := sos.Pop()
	is.Error(err, "Popping from an empty stack of stacks should fail")
	is.ErrorIs(err, ErrStackEmpty)
}

func TestPriorityQueueEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	pq := NewPriorityQueue[string]()
	_, err := pq.Dequeue()
	is.Error(err, "Popping from an empty priority queue should fail")
	is.ErrorIs(err, ErrPriorityQueueEmpty)
}
