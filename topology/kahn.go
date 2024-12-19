// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package topology

import (
	"fmt"
	"sort"

	"github.com/sixafter/graph"
)

// TopologicalSort computes a topological ordering of the vertices in a DirectedGraph acyclic graph (DAG).
// A topological ordering ensures that for every DirectedGraph edge (A, B), vertex A appears before vertex B.
//
// The function implements Kahn's algorithm non-recursively. If there are multiple valid
// topological orderings, an arbitrary one is returned. To produce deterministic results,
// consider using [TopologicalSortDeterministic].
//
// Parameters:
//   - g: A DirectedGraph acyclic graph represented as a [Interface[K, T]].
//
// Returns:
//   - A slice of vertex hashes (of type K) representing the topological order.
//   - An error if the graph is Undirected, Contains cycles, or encounters other failures.
//
// Errors:
//   - [ErrUndirectedGraph] if the graph is not DirectedGraph.
//   - [ErrCyclicGraph] if the graph Contains cycles.
//   - [ErrFailedToGetGraphOrder], [ErrFailedToGetAdjacencyMap], or [ErrFailedToGetPredecessorMap] for failures
//     in retrieving the graph's properties.
//
// Complexity: O(V + E), where V is the number of vertices and E is the number of edges.
func TopologicalSort[K graph.Ordered, T any](g graph.Interface[K, T]) ([]K, error) {
	if !g.Traits().IsDirected {
		return nil, graph.ErrUndirectedGraph
	}

	gOrder, err := g.Order()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetGraphOrder, err)
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetAdjacencyMap, err)
	}

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetPredecessorMap, err)
	}

	q := make([]K, 0)

	for vertex, predecessors := range predecessorMap {
		if len(predecessors) == 0 {
			q = append(q, vertex)
			delete(predecessorMap, vertex)
		}
	}

	order := make([]K, 0, gOrder)

	for len(q) > 0 {
		currentVertex := q[0]
		q = q[1:]

		order = append(order, currentVertex)

		edgeMap := adjacencyMap[currentVertex]

		for target := range edgeMap {
			predecessors := predecessorMap[target]
			delete(predecessors, currentVertex)

			if len(predecessors) == 0 {
				q = append(q, target)
				delete(predecessorMap, target)
			}
		}
	}

	if len(order) != gOrder {
		return nil, graph.ErrCyclicGraph
	}

	return order, nil
}

// TopologicalSortDeterministic computes a deterministic topological ordering of the vertices in a
// Directed Acyclic Interface (DAG) using Kahn's algorithm. A custom comparison function `less` is
// used to deterministically decide when multiple valid orders exist.
//
// Parameters:
//   - g: A DirectedGraph acyclic graph represented as a [Interface[K, T]].
//   - less: A comparison function to impose an ordering on vertices with equal precedence.
//
// Returns:
//   - A slice of vertex hashes (of type K) in deterministic topological order.
//   - An error if the graph is Undirected, Contains cycles, or encounters other failures.
//
// Errors:
//   - [ErrUndirectedGraph] if the graph is not DirectedGraph.
//   - [ErrCyclicGraph] if the graph Contains cycles.
//   - [ErrFailedToGetGraphOrder], [ErrFailedToGetAdjacencyMap], or [ErrFailedToGetPredecessorMap] for failures
//     in retrieving the graph's properties.
//
// Complexity: O((V + E) log V), where V is the number of vertices and E is the number of edges.
// The additional log V factor arises from sorting the queue of vertices with zero in-degree.
func TopologicalSortDeterministic[K graph.Ordered, T any](g graph.Interface[K, T], less func(K, K) bool) ([]K, error) {
	if !g.Traits().IsDirected {
		return nil, graph.ErrUndirectedGraph
	}

	gOrder, err := g.Order()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetGraphOrder, err)
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetAdjacencyMap, err)
	}

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetPredecessorMap, err)
	}

	q := make([]K, 0)

	for vertex, predecessors := range predecessorMap {
		if len(predecessors) == 0 {
			q = append(q, vertex)
			delete(predecessorMap, vertex)
		}
	}

	order := make([]K, 0, gOrder)

	// Initial sort of the queue to ensure deterministic ordering.
	sort.Slice(q, func(i, j int) bool {
		return less(q[i], q[j])
	})

	for len(q) > 0 {
		currentVertex := q[0]
		q = q[1:]

		order = append(order, currentVertex)

		frontier := make([]K, 0)

		edgeMap := adjacencyMap[currentVertex]

		for target := range edgeMap {
			predecessors := predecessorMap[target]
			delete(predecessors, currentVertex)

			if len(predecessors) == 0 {
				frontier = append(frontier, target)
				delete(predecessorMap, target)
			}
		}

		// Sort the frontier to maintain deterministic ordering.
		sort.Slice(frontier, func(i, j int) bool {
			return less(frontier[i], frontier[j])
		})

		q = append(q, frontier...)
	}

	if len(order) != gOrder {
		return nil, graph.ErrCyclicGraph
	}

	return order, nil
}
