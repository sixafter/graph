// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package traverse

import (
	"fmt"

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
	return BFSWithDepth(g, start, ignoreDepth)
}

// BFSWithDepth performs a breadth-first search (BFS) on the graph, starting from the given vertex.
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
//	_ = graph.BFSWithDepth(g, 1, func(value int, depth int) bool {
//	    fmt.Printf("Visited vertex: %d at depth: %d\n", value, depth)
//	    return depth > 3 // Stop traversal if depth exceeds 3
//	})
//
// Errors:
// - Returns an error if the adjacency map cannot be retrieved.
// - Returns an error if the start vertex is not found in the graph.
//
// Complexity: O(V + E), where V is the number of vertices and E is the number of edges.
func BFSWithDepth[K graph.Ordered, T any](g graph.Interface[K, T], start K, visit func(K, int) bool) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map: %w", err)
	}

	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	// Define a queue node with vertex and depth
	type queueNode struct {
		vertex K
		depth  int
	}

	q := []queueNode{{vertex: start, depth: 1}}
	visited := make(map[K]bool)
	visited[start] = true

	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		// Visit the current vertex
		if stop := visit(current.vertex, current.depth); stop {
			return nil
		}

		for adjacency := range adjacencyMap[current.vertex] {
			if !visited[adjacency] {
				visited[adjacency] = true
				q = append(q, queueNode{vertex: adjacency, depth: current.depth + 1})
			}
		}
	}

	return nil
}
