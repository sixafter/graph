// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package traverse

import (
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// DFS performs a depth-first traversal on the given graph, starting from the specified vertex.
// The traversal visits each vertex and calls the provided `visit` function. If `visit` returns
// true, the traversal stops early.
//
// Parameters:
// - g: The graph to traverse.
// - start: The starting vertex for the traversal.
// - visit: A callback function invoked for each visited vertex. If it returns true, the traversal stops.
//
// Returns:
// - An error if the adjacency map cannot be retrieved or if the starting vertex is not found.
//
// Complexity: O(V + E), where V is the number of vertices and E is the number of edges.
//
// DFS is non-recursive and maintains a Stack instead.
func DFS[K graph.Ordered, T any](g graph.Interface[K, T], start K, visit func(K) bool) error {
	// Retrieve the adjacency map for the graph, which provides edges for each vertex.
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map: %w", err)
	}

	// Ensure that the starting vertex exists in the graph.
	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	// Initialize a stack to simulate the recursion used in traditional DFS.
	// The stack contains vertices to be processed.
	stack := queue.NewStack[K]()

	// A map to track which vertices have already been visited.
	visited := make(map[K]bool)

	// Push the starting vertex onto the stack to begin the traversal.
	stack.Push(start)

	// Continue the traversal until the stack is empty.
	for !stack.IsEmpty() {
		// Pop the top vertex from the stack to process it.
		current, _ := stack.Pop()

		// If the vertex has not been visited, process it.
		if _, ok := visited[current]; !ok {
			// Call the visit function for the current vertex.
			// If the function returns true, stop the traversal.
			if stop := visit(current); stop {
				break
			}

			// Mark the vertex as visited to prevent re-processing.
			visited[current] = true

			// Push all adjacent vertices of the current vertex onto the stack.
			// These vertices will be processed later.
			for adjacency := range adjacencyMap[current] {
				stack.Push(adjacency)
			}
		}
	}

	// Return nil to indicate that the traversal completed successfully.
	return nil
}
