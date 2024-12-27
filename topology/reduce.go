// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package topology

import (
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// TransitiveReduction computes the transitive reduction of a DirectedGraph graph.
// A transitive reduction retains the same reachability between vertices while minimizing the number of edges.
//
// The function identifies and removes redundant edges by checking for alternative paths
// between pairs of vertices.
//
// Parameters:
//   - g: A DirectedGraph graph represented as a [Interface[K, T]].
//
// Returns:
//   - A new graph that is the transitive reduction of the input graph.
//   - An error if the graph is Undirected, Contains cycles, or encounters other failures.
//
// Errors:
//   - [ErrUndirectedGraph] if the graph is not DirectedGraph.
//   - [ErrCyclicGraph] if the graph Contains cycles.
//   - [ErrFailedToCloneGraph], [ErrFailedToGetGraphOrder], or [ErrFailedToGetAdjacencyMap] for failures
//     in graph operations.
//
// Complexity: O(V * (V + E)), where V is the number of vertices and E is the number of edges.
// This makes it computationally expensive for large graphs.
//
// Example:
//
//	g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())
//	g.AddEdge(1, 2)
//	g.AddEdge(2, 3)
//	g.AddEdge(1, 3) // Redundant edge
//
//	reducedGraph, err := TransitiveReduction(g)
//	if err != nil {
//	    log.Fatalf("Failed to perform transitive reduction: %v", err)
//	}
//
//	// `reducedGraph` will no longer contain the redundant edge (1 -> 3).
func TransitiveReduction[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	if !g.Traits().IsDirected {
		return nil, graph.ErrUndirectedGraph
	}

	reduced, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToCloneGraph, err)
	}

	adjacencyMap, err := reduced.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetAdjacencyMap, err)
	}

	// Iterate over each vertex and its direct successors in the adjacency map.
	// For each vertex, perform the following steps:
	// 1. Retrieve the graph's order (number of vertices) to initialize necessary data structures.
	// 2. For each direct successor of the current vertex:
	//    a. Use a Stack-based depth-first search (DFS) to explore all vertices reachable from the successor.
	//    b. Maintain a `visited` map to track explored vertices and prevent revisiting.
	//    c. Push the successor onto the Stack and continue DFS until the Stack is empty.
	//
	// Inside the DFS loop:
	// - Check if the current vertex has already been visited to avoid redundant processing.
	// - For each adjacency (neighbor) of the current vertex:
	//   a. If the adjacency has already been visited and is also present in the Stack,
	//      a cycle exists in the graph, and the transitive reduction cannot proceed.
	//   b. If the adjacency has not been visited and there is a direct edge between the Top-level
	//      vertex and the adjacency, this edge is redundant and is removed from the graph.
	//   c. Push the adjacency onto the Stack for further exploration.
	//
	// This portion of the code ensures that for each vertex, edges that are reachable via
	// transitive paths are identified and removed, effectively performing the transitive reduction.
	//
	// Error Handling:
	// - If retrieving the graph order fails, return an error indicating the failure.
	// - If a cycle is detected during the DFS, return an error indicating that transitive
	//   reduction cannot be performed on graphs with cycles.
	for vertex, successors := range adjacencyMap {
		order, err := reduced.Order()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetGraphOrder, err)
		}

		for successor := range successors {
			s := queue.NewStack[K]()
			visited := make(map[K]struct{}, order)

			s.Push(successor)

			for !s.IsEmpty() {
				current, _ := s.Pop()

				if _, ok := visited[current]; ok {
					continue
				}

				visited[current] = struct{}{}
				s.Push(current)

				for adjacency := range adjacencyMap[current] {
					if _, ok := visited[adjacency]; ok {
						if s.Contains(adjacency) {
							// If the current adjacency is both on the Stack and
							// has already been visited, there is a cycle.
							return nil, graph.ErrCyclicGraph
						}
						continue
					}

					if _, ok := adjacencyMap[vertex][adjacency]; ok {
						_ = reduced.RemoveEdge(vertex, adjacency)
					}
					s.Push(adjacency)
				}
			}
		}
	}

	return reduced, nil
}
