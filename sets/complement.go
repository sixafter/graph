// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"errors"
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
)

func Complement[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	// Determine if the graph is directed
	if g.Traits().IsDirected {
		return complementDirected(g)
	}

	return complementUndirected(g)
}

// complementDirected computes the complement of a directed graph.
// It returns a new graph containing all the vertices of the original graph
// and adds directed edges that do not exist in the original graph.
func complementDirected[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	complement, err := simple.NewLike(g)
	if err != nil {
		return nil, fmt.Errorf("failed to clone the graph: %w", err)
	}

	// Add all vertices from the original graph to the complement graph
	// This ensures that even vertices without outgoing or incoming edges are included
	err = complement.AddVerticesFrom(g)
	if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
		return nil, fmt.Errorf("failed to add vertices from original graph: %w", err)
	}

	// Retrieve all vertices from the original graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vertices: %w", err)
	}

	// Initialize a slice to store all unique keys (K) from the vertices
	var keys []K

	// Populate keys using the Hash function on vertex values
	for _, vertex := range vertices {
		key := g.Hash()(vertex.Value())
		keys = append(keys, key)
	}

	// Iterate over all ordered pairs of distinct vertices
	for _, sourceKey := range keys {
		for _, targetKey := range keys {
			// Skip self-loops if they are not allowed
			if sourceKey == targetKey {
				continue
			}

			// Check if the edge (sourceKey -> targetKey) exists in the original graph
			var hasEdge bool
			hasEdge, err = g.HasEdge(sourceKey, targetKey)
			if err != nil {
				return nil, fmt.Errorf("failed to check existence of edge (%v, %v): %w", sourceKey, targetKey, err)
			}

			// If the edge does not exist, add it to the complement graph
			if !hasEdge {
				// Add the edge without any additional properties
				// If your graph implementation requires properties, use AddEdgeWithOptions accordingly
				err = complement.AddEdgeWithOptions(sourceKey, targetKey)
				if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, fmt.Errorf("failed to add edge from %v to %v: %w", sourceKey, targetKey, err)
				}
			}
		}
	}

	return complement, nil
}

// complementUndirected computes the complement of an undirected graph.
// It returns a new graph containing all the vertices of the original graph
// and adds edges that do not exist in the original graph.
func complementUndirected[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	complement, err := simple.NewLike(g)
	if err != nil {
		return nil, fmt.Errorf("failed to clone the graph: %w", err)
	}

	// Add all vertices from the original graph to the complement graph
	// This ensures that even vertices without any edges are included
	err = complement.AddVerticesFrom(g)
	if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
		return nil, fmt.Errorf("failed to add vertices from original graph: %w", err)
	}

	// Retrieve the adjacency map of the original graph for efficient lookups
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	// Retrieve all vertices from the original graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vertices: %w", err)
	}

	// Initialize a set to store all unique keys (K) from the vertices
	keysSet := make(map[K]struct{})

	// Populate keysSet using the Hash function on vertex values
	for _, vertex := range vertices {
		key := g.Hash()(vertex.Value())
		keysSet[key] = struct{}{}
	}

	// Additionally, ensure all target keys from the adjacency map are included
	for sourceKey, targets := range adjMap {
		keysSet[sourceKey] = struct{}{}
		for targetKey := range targets {
			keysSet[targetKey] = struct{}{}
		}
	}

	// Convert keysSet to a slice for iteration
	keys := make([]K, 0, len(keysSet))
	for key := range keysSet {
		keys = append(keys, key)
	}

	// Iterate over all unique unordered pairs of distinct vertices
	// Since the graph is undirected, (A,B) is equivalent to (B,A), so we ensure each pair is considered only once
	for i, sourceKey := range keys {
		for j := i + 1; j < len(keys); j++ {
			targetKey := keys[j]

			// Check if the edge (sourceKey <-> targetKey) exists in the original graph using HasEdge
			var hasEdge bool
			hasEdge, err = g.HasEdge(sourceKey, targetKey)
			if err != nil {
				return nil, fmt.Errorf("failed to check existence of edge (%v, %v): %w", sourceKey, targetKey, err)
			}

			// If the edge does not exist, add it to the complement graph
			if !hasEdge {
				// Add the edge without any additional properties
				// If your graph implementation requires properties, use AddEdgeWithOptions accordingly
				err = complement.AddEdgeWithOptions(sourceKey, targetKey)
				if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, fmt.Errorf("failed to add edge between %v and %v: %w", sourceKey, targetKey, err)
				}
			}
		}
	}

	return complement, nil
}
