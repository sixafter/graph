// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// Transitivity calculates the transitivity of the given graph.
func Transitivity[K graph.Ordered, T any](g graph.Interface[K, T]) (float64, error) {
	if g == nil {
		return 0, graph.ErrNilInputGraph
	}

	// Retrieve the adjacency map
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve adjacency map: %w", err)
	}

	traits := g.Traits()
	var triangles, triplets float64

	// Iterate over all vertices
	for v, neighbors := range adjMap {
		neighborList := make([]K, 0, len(neighbors))
		for neighbor := range neighbors {
			neighborList = append(neighborList, neighbor)
		}

		// Count triplets and triangles for this vertex
		for i := 0; i < len(neighborList); i++ {
			for j := i + 1; j < len(neighborList); j++ {
				triplets++
				n1 := neighborList[i]
				n2 := neighborList[j]

				// Check for triangles
				if traits.IsDirected {
					// Directed graph: check edge directions
					if adjMap[n1][n2] != nil && adjMap[n2][v] != nil && adjMap[v][n1] != nil {
						triangles++
					}
				} else {
					// Undirected graph: ensure consistent counting
					if _, exists := adjMap[n1][n2]; exists {
						triangles += 1.0 / 3.0 // Avoid triple counting
					}
				}
			}
		}
	}

	// Handle cases with no triplets
	if triplets == 0 {
		return 0, nil
	}

	// Return transitivity
	return (3 * triangles) / triplets, nil
}
