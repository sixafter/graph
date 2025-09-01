// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_FIFO(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a FIFO queue
	q := NewQueue(WithMode(ModeFIFO))

	// Enqueue elements
	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	// Check size
	is.Equal(3, q.Len())

	// Peek at the front element
	value, err := q.Peek()
	is.NoError(err)
	is.Equal("first", value)

	// Dequeue elements and validate FIFO order
	expected := []string{"first", "second", "third"}
	for _, exp := range expected {
		value, err := q.Dequeue()
		is.NoError(err)
		is.Equal(exp, value)
	}

	// Check if the queue is empty
	is.True(q.IsEmpty())

	// Attempt to dequeue from an empty queue
	_, err = q.Dequeue()
	is.Error(err)
	is.Equal("queue is empty", err.Error())
}

func TestQueue_LIFO(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a LIFO queue
	q := NewQueue(WithMode(ModeLIFO))

	// Enqueue elements
	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	// Check size
	is.Equal(3, q.Len())

	// Peek at the top element
	value, err := q.Peek()
	is.NoError(err)
	is.Equal("third", value)

	// Dequeue elements and validate LIFO order
	expected := []string{"third", "second", "first"}
	for _, exp := range expected {
		value, err := q.Dequeue()
		is.NoError(err)
		is.Equal(exp, value)
	}

	// Check if the queue is empty
	is.True(q.IsEmpty())

	// Attempt to dequeue from an empty queue
	_, err = q.Dequeue()
	is.Error(err)
	is.Equal("queue is empty", err.Error())
}

func TestQueue_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a FIFO queue
	q := NewQueue(WithMode(ModeFIFO))

	// Define the number of operations
	const numOps = 1000

	// Use wait groups to synchronize enqueue and dequeue operations
	var wg sync.WaitGroup
	wg.Add(2) // Two goroutines: one for enqueue and one for dequeue

	// Enqueue elements in a goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			q.Enqueue(i)
		}
	}()

	// Dequeue elements in a goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			// Dequeue with retry to handle empty queue errors
			for {
				_, err := q.Dequeue()
				if err == nil {
					break
				}
			}
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Verify the queue is empty after concurrent operations
	is.True(q.IsEmpty(), "Queue should be empty after all operations")
}

func TestQueue_DefaultMode(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a queue with the default mode (FIFO)
	q := NewQueue()

	// Enqueue elements
	q.Enqueue("first")
	q.Enqueue("second")

	// Dequeue elements and validate FIFO order
	expected := []string{"first", "second"}
	for _, exp := range expected {
		value, err := q.Dequeue()
		is.NoError(err)
		is.Equal(exp, value)
	}
}
