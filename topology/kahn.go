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
// Complexity: O(Items + E), where Items is the number of vertices and E is the number of edges.
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

		edges := adjacencyMap[currentVertex]

		for target := range edges {
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

// TopologicalSortDeterministic performs a deterministic topological sort on a directed graph
// and is an enhanced version of Kahn's algorithm that guarantees a deterministic topological
// order by adding a sorting mechanism.
//
// The function ensures that the resulting order is consistent across executions by applying
// a custom comparison function (`less`) to resolve ties between vertices.
//
// Parameters:
//   - g: A graph implementing the `graph.Interface[K, T]`, representing the directed graph to sort.
//   - less: A comparison function that determines the ordering of vertices with the same topological rank.
//
// Returns:
//   - A slice of vertices (`[]K`) in topological order.
//   - An error if the graph is not directed or contains cycles.
//
// Errors:
//   - Returns `graph.ErrUndirectedGraph` if the graph is undirected.
//   - Returns `graph.ErrCyclicGraph` if the graph contains a cycle.
//   - Returns `graph.ErrFailedToGetGraphOrder` or `graph.ErrFailedToGetAdjacencyMap` if there is an issue retrieving graph properties.
//
// Key Details:
//   - The function starts by identifying all vertices with no predecessors and adding them to a processing queue (`q`).
//   - During each iteration, the next vertex is dequeued, added to the result, and its successors are checked.
//   - The `frontier` variable holds the successors of the current vertex that are ready for processing (i.e., all their predecessors have been processed).
//   - The `frontier` is sorted using the `less` function before being added to the processing queue to ensure deterministic ordering.
//
// Time Complexity:
//   - O(Items + E) for traversing the graph, where Items is the number of vertices and E is the number of edges.
//   - Additional cost for sorting the `frontier` in each iteration, which depends on the graph structure and the `less` function.
//
// Example Usage:
//
//	g := NewDirectedGraph[int, any]()
//	g.AddEdge(1, 2)
//	g.AddEdge(2, 3)
//	order, err := TopologicalSortDeterministic(g, func(a, b int) bool { return a < b })
//	if err != nil {
//	    fmt.Println(err)
//	} else {
//	    fmt.Println(order) // Output: [1, 2, 3]
//
// Frontier Explanation:
//   - The `frontier` holds the set of vertices that are ready to be processed next. A vertex is added to the
//     `frontier` only when all its predecessors have been processed and added to the result.
//   - Sorting the `frontier` ensures a consistent and deterministic order of processing, even when the graph
//     has multiple valid topological orderings.
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
		edges := adjacencyMap[currentVertex]

		for target := range edges {
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
