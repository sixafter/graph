// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package traverse

import (
	"fmt"
	"sort"

	"github.com/sixafter/graph"
)

// BFS performs a breadth-first search on the graph, starting from the given vertex. The visit
// function will be invoked with the hash of the vertex currently visited. If it returns false, BFS
// will continue traversing the graph, and if it returns true, the traversal will be stopped. In
// case the graph is disconnected, only the vertices joined with the starting vertex are visited.
//
// This example prints all vertices of the graph in BFS-order:
//
//	g := graph.New(graph.IntHash)
//
//	_ = g.AddVertex(1)
//	_ = g.AddVertex(2)
//	_ = g.AddVertex(3)
//
//	_ = g.AddEdge(1, 2)
//	_ = g.AddEdge(2, 3)
//	_ = g.AddEdge(3, 1)
//
//	_ = graph.BFS(g, 1, func(value int) bool {
//		fmt.Println(value)
//		return false
//	})
//
// Similarly, if you have a graph of City vertices and the traversal should stop at London, the
// visit function would look as follows:
//
//	func(c City) bool {
//		return c.Name == "London"
//	}
//
// BFS is non-recursive and maintains a Stack instead.
func BFS[K graph.Ordered, T any](g graph.Interface[K, T], start K, visit func(K) bool) error {
	ignoreDepth := func(vertex K, _ int) bool {
		return visit(vertex)
	}
	return BFSWithDepthTracking(g, start, ignoreDepth)
}

// BFSWithDepthTracking performs a breadth-first search (BFS) on the graph, starting from the given vertex.
// The `visit` function is invoked for each vertex visited, and the current depth is passed as the
// second argument. The traversal continues unless the `visit` function returns `true`, in which case
// the traversal stops early.
//
// Parameters:
// - g: The graph on which BFS is performed.
// - start: The starting vertex for BFS.
// - visit: A function that is called for each visited vertex. If `visit` returns `true`, BFS stops.
//
// Example usage:
//
//	_ = graph.BFSWithDepthTracking(g, 1, func(value int, depth int) bool {
//	    fmt.Printf("Visited vertex: %d at depth: %d\n", value, depth)
//	    return depth > 3 // Stop traversal if depth exceeds 3
//	})
//
// Errors:
// - Returns an error if the adjacency map cannot be retrieved.
// - Returns an error if the start vertex is not found in the graph.
//
// Complexity: O(Items + E), where Items is the number of vertices and E is the number of edges.
func BFSWithDepthTracking[K graph.Ordered, T any](g graph.Interface[K, T], start K, visit func(K, int) bool) error {
	// Retrieve the adjacency map of the graph.
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map: %w", err)
	}

	// Ensure the starting vertex exists in the graph.
	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	// Define a helper type for the BFS queue, which stores the vertex and its depth.
	type queueNode struct {
		vertex K   // The current vertex.
		depth  int // The depth of the vertex from the starting point.
	}

	// Initialize the queue with the starting vertex at depth 1.
	q := []queueNode{{vertex: start, depth: 1}}

	// Create a map to track visited vertices, starting with the initial vertex.
	visited := make(map[K]bool)
	visited[start] = true

	// Process the queue until it is empty.
	for len(q) > 0 {
		// Dequeue the first element in the queue.
		current := q[0]
		q = q[1:]

		// Call the visit function for the current vertex. If it returns true, stop traversal.
		if stop := visit(current.vertex, current.depth); stop {
			return nil
		}

		// Gather and sort neighbors for deterministic processing.
		neighbors := make([]K, 0)
		for neighbor := range adjacencyMap[current.vertex] {
			if !visited[neighbor] {
				neighbors = append(neighbors, neighbor)
			}
		}
		// Sort neighbors to ensure deterministic order.
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i] < neighbors[j]
		})

		// Add sorted neighbors to the queue.
		for _, neighbor := range neighbors {
			visited[neighbor] = true
			q = append(q, queueNode{vertex: neighbor, depth: current.depth + 1})
		}
	}

	// Return nil to indicate that the traversal completed successfully.
	return nil
}
