// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package queue

import (
	"errors"
	"sync"
)

// Mode defines the type for the queue's operation mode.
type Mode int

const (
	ModeFIFO Mode = iota // First-In-First-Out
	ModeLIFO             // Last-In-First-Out
)

// Queue is a thread-safe queue that can operate in FIFO or LIFO mode.
type Queue struct {
	queue []any
	lock  sync.RWMutex
	mode  Mode
}

// NewQueue creates a new queue with the specified mode.
func NewQueue(options ...func(*Queue)) *Queue {
	q := &Queue{
		queue: make([]any, 0),
		mode:  ModeFIFO, // Default to FIFO
	}
	for _, option := range options {
		option(q)
	}
	return q
}

// WithMode sets the mode of the queue (FIFO or LIFO).
func WithMode(mode Mode) func(*Queue) {
	return func(q *Queue) {
		q.mode = mode
	}
}

// Enqueue adds an element to the queue.
func (q *Queue) Enqueue(value any) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, value)
}

// Dequeue removes and returns an element from the queue, depending on the mode.
// Returns an error if the queue is empty.
func (q *Queue) Dequeue() (any, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 {
		return nil, errors.New("queue is empty")
	}

	var value any
	if q.mode == ModeFIFO {
		// FIFO: Remove from the front
		value = q.queue[0]
		q.queue = q.queue[1:]
	} else {
		// LIFO: Remove from the back
		value = q.queue[len(q.queue)-1]
		q.queue = q.queue[:len(q.queue)-1]
	}
	return value, nil
}

// Peek returns the next element without removing it, depending on the mode.
// Returns an error if the queue is empty.
func (q *Queue) Peek() (any, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()
	if len(q.queue) == 0 {
		return nil, errors.New("queue is empty")
	}

	if q.mode == ModeFIFO {
		return q.queue[0], nil // FIFO: Peek at the front
	}
	return q.queue[len(q.queue)-1], nil // LIFO: Peek at the back
}

// IsEmpty checks if the queue is empty.
func (q *Queue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.queue) == 0
}

// Len returns the number of elements in the queue.
func (q *Queue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.queue)
}
