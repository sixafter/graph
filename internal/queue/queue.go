// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package queue

import (
	"container/heap"
	"errors"
)

var (
	// ErrStackEmpty is returned when an attempt is made to Pop or access the Top
	// element from an empty Stack or Stack of stacks. This error helps prevent
	// invalid operations on empty data structures.
	ErrStackEmpty = errors.New("no element in Stack")

	// ErrPriorityQueueEmpty is returned when an attempt is made to remove an item
	// from an empty priority queue. This ensures that the caller is informed
	// that the operation could not be completed due to the absence of elements.
	ErrPriorityQueueEmpty = errors.New("priority queue is empty")
)

// PriorityQueue is a generic data structure that provides an efficient way to
// manage elements with associated priorities. It implements a minimum priority
// queue, where elements with smaller priority values are dequeued before elements
// with larger priority values.
//
// Internally, the PriorityQueue uses a binary heap to maintain the heap property,
// ensuring that the smallest-priority element is always at the front. The queue
// supports dynamic updates to element priorities and efficient insertion and
// removal operations.
//
// Type Parameters:
//
//	T - The type of the elements in the queue. It must be comparable.
//
// Key Features:
//   - Fast insertion of new elements (O(log n)).
//   - Efficient removal of the smallest-priority element (O(log n)).
//   - Priority updates in logarithmic time (O(log n)).
//   - Constant-time element existence checks via an internal cache.
//
// Use Cases:
//   - Task scheduling based on priority.
//   - Pathfinding algorithms (e.g., Dijkstra’s or A*).
//   - Event simulation systems.
//
// Example Usage:
//
//	pq := NewPriorityQueue[string]()
//	pq.Push("task1", 2.0)
//	pq.Push("task2", 1.0)
//	task, err := pq.Pop()
//	if err != nil {
//	    fmt.Println(err)
//	} else {
//	    fmt.Println(task) // Output: task2
//	}
type PriorityQueue[T comparable] struct {
	items *minHeap[T]
	cache map[T]*PriorityItem[T]
}

// PriorityItem represents an item in the binary heap, consisting of a priority value
// and an actual payload value.
type PriorityItem[T comparable] struct {
	value    T
	priority float64
	index    int
}

// NewPriorityQueue creates and returns a new empty priority queue.
//
// Example:
//
//	pq := NewPriorityQueue[string]()
//	pq.Push("task1", 2.0)
//	pq.Push("task2", 1.0)
func NewPriorityQueue[T comparable]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: &minHeap[T]{},
		cache: map[T]*PriorityItem[T]{},
	}
}

// Len returns the total number of items in the priority queue.
//
// Example:
//
//	pq := NewPriorityQueue[string]()
//	fmt.Println(pq.Len()) // Output: 0
func (p *PriorityQueue[T]) Len() int {
	return p.items.Len()
}

// Push adds a new item with the given priority to the priority queue.
// If the item already exists, it does nothing.
//
// Example:
//
//	pq := NewPriorityQueue[string]()
//	pq.Push("task1", 2.0)
//	pq.Push("task2", 1.0)
func (p *PriorityQueue[T]) Push(item T, priority float64) {
	if _, ok := p.cache[item]; ok {
		return
	}

	newItem := &PriorityItem[T]{
		value:    item,
		priority: priority,
		index:    0,
	}

	heap.Push(p.items, newItem)
	p.cache[item] = newItem
}

// Pop removes and returns the item with the lowest priority from the priority queue.
// Returns ErrPriorityQueueEmpty if the queue is empty.
//
// Example:
//
//	pq := NewPriorityQueue[string]()
//	pq.Push("task1", 2.0)
//	task, err := pq.Pop()
//	if err != nil {
//	    fmt.Println(err)
//	} else {
//	    fmt.Println(task) // Output: task1
//	}
func (p *PriorityQueue[T]) Pop() (T, error) {
	if len(*p.items) == 0 {
		var empty T
		return empty, ErrPriorityQueueEmpty
	}

	item := heap.Pop(p.items).(*PriorityItem[T])
	delete(p.cache, item.value)

	return item.value, nil
}

// SetPriority updates the priority of a given item in the queue.
// If the item does not exist, this operation does nothing.
//
// Example:
//
//	pq := NewPriorityQueue[string]()
//	pq.Push("task1", 2.0)
//	pq.SetPriority("task1", 1.0)
func (p *PriorityQueue[T]) SetPriority(item T, priority float64) {
	i, ok := p.cache[item]
	if !ok {
		return
	}

	i.priority = priority
	heap.Fix(p.items, i.index)
}

// minHeap is a minimum binary heap that implements heap.Interface.
type minHeap[T comparable] []*PriorityItem[T]

// Len returns the number of items in the heap.
func (m *minHeap[T]) Len() int {
	return len(*m)
}

// Less determines the heap order based on the priority value.
func (m *minHeap[T]) Less(i, j int) bool {
	return (*m)[i].priority < (*m)[j].priority
}

