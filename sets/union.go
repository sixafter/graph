// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"errors"
	"fmt"

	"github.com/sixafter/graph"
)

func Union[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	if nil == g || nil == h {
		return nil, graph.ErrNilInputGraph
	}

	// Ensure the graphs have the same traits
	if !g.Traits().Equals(h.Traits()) {
		return nil, graph.ErrGraphTypeMismatch
	}

	// Determine if the graphs are directed
	if g.Traits().IsDirected {
		return unionDirected(g, h)
	}

	return unionUndirected(g, h)
}

// unionDirected performs the union operation for directed graphs.
func unionDirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Clone graph g to start with the union
	u, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone g: %w", err)
	}

	// Get adjacency map of h
	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Add vertices and edges from h to the union graph
	for sourceHash, targets := range hAdj {
		// Retrieve the source vertex value from h
		var sourceValue graph.Vertex[K, T]
		sourceValue, err = h.Vertex(sourceHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get vertex %v: %w", sourceHash, err)
		}

		// Add source vertex to the union graph
		err = u.AddVertex(sourceValue)
		if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, fmt.Errorf("failed to add vertex %v: %w", sourceHash, err)
		}

		for targetHash, edge := range targets {
			// Retrieve the target vertex value from h
			var targetValue graph.Vertex[K, T]
			targetValue, err = h.Vertex(targetHash)
			if err != nil {
				return nil, fmt.Errorf("failed to get vertex %v: %w", targetHash, err)
			}

			// Add target vertex to the union graph
			err = u.AddVertex(targetValue)
			if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, fmt.Errorf("failed to add vertex %v: %w", targetHash, err)
			}

			// Add edge to the union graph using the arguments from copyEdge
			err = u.AddEdge(edge.Clone())
			if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, fmt.Errorf("failed to add edge (%v, %v): %w", edge.Source(), edge.Target(), err)
			}
		}
	}

	return u, nil
}

// unionUndirected performs the union operation for undirected graphs.
func unionUndirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Clone graph g to start with the union
	u, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone g: %w", err)
	}

	// Get adjacency map of h
	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// To avoid processing each undirected edge twice
	processedEdges := make(map[string]struct{})

	// Add vertices and edges from h to the union graph
	for sourceHash, targets := range hAdj {
		// Retrieve the source vertex value from h
		var sourceValue graph.Vertex[K, T]
		sourceValue, err = h.Vertex(sourceHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get vertex %v: %w", sourceHash, err)
		}

		// Add source vertex to the union graph
		err = u.AddVertex(sourceValue)
		if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, fmt.Errorf("failed to add vertex %v: %w", sourceHash, err)
		}

		for targetHash, edge := range targets {
			// Create a unique key for the undirected edge
			edgeKey := undirectedEdgeKey(sourceHash, targetHash)
			if _, exists := processedEdges[edgeKey]; exists {
				continue // Edge already processed
			}
			processedEdges[edgeKey] = struct{}{}

			// Retrieve the target vertex value from h
			var targetValue graph.Vertex[K, T]
			targetValue, err = h.Vertex(targetHash)
			if err != nil {
				return nil, fmt.Errorf("failed to get vertex %v: %w", targetHash, err)
			}

			// Add target vertex to the union graph
			err = u.AddVertex(targetValue)
			if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, fmt.Errorf("failed to add vertex %v: %w", targetHash, err)
			}

			// Add edge to the union graph using the arguments from copyEdge
			err = u.AddEdge(edge.Clone())
			if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, fmt.Errorf("failed to add edge (%v, %v): %w", edge.Source(), edge.Target(), err)
			}
		}
	}

	return u, nil
}

func undirectedEdgeKey[K graph.Ordered](a, b K) string {
	aStr, bStr := fmt.Sprint(a), fmt.Sprint(b)
	if aStr < bStr {
		return aStr + "-" + bStr
	}
	return bStr + "-" + aStr
}
