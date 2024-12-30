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
