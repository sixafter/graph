// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package paths

import (
	"errors"
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// FindAllPaths computes all possible paths between two vertices in a directedGraph graph.
// The function uses a non-recursive Stack-based algorithm to explore all paths from
// the `start` vertex to the `end` vertex, avoiding cycles in the process.
//
// Parameters:
//   - g: The graph of type `Interface[K, T, Items]` in which paths are to be found.
//   - start: The starting vertex, identified by its hash of type K.
//   - end: The ending vertex, identified by its hash of type K.
//
// Returns:
//   - A slice of paths, where each path is represented as a slice of vertex hashes ([]K).
//   - An error if an issue occurs during graph traversal or Stack operations.
//
// Algorithm Details:
//
//	The algorithm computes all paths between two vertices in a graph using a non-recursive,
//	Stack-based approach. It avoids recursion by maintaining two stacks:
//	- The main Stack tracks the current path being explored.
//	- The vice Stack (a Stack of stacks) tracks unexplored neighbors for each vertex in the path.
//
//	At each step, the algorithm either explores a new neighbor by pushing it onto both stacks
//	or backtracks by popping vertices when no neighbors remain. Cycles are avoided by ensuring
//	no vertex is revisited within the same path.
//
//	This method is beneficial for deep graphs, as it avoids Stack overflow issues
//	associated with recursion, and systematically enumerates all valid paths.
//
//	For more detail, see https://boycgit.github.io/all-paths-between-two-vertex/
//
// Constraints:
//   - Cycles in the graph are avoided by ensuring no vertex in the current path appears twice.
//
// Complexity:
//   - Time Complexity: Exponential in the worst case, as it depends on the number of paths.
//     The actual runtime is influenced by the graph's structure and connectivity.
//   - Space Complexity: Proportional to the size of the stacks and the number of stored paths.
//
// Errors:
//   - Returns an error if the adjacency map cannot be retrieved.
//   - Returns an error for internal Stack inconsistencies or empty stacks during traversal.
//
// Example Usage:
//
//	g := simple.New(graph.IntHash, graph.Directed())
//	g.AddEdge(1, 2)
//	g.AddEdge(1, 3)
//	g.AddEdge(2, 4)
//	g.AddEdge(3, 4)
//
//	paths, err := FindAllPaths(g, 1, 4)
//	if err != nil {
//	    log.Fatalf("Error finding paths: %v", err)
//	}
//
//	// Output: [[1, 2, 4], [1, 3, 4]]
func FindAllPaths[K graph.Ordered, T any](g graph.Interface[K, T], start, end K) ([][]K, error) {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	mainStack := queue.NewStack[K]()
	viceStack := queue.NewStackOfStacks[K]()

	checkEmpty := func() error {
		if mainStack.IsEmpty() || viceStack.IsEmpty() {
			return errors.New("empty Stack")
		}
		return nil
	}

	buildLayer := func(element K) {
		mainStack.Push(element)
		newElements := queue.NewStack[K]()

		for e := range adjacencyMap[element] {
			var contains bool
			mainStack.ForEach(func(k K) {
				if e == k {
					contains = true
				}
			})
			if contains {
				continue
			}
			newElements.Push(e)
		}
		viceStack.Push(newElements)
	}

	buildStack := func() error {
		if err = checkEmpty(); err != nil {
			return fmt.Errorf("unable to build Stack: %w", err)
		}

		elements, _ := viceStack.Top()

		for !elements.IsEmpty() {
			element, _ := elements.Pop()
			buildLayer(element)
			elements, _ = viceStack.Top()
		}

		return nil
	}

	removeLayer := func() error {
		if err = checkEmpty(); err != nil {
			return fmt.Errorf("unable to remove layer: %w", err)
		}

		if e, _ := viceStack.Top(); !e.IsEmpty() {
			return errors.New("the Top element of vice-Stack is not empty")
		}

		_, _ = mainStack.Pop()
		_, _ = viceStack.Pop()

		return nil
	}

	buildLayer(start)

	allPaths := make([][]K, 0)

	for !mainStack.IsEmpty() {
		v, _ := mainStack.Top()
		adjacencies, _ := viceStack.Top()

		if adjacencies.IsEmpty() {
			if v == end {
				path := make([]K, 0)
				mainStack.ForEach(func(k K) {
					path = append(path, k)
				})
				allPaths = append(allPaths, path)
			}

			err = removeLayer()
			if err != nil {
				return nil, err
			}
		} else {
			if err = buildStack(); err != nil {
				return nil, err
			}
		}
	}

	return allPaths, nil
}
