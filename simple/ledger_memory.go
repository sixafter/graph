// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"fmt"
	"sort"
	"sync"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// memoryLedger is an in-memory implementation of the ledger interface,
// providing methods to manage the vertices and edges of a graph.
//
// This implementation uses maps to efficiently ledger and retrieve graph data,
// including vertex data, vertex properties, and both incoming and outgoing edges.
//
// Type Parameters:
//   - K: The type used to uniquely identify vertices (e.g., string, int). Must be comparable.
//   - T: The type of data stored in each vertex (e.g., a custom struct or primitive type).
type memoryLedger[K graph.Ordered, T any] struct {
	// vertices maps vertex identifiers (keys of type K) to their associated data (of type T).
	vertices map[K]T

	// vertexProps maps vertex identifiers to their associated properties,
	// such as metadata or attributes specific to each vertex.
	vertexProps map[K]graph.VertexProperties

	// outEdges maps each vertex to a map of its outgoing edges.
	// The inner map maps the target vertex identifiers to their corresponding Edge objects.
	outEdges map[K]map[K]graph.Edge[K]

	// inEdges maps each vertex to a map of its incoming edges.
	// The inner map maps the source vertex identifiers to their corresponding Edge objects.
	inEdges map[K]map[K]graph.Edge[K]

	// lock is a read-write mutex used to ensure thread-safe access to the graph's data.
	lock sync.RWMutex

	// edgeCount tracks the total number of edges in the graph,
	// including both incoming and outgoing edges.
	edgeCount int
}

// newMemoryStore initializes a new in-memory graph ledger.
//
// Returns:
//   - A pointer to a memoryLedger instance.
func newMemoryStore[K graph.Ordered, T any]() (ledger[K, T], error) {
	return &memoryLedger[K, T]{
		vertices:    make(map[K]T),
		vertexProps: make(map[K]graph.VertexProperties),
		outEdges:    make(map[K]map[K]graph.Edge[K]),
		inEdges:     make(map[K]map[K]graph.Edge[K]),
	}, nil
}

// New creates and returns a new instance of the ledger implementation.
// This method allows for creating fresh, independent ledger instances that adhere to the ledger interface.
func (ms *memoryLedger[K, T]) New() (ledger[K, T], error) {
	return newMemoryStore[K, T]()
}

// AddVertex adds a vertex with the specified hash, value, and properties to the graph.
// If the vertex already exists, ErrVertexAlreadyExists is returned.
func (ms *memoryLedger[K, T]) AddVertex(hash K, value T, properties graph.VertexProperties) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if _, exists := ms.vertices[hash]; exists {
		return graph.ErrVertexAlreadyExists
	}

	ms.vertices[hash] = value
	ms.vertexProps[hash] = properties
	return nil
}

// AddEdge adds an edge between the specified source and target vertices.
// If either vertex does not exist, ErrVertexNotFound is returned.
func (ms *memoryLedger[K, T]) AddEdge(source, target K, edge graph.Edge[K]) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if _, exists := ms.vertices[source]; !exists {
		return graph.ErrVertexNotFound
	}
	if _, exists := ms.vertices[target]; !exists {
		return graph.ErrVertexNotFound
	}

	if ms.outEdges[source] == nil {
		ms.outEdges[source] = make(map[K]graph.Edge[K])
	}
	if ms.inEdges[target] == nil {
		ms.inEdges[target] = make(map[K]graph.Edge[K])
	}

	ms.outEdges[source][target] = edge
	ms.inEdges[target][source] = edge
	ms.edgeCount++
	return nil
}

// FindVertex retrieves a vertex and its properties by its hash.
// If the vertex does not exist, ErrVertexNotFound is returned.
func (ms *memoryLedger[K, T]) FindVertex(key K) (T, graph.VertexProperties, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	vertex, exists := ms.vertices[key]
	if !exists {
		return *new(T), nil, graph.ErrVertexNotFound
	}

	return vertex, ms.vertexProps[key], nil
}

// FindEdge retrieves an edge between the specified source and target vertices.
// If the edge does not exist, ErrEdgeNotFound is returned.
func (ms *memoryLedger[K, T]) FindEdge(source, target K) (graph.Edge[K], error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	if edge, exists := ms.outEdges[source][target]; exists {
		return edge, nil
	}
	return nil, graph.ErrEdgeNotFound
}

// ModifyVertex updates the properties of an existing vertex.
// If the vertex does not exist, ErrVertexNotFound is returned.
func (ms *memoryLedger[K, T]) ModifyVertex(key K, properties graph.VertexProperties) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if _, exists := ms.vertices[key]; !exists {
		return graph.ErrVertexNotFound
	}

	ms.vertexProps[key] = properties
	return nil
}

