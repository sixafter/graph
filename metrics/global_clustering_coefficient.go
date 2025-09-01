// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// GlobalClusteringCoefficient calculates the global clustering coefficient of the given graph.
//
// The global clustering coefficient is a measure of the degree to which nodes in a graph
// tend to cluster together. It is defined as the ratio of the number of closed triplets
// (triangles) to the total number of triplets (both open and closed) in the graph.
//
// Parameters:
//   - g: A graph.Interface representing the graph. The graph must provide the necessary
//     methods to traverse edges and count triangles or triplets.
//
// Returns:
//   - A float64 representing the global clustering coefficient, which is a value between 0 and 1.
//   - An error if the computation encounters an issue, such as a graph structure that prevents
//     the calculation.
//
// Type Parameters:
// - K: The type of the graph's vertex keys, which must implement the graph.Ordered interface.
// - T: The type of the graph's vertex data, which can be any type.
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
