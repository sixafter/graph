// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"fmt"

	"github.com/sixafter/graph"
)

func IsDisjoint[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Check if graph traits match
	if nil == g || nil == h {
		return false, nil
	}

	if !g.Traits().Equals(h.Traits()) {
		return false, graph.ErrGraphTypeMismatch
	}

	// Determine if the graphs are directed
	if g.Traits().IsDirected {
		return isDisjointDirected(g, h)
	}
	return isDisjointUndirected(g, h)
}

// isDisjointDirected checks if two directed graphs are disjoint.
func isDisjointDirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Get adjacency maps for both graphs
	gAdjMap, err := g.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}

	hAdjMap, err := h.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Check for common vertices
	for v := range gAdjMap {
		if _, exists := hAdjMap[v]; exists {
			return false, nil
		}
	}

	// Check for common edges
	for gSource, gTargets := range gAdjMap {
		if hTargets, exists := hAdjMap[gSource]; exists {
			for gTarget := range gTargets {
				if _, edgeExists := hTargets[gTarget]; edgeExists {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

// isDisjointUndirected checks if two undirected graphs are disjoint.
func isDisjointUndirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Get adjacency maps for both graphs
	gAdjMap, err := g.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}

	hAdjMap, err := h.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Use a set to avoid processing undirected edges twice
	processedEdges := make(map[string]struct{})

	// Check for common vertices
	for v := range gAdjMap {
		if _, exists := hAdjMap[v]; exists {
			return false, nil
		}
	}

	// Check for common edges
	for gSource, gTargets := range gAdjMap {
		for gTarget := range gTargets {
			// Create a unique key for the undirected edge
			edgeKey := undirectedEdgeKey(gSource, gTarget)
			if _, processed := processedEdges[edgeKey]; processed {
				continue // Edge already processed
			}

			// Check if the same edge exists in h
			if hTargets, exists := hAdjMap[gSource]; exists {
				if _, edgeExists := hTargets[gTarget]; edgeExists {
					return false, nil
				}
			}
			if hTargets, exists := hAdjMap[gTarget]; exists {
				if _, edgeExists := hTargets[gSource]; edgeExists {
					return false, nil
				}
			}

			// Mark the edge as processed
			processedEdges[edgeKey] = struct{}{}
		}
	}

	return true, nil
}
