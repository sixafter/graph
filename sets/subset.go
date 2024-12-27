// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"errors"
	"fmt"

	"github.com/sixafter/graph"
)

func IsSubset[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Ensure the graphs have the same traits
	if nil == g || nil == h {
		return false, nil
	}

	if !g.Traits().Equals(h.Traits()) {
		return false, graph.ErrGraphTypeMismatch
	}

	// Determine if the graphs are directed
	if g.Traits().IsDirected {
		return isSubsetDirected(g, h)
	}
	return isSubsetUndirected(g, h)
}

// isSubsetDirected checks if h is a subset of g for directed graphs.
func isSubsetDirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Get adjacency maps for both graphs
	gAdjMap, err := g.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}
	hAdjMap, err := h.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Get predecessor maps for both graphs
	gPredMap, err := g.PredecessorMap()
	if err != nil {
		return false, fmt.Errorf("failed to get predecessor map of g: %w", err)
	}
	_, err = h.PredecessorMap()
	if err != nil {
		return false, fmt.Errorf("failed to get predecessor map of h: %w", err)
	}

	// Collect all vertices in g
	gVertices := make(map[K]struct{})
	for v := range gAdjMap {
		gVertices[v] = struct{}{}
	}
	for v := range gPredMap {
		gVertices[v] = struct{}{}
	}

	// Check that all vertices in g are present in h
	for v := range gVertices {
		_, err := h.Vertex(v)
		if err != nil {
			if errors.Is(err, graph.ErrVertexNotFound) {
				return false, nil
			}
			return false, fmt.Errorf("failed to get vertex %v from h: %w", v, err)
		}
	}

	// Check that all edges in g are present in h with matching properties
	for sourceHash, gTargets := range gAdjMap {
		for targetHash, gEdge := range gTargets {
			hTargets, exists := hAdjMap[sourceHash]
			if !exists {
				return false, nil
			}
			hEdge, edgeExists := hTargets[targetHash]
			if !edgeExists {
				return false, nil
			}
			// Use edgePropertiesEqual to compare edge properties
			if !edgePropertiesEqual(gEdge.Properties(), hEdge.Properties()) {
				return false, nil
			}
		}
	}

	return true, nil
}

// isSubsetUndirected checks if h is a subset of g for undirected graphs.
func isSubsetUndirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Get adjacency maps for both graphs
	gAdjMap, err := g.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}
	hAdjMap, err := h.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Collect all vertices in g
	gVertices := make(map[K]struct{})
	for v := range gAdjMap {
		gVertices[v] = struct{}{}
	}
	// Include vertices that are targets but not sources
	for _, targets := range gAdjMap {
		for targetHash := range targets {
			gVertices[targetHash] = struct{}{}
		}
	}

	// Check that all vertices in g are present in h
	for v := range gVertices {
		_, err := h.Vertex(v)
		if err != nil {
			if errors.Is(err, graph.ErrVertexNotFound) {
				return false, nil
			}
			return false, fmt.Errorf("failed to get vertex %v from h: %w", v, err)
		}
	}

	// Check that all edges in g are present in h with matching properties
	for sourceHash, gNeighbors := range gAdjMap {
		for targetHash, gEdge := range gNeighbors {
			// For undirected graphs, check both directions
			edgeFound := false
			if hNeighbors, exists := hAdjMap[sourceHash]; exists {
				if hEdge, edgeExists := hNeighbors[targetHash]; edgeExists {
					if edgePropertiesEqual(gEdge.Properties(), hEdge.Properties()) {
						edgeFound = true
					}
				}
			}
			if !edgeFound {
				if hNeighbors, exists := hAdjMap[targetHash]; exists {
					if hEdge, edgeExists := hNeighbors[sourceHash]; edgeExists {
						if edgePropertiesEqual(gEdge.Properties(), hEdge.Properties()) {
							edgeFound = true
						}
					}
				}
			}
			if !edgeFound {
				return false, nil
			}
		}
	}

	return true, nil
}

func edgePropertiesEqual(a, b graph.EdgeProperties) bool {
	// Only compare properties that are significant
	return a.Weight() == b.Weight() // Or ignore weight if not significant
}
