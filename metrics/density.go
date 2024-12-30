// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// Density calculates the density of the given graph.
//
// The density of a graph is defined as the ratio of the number of edges present
// in the graph to the maximum number of edges possible in a graph with the same
// number of vertices. For directed graphs, the maximum possible edges is n*(n-1),
// and for undirected graphs, it is n*(n-1)/2, where n is the number of vertices.
//
// Parameters:
//   - g: A graph.Interface representing the graph. The graph must provide the necessary
//     methods to retrieve the vertex count and edge count.
//
// Returns:
// - A float64 representing the density of the graph, which is a value between 0 and 1.
// - An error if the calculation cannot be performed (e.g., if the graph has no vertices).
//
// Type Parameters:
// - K: The type of the graph's vertex keys, which must implement the graph.Ordered interface.
// - T: The type of the graph's vertex data, which can be any type.
func Density[K graph.Ordered, T any](g graph.Interface[K, T]) (float64, error) {
	if g == nil {
		return 0, graph.ErrNilInputGraph
	}

	traits := g.Traits()

	// Get the number of vertices and edges
	vertexCount, err := g.Order()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve vertex count: %w", err)
	}

	edgeCount, err := g.Size()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve edge count: %w", err)
	}

	// If there are fewer than 2 vertices, the density is 0
	if vertexCount < 2 {
		return 0, nil
	}

	// Calculate density based on graph traits
	if traits.IsDirected {
		// Directed graph
		return float64(edgeCount) / float64(vertexCount*(vertexCount-1)), nil
	}

	// Undirected graph
	return float64(2*edgeCount) / float64(vertexCount*(vertexCount-1)), nil
}
