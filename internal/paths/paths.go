// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package paths

import (
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// WouldCreateCycle determines whether adding an edge between the two given vertices
// would introduce a cycle in the graph. WouldCreateCycle will not create an edge.
//
// A potential edge would create a cycle if the target vertex is also a parent
// of the source vertex. To determine this, WouldCreateCycle runs a Depth-First Search (DFS).
//
// Returns true if adding the edge would introduce a cycle, otherwise false.
//
// Example:
//
//	isCycle, err := WouldCreateCycle(graph, "A", "B")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if isCycle {
//		fmt.Println("Adding the edge would create a cycle")
//	} else {
//		fmt.Println("Adding the edge would not create a cycle")
//	}
func WouldCreateCycle[K graph.Ordered, T any](g graph.Interface[K, T], source, target K) (bool, error) {
	if _, err := g.Vertex(source); err != nil {
		return false, fmt.Errorf("%w: %v", graph.ErrVertexNotFound, source)
	}

	if _, err := g.Vertex(target); err != nil {
		return false, fmt.Errorf("%w: %v", graph.ErrVertexNotFound, target)
	}

	if source == target {
		return true, graph.ErrSameSourceAndTarget
	}

	predecessors, err := g.PredecessorMap()
	if err != nil {
		return false, fmt.Errorf("%w: %v", graph.ErrPredecessorMapFailed, err)
	}

	s := queue.NewStack[K]()
	visited := make(map[K]bool)

	s.Push(source)

	for !s.IsEmpty() {
		current, _ := s.Pop()

		if _, ok := visited[current]; !ok {
			if current == target {
				return true, nil
			}

			visited[current] = true

			for adjacency := range predecessors[current] {
				s.Push(adjacency)
			}
		}
	}

	return false, nil
}