// ModifyEdge updates the properties of an existing edge between the specified source and target vertices.
// If the edge does not exist, ErrEdgeNotFound is returned.
func (ms *memoryLedger[K, T]) ModifyEdge(source, target K, edge graph.Edge[K]) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if _, exists := ms.outEdges[source][target]; !exists {
		return graph.ErrEdgeNotFound
	}

	ms.outEdges[source][target] = edge
	ms.inEdges[target][source] = edge
	return nil
}

// RemoveVertex removes a vertex by its hash. If the vertex is connected by any edges,
// ErrVertexHasEdges is returned. If the vertex does not exist, ErrVertexNotFound is returned.
func (ms *memoryLedger[K, T]) RemoveVertex(key K) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if _, exists := ms.vertices[key]; !exists {
		return graph.ErrVertexNotFound
	}

	if len(ms.outEdges[key]) > 0 || len(ms.inEdges[key]) > 0 {
		return graph.ErrVertexHasEdges
	}

	delete(ms.vertices, key)
	delete(ms.vertexProps, key)
	delete(ms.outEdges, key)
	delete(ms.inEdges, key)
	return nil
}

// RemoveEdge removes an edge between the specified source and target vertices.
// If the edge does not exist, ErrEdgeNotFound is returned.
func (ms *memoryLedger[K, T]) RemoveEdge(source, target K) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	if _, exists := ms.outEdges[source][target]; !exists {
		return graph.ErrEdgeNotFound
	}

	delete(ms.outEdges[source], target)
	delete(ms.inEdges[target], source)
	ms.edgeCount--
	return nil
}

// ListVertices retrieves all vertex hashes in the graph.
func (ms *memoryLedger[K, T]) ListVertices() ([]K, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	keys := make([]K, 0, len(ms.vertices))
	for key := range ms.vertices {
		keys = append(keys, key)
	}

	// Sort hashes to ensure consistent ordering
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j] // Assumes K supports < via Ordered[K]
	})

	return keys, nil
}

// ListEdges retrieves all edges in the graph as a slice.
func (ms *memoryLedger[K, T]) ListEdges() ([]graph.Edge[K], error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	allEdges := make([]graph.Edge[K], 0, ms.edgeCount)
	for _, targets := range ms.outEdges {
		for _, edge := range targets {
			allEdges = append(allEdges, edge)
		}
	}

	// Sort edges by source, then by target
	sort.Slice(allEdges, func(i, j int) bool {
		if allEdges[i].Source() == allEdges[j].Source() {
			return allEdges[i].Target() < allEdges[j].Target()
		}
		return allEdges[i].Source() < allEdges[j].Source()
	})

	return allEdges, nil
}

// CountVertices returns the total number of vertices in the ledger.
func (ms *memoryLedger[K, T]) CountVertices() (int, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	return len(ms.vertices), nil
}

// CountEdges returns the total number of edges in the ledger.
func (ms *memoryLedger[K, T]) CountEdges() (int, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	return ms.edgeCount, nil
}

// WouldCreateCycle checks if adding an edge from source to target would create a cycle in the graph.
// It requires access to both vertex and edge data to verify vertex existence and traverse edges.
//
// Parameters:
//   - source: The unique identifier of the source vertex.
//   - target: The unique identifier of the target vertex.
//
// Returns:
//   - true if adding the edge would create a cycle.
//   - false if no cycle would be created.
//   - An error if either the source or target vertex does not exist.
func (ms *memoryLedger[K, T]) WouldCreateCycle(source, target K) (bool, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	// Verify that both the source and target vertices exist
	if _, exists := ms.vertices[source]; !exists {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", source, graph.ErrVertexNotFound)
	}
	if _, exists := ms.vertices[target]; !exists {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", target, graph.ErrVertexNotFound)
	}

	// Check for trivial cycle (self-loop)
	if source == target {
		return true, nil
	}

	// Perform a depth-first search (DFS) using inEdges to detect cycles
	stack := queue.NewStack[K]()
	visited := make(map[K]struct{})

	stack.Push(source)

	for !stack.IsEmpty() {
		current, _ := stack.Pop()

		if _, ok := visited[current]; !ok {
			// If the target is reachable from the source, adding the edge creates a cycle
			if current == target {
				return true, nil
			}

			visited[current] = struct{}{}

			// Traverse adjacent vertices using the inEdges map
			if neighbors, exists := ms.inEdges[current]; exists {
				for neighbor := range neighbors {
					stack.Push(neighbor)
				}
			}
		}
	}

	return false, nil
}
