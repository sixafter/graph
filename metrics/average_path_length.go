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

// AveragePathLength calculates the average shortest path length in the given graph.
func AveragePathLength[K graph.Ordered, T any](g graph.Interface[K, T]) (float64, error) {
	if g == nil {
		return 0, graph.ErrNilInputGraph
	}

	// Retrieve all vertices
	vertices, err := g.Vertices()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve vertices: %w", err)
	}

	vertexCount := len(vertices)
	if vertexCount < 2 {
		return 0, nil // APL is 0 for graphs with fewer than 2 vertices
	}

	var totalPathLength float64
	var pathCount int

	hashFunc := g.Hash()

	// Iterate over all pairs of vertices
	for i := 0; i < vertexCount; i++ {
		src := hashFunc(vertices[i].Value())

		for j := 0; j < vertexCount; j++ {
			if i == j {
				continue // Skip self-loops
			}

			target := hashFunc(vertices[j].Value())

			// Find the shortest path using DijkstraFrom
			path, err := paths.DijkstraFrom(g, src, target)
			if err != nil {
				if errors.Is(err, graph.ErrTargetNotReachable) {
					return 0, err
				}
				return 0, fmt.Errorf("failed to calculate shortest path from %v to %v: %w", src, target, err)
			}

			// Add path length to total
			totalPathLength += float64(len(path) - 1) // Number of edges in the path
			pathCount++
		}
	}

	if pathCount == 0 {
		return 0, fmt.Errorf("graph is disconnected, no valid paths")
	}

	return totalPathLength / float64(pathCount), nil
}
