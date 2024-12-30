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
//
// It computes the average of the shortest path lengths between all pairs of vertices
// in the graph. If the graph is disconnected, the function returns an error as the
// calculation is not well-defined in such cases.
//
// Parameters:
//   - g: A graph.Interface representing the graph. The graph must implement the
//     necessary methods for traversal and shortest path calculation.
//
// Returns:
// - A float64 representing the average shortest path length.
// - An error if the graph is disconnected or any other issue arises during computation.
//
// Type Parameters:
// - K: The type of the graph's vertex keys, which must implement the graph.Ordered interface.
// - T: The type of the graph's vertex data, which can be any type.
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
