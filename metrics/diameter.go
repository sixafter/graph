// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"errors"
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/paths"
)

// Diameter calculates the diameter of the given graph.
func Diameter[K graph.Ordered, T any](g graph.Interface[K, T]) (int, error) {
	if g == nil {
		return 0, graph.ErrNilInputGraph
	}

	// Retrieve all vertices
	vertices, err := g.Vertices()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve vertices: %w", err)
	}

	vertexCount := len(vertices)
	if vertexCount == 0 {
		return 0, nil // Diameter of an empty graph is 0
	}

	maxShortestPath := 0

	// Iterate over all pairs of vertices
	for i := 0; i < vertexCount; i++ {
		src := vertices[i].ID()

		for j := 0; j < vertexCount; j++ {
			if i == j {
				continue // Skip self-loops
			}
			target := vertices[j].ID()

			// Find the shortest path using DijkstraFrom
			path, err := paths.DijkstraFrom(g, src, target)
			if err != nil {
				if errors.Is(err, graph.ErrTargetNotReachable) {
					return 0, err
				}

				return 0, fmt.Errorf("failed to calculate shortest path from %v to %v: %w", src, target, err)
			}

			// Calculate the length of the path
			pathLength := len(path) - 1 // Number of edges in the path
			if pathLength > maxShortestPath {
				maxShortestPath = pathLength
			}
		}
	}

	return maxShortestPath, nil
}