// Swap swaps two items in the heap.
func (m *minHeap[T]) Swap(i, j int) {
	(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
	(*m)[i].index = i
	(*m)[j].index = j
}

// Push adds an item to the heap.
func (m *minHeap[T]) Push(item interface{}) {
	i := item.(*PriorityItem[T])
	i.index = len(*m)
	*m = append(*m, i)
}

// Pop removes and returns the last item in the heap.
func (m *minHeap[T]) Pop() interface{} {
	old := *m
	item := old[len(old)-1]
	*m = old[:len(old)-1]

	return item
}

// Stack is a generic Stack implementation with constant-time membership checks.
type Stack[T comparable] struct {
	registry map[T]struct{}
	elements []T
}

// NewStack creates and returns a new empty Stack.
//
// Example:
//
//	s := NewStack[int]()
//	s.Push(10)
//	s.Push(20)
func NewStack[T comparable]() *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0),
		registry: make(map[T]struct{}),
	}
}

// Push adds an item to the Stack.
//
// Example:
//
//	s := NewStack[int]()
//	s.Push(10)
func (s *Stack[T]) Push(t T) {
	s.elements = append(s.elements, t)
	s.registry[t] = struct{}{}
}

// Pop removes and returns the Top item from the Stack.
// Returns an error if the Stack is empty.
//
// Example:
//
//	s := NewStack[int]()
//	s.Push(10)
//	val, ok := s.Pop()
//	if ok {
//	    fmt.Println(val) // Output: 10
//	}
func (s *Stack[T]) Pop() (T, bool) {
	element, ok := s.Top()
	if !ok {
		return element, false
	}

	s.elements = s.elements[:len(s.elements)-1]
	delete(s.registry, element)

	return element, true
}

// Top returns the Top item from the Stack without removing it.
// Returns false if the Stack is empty.
//
// Example:
//
//	s := NewStack[int]()
//	s.Push(10)
//	val, ok := s.Top()
//	fmt.Println(val, ok) // Output: 10 true
func (s *Stack[T]) Top() (T, bool) {
	if s.IsEmpty() {
		var defaultValue T
		return defaultValue, false
	}

	return s.elements[len(s.elements)-1], true
}

// IsEmpty checks if the Stack is empty.
//
// Example:
//
//	s := NewStack[int]()
//	fmt.Println(s.IsEmpty()) // Output: true
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// ForEach executes a given function for each element in the Stack.
//
// Example:
//
//	s := NewStack[int]()
//	s.Push(10)
//	s.ForEach(func(i int) { fmt.Println(i) }) // Output: 10
func (s *Stack[T]) ForEach(f func(T)) {
	for _, e := range s.elements {
		f(e)
	}
}

// Contains checks if the Stack Contains a specific element.
//
// Example:
//
//	s := NewStack[int]()
//	s.Push(10)
//	fmt.Println(s.Contains(10)) // Output: true
func (s *Stack[T]) Contains(element T) bool {
	_, ok := s.registry[element]
	return ok
}

// StackOfStacks is a Stack containing multiple stacks.
type StackOfStacks[T comparable] struct {
	stacks []*Stack[T]
}

// NewStackOfStacks creates and returns a new empty Stack of stacks.
//
// Example:
//
//	sos := NewStackOfStacks[int]()
func NewStackOfStacks[T comparable]() *StackOfStacks[T] {
	return &StackOfStacks[T]{
		stacks: make([]*Stack[T], 0),
	}
}

// Push adds a new Stack to the Stack of stacks.
//
// Example:
//
//	sos := NewStackOfStacks[int]()
//	s := NewStack[int]()
//	sos.Push(s)
func (s *StackOfStacks[T]) Push(stack *Stack[T]) {
	s.stacks = append(s.stacks, stack)
}

// Pop removes and returns the Top Stack from the Stack of stacks.
// Returns an error if the Stack of stacks is empty.
//
// Example:
//
//	sos := NewStackOfStacks[int]()
//	s := NewStack[int]()
//	sos.Push(s)
//	topStack, err := sos.Pop()
func (s *StackOfStacks[T]) Pop() (*Stack[T], error) {
	e, err := s.Top()
	if err != nil {
		return &Stack[T]{}, err
	}

	s.stacks = s.stacks[:len(s.stacks)-1]
	return e, nil
}

// Top returns the Top Stack from the Stack of stacks without removing it.
// Returns an error if the Stack of stacks is empty.
//
// Example:
//
//	sos := NewStackOfStacks[int]()
//	s := NewStack[int]()
//	sos.Push(s)
//	topStack, err := sos.Top()
func (s *StackOfStacks[T]) Top() (*Stack[T], error) {
	if s.IsEmpty() {
		return &Stack[T]{}, ErrStackEmpty
	}

	return s.stacks[len(s.stacks)-1], nil
}

// IsEmpty checks if the Stack of stacks is empty.
//
// Example:
//
//	sos := NewStackOfStacks[int]()
//	fmt.Println(sos.IsEmpty()) // Output: true
func (s *StackOfStacks[T]) IsEmpty() bool {
	return len(s.stacks) == 0
}
