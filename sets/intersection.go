// Copyright (c) 2024 Six After, Inc
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

func Intersection[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	if g == nil || h == nil {
		return nil, fmt.Errorf("one or both input graphs are nil")
	}

	// Ensure the graphs have the same traits
	if !g.Traits().Equals(h.Traits()) {
		return nil, graph.ErrGraphTypeMismatch
	}

	// Determine if the graphs are directed
	if g.Traits().IsDirected {
		return intersectionDirected(g, h)
	}

	return intersectionUndirected(g, h)
}

// intersectionDirected computes the intersection of two directed graphs.
// It ensures that all common vertices are included and only the edges present
// in both graphs with the same direction are added.
// Edge properties from the first graph (g) are preserved in the intersection graph.
func intersectionDirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Step 1: Verify that both graphs are directed.
	if !g.Traits().IsDirected || !h.Traits().IsDirected {
		return nil, graph.ErrGraphTypeMismatch
	}

	// Step 2: Create a new graph for the intersection, inheriting traits from g.
	result, err := simple.NewLike(g)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new graph like g: %w", err)
	}

	// Step 3: Retrieve adjacency maps for both graphs.
	gAdjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve adjacency map of graph g: %w", err)
	}

	hAdjMap, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve adjacency map of graph h: %w", err)
	}

	// Step 4: Add common vertices to the result graph.
	for vertexID := range gAdjMap {
		if _, exists := hAdjMap[vertexID]; exists {
			// Retrieve the vertex from graph g.
			vertexG, err := g.Vertex(vertexID)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve vertex %v from graph g: %w", vertexID, err)
			}

			// Clone the vertex to preserve its properties.
			clonedVertex := vertexG.Clone()

			// Add the cloned vertex to the result graph.
			err = result.AddVertex(clonedVertex)
			if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, fmt.Errorf("failed to add vertex %v to the intersection graph: %w", vertexID, err)
			}
		}
	}

	// Step 5: Add common edges to the result graph.
	for source, gTargets := range gAdjMap {
		for target := range gTargets {
			// Check if the edge exists in both graphs with the same direction.
			if _, exists := hAdjMap[source][target]; !exists {
				continue
			}

			// Retrieve the edge from graph g.
			edgeG, err := g.Edge(source, target)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve edge (%v -> %v) from graph g: %w", source, target, err)
			}

			// Clone the edge using the Clone() method to preserve its properties and direction.
			clonedEdge := simple.NewEdge(source, target, edgeG.Properties().Clone())

			// Add the cloned edge to the result graph.
			err = result.AddEdge(clonedEdge)
			if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, fmt.Errorf("failed to add edge (%v -> %v) to the intersection graph: %w", source, target, err)
			}
		}
	}

	return result, nil
}

// intersectionUndirected computes the intersection of two undirected graphs.
// It ensures that edges are treated symmetrically and properties are preserved.
//
// Parameters:
//   - g: The first input graph.
//   - h: The second input graph.
//
// Returns:
//   - A new graph representing the intersection of g and h.
//   - An error if any operation fails (e.g., incompatible graph types).
func intersectionUndirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Create a new graph like g for the result.
	result, err := simple.NewLike(g)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new graph like g: %w", err)
	}

	// Retrieve adjacency maps for both graphs.
	gAdj, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}

	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Phase 1: Add all common vertices to the result graph.
	for vertex := range gAdj {
		if _, exists := hAdj[vertex]; exists {
			// Retrieve the vertex from graph g.
			v, err := g.Vertex(vertex)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve vertex %v from g: %w", vertex, err)
			}
			clonedV := v.Clone()

			// Add the cloned vertex to the result graph.
			err = result.AddVertex(clonedV)
			if err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, fmt.Errorf("failed to add vertex %v to result graph: %w", vertex, err)
			}
		}
	}

	// Phase 2: Add all common edges to the result graph.
	for source, gTargets := range gAdj {
		for target := range gTargets {
			// Enforce source < target to prevent duplicate undirected edges.
			if !isOrdered(source, target) {
				continue
			}

			// Check if the edge exists in both graphs.
			if _, exists := hAdj[source][target]; !exists {
				continue
			}

			// Retrieve the edge from graph g to preserve its properties.
			edgeG, err := g.Edge(source, target)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve edge (%v -- %v) from g: %w", source, target, err)
			}

			// Clone the edge using the Clone method.
			clonedEdge := simple.NewEdge(source, target, edgeG.Properties().Clone())

			// Add the cloned edge to the result graph.
			err = result.AddEdge(clonedEdge)
			if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, fmt.Errorf("failed to add edge (%v -- %v) to result graph: %w", source, target, err)
			}
		}
	}

	return result, nil
}

// isOrdered enforces an order on the keys to prevent duplicate undirected edges.
// It returns true if source < target.
func isOrdered[K graph.Ordered](source, target K) bool {
	return source < target
}
