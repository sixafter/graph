// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// GlobalClusteringCoefficient calculates the global clustering coefficient of the given graph.
func GlobalClusteringCoefficient[K graph.Ordered, T any](g graph.Interface[K, T]) (float64, error) {
	if g == nil {
		return 0, graph.ErrNilInputGraph
	}

	traits := g.Traits()

	// Get adjacency and predecessor maps for directed graphs
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve adjacency map: %w", err)
	}

	var closedTriplets, totalTriplets float64

	for _, neighbors := range adjMap {
		neighborList := make([]K, 0, len(neighbors))
		for neighbor := range neighbors {
			neighborList = append(neighborList, neighbor)
		}

		numNeighbors := len(neighborList)
		if numNeighbors < 2 {
			continue // No triplets possible
		}

		// Count triplets
		for i := 0; i < numNeighbors; i++ {
			for j := i + 1; j < numNeighbors; j++ {
				totalTriplets++
				if traits.IsDirected {
					// For directed graphs, check both directions between neighbors
					if _, exists := adjMap[neighborList[i]][neighborList[j]]; exists {
						closedTriplets++
					} else if _, exists := adjMap[neighborList[j]][neighborList[i]]; exists {
						closedTriplets++
					}
				} else {
					// For undirected graphs, only check one direction
					if _, exists := adjMap[neighborList[i]][neighborList[j]]; exists {
						closedTriplets++
					}
				}
			}
		}
	}

	if totalTriplets == 0 {
		return 0, nil // Avoid division by zero
	}

	return closedTriplets / totalTriplets, nil
}
