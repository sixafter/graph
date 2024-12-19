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

// Difference returns a new graph containing the vertices and edges that are present in graph g
// but not in graph h. Both graphs must have the same traits (e.g., directedness).
// The resulting graph inherits traits from graph g.
//
// Parameters:
//   - g: The first input graph.
//   - h: The second input graph.
//
// Returns:
//   - A new graph representing the difference of g and h (g - h).
//   - An error if the graph traits do not match or if any operation fails.
//
// Usage Example:
//
//	difference, err := Difference(g, h)
//	if err != nil {
//	    log.Fatal(err)
//	}
func Difference[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Ensure the graphs have the same traits
	if nil == g || nil == h {
		return nil, nil
	}

	if !g.Traits().Equals(h.Traits()) {
		return nil, graph.ErrGraphTypeMismatch
	}

	// Determine if the graphs are directed
	if g.Traits().IsDirected {
		return differenceDirected(g, h)
	}

	return differenceUndirected(g, h)
}

// differenceDirected computes the difference of two directed graphs (g - h).
func differenceDirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Clone graph g to start with the difference
	diff, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone g: %w", err)
	}

	// Get adjacency map of h
	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Remove edges present in h from the diff graph
	for hSourceHash, hTargets := range hAdj {
		for hTargetHash := range hTargets {
			err := diff.RemoveEdge(hSourceHash, hTargetHash)
			if err != nil && !errors.Is(err, graph.ErrEdgeNotFound) {
				return nil, fmt.Errorf("failed to remove edge (%v, %v): %w", hSourceHash, hTargetHash, err)
			}
		}
	}

	return diff, nil
}

// differenceUndirected computes the difference of two undirected graphs (g - h).
func differenceUndirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Clone graph g to start with the difference
	diff, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone g: %w", err)
	}

	// Get adjacency map of h
	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Remove edges present in h from the diff graph
	for hSourceHash, hNeighbors := range hAdj {
		for hTargetHash := range hNeighbors {
			// Remove the edge from diff
			err = diff.RemoveEdge(hSourceHash, hTargetHash)
			if err != nil && !errors.Is(err, graph.ErrEdgeNotFound) {
				return nil, fmt.Errorf("failed to remove edge (%v, %v): %w", hSourceHash, hTargetHash, err)
			}

			// No need to remove the reverse edge separately since RemoveEdge in undirected graphs should handle it
		}
	}

	// Remove isolated vertices
	diffAdj, err := diff.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of diff: %w", err)
	}
	for vertexHash := range diffAdj {
		// Check if the vertex has any neighbors
		if len(diffAdj[vertexHash]) == 0 {
			err = diff.RemoveVertex(vertexHash)
			if err != nil && !errors.Is(err, graph.ErrVertexNotFound) {
				return nil, fmt.Errorf("failed to remove isolated vertex %v: %w", vertexHash, err)
			}
		}
	}

	return diff, nil
}

// SymmetricDifference returns a new graph that represents the symmetric difference between graphs g and h.
// The symmetric difference includes vertices and edges that are in either graph, but not in both.
func SymmetricDifference[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Ensure the graphs are not nil
	if g == nil || h == nil {
		return nil, fmt.Errorf("one or both input graphs are nil")
	}

	// Ensure the graphs have the same traits
	if !g.Traits().Equals(h.Traits()) {
		return nil, graph.ErrGraphTypeMismatch
	}

	// Determine if the graphs are directed
	if g.Traits().IsDirected {
		return symmetricDifferenceDirected(g, h)
	}
	return symmetricDifferenceUndirected(g, h)
}

// symmetricDifferenceDirected computes the symmetric difference of two directed graphs.
func symmetricDifferenceDirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Clone graph g to start with the result
	result, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone g: %w", err)
	}

	// Get adjacency maps for both graphs
	gAdj, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}

	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Process edges from g
	for gSource, gTargets := range gAdj {
		for gTarget := range gTargets {
			if _, exists := hAdj[gSource][gTarget]; exists {
				// Edge exists in both graphs, remove it
				err = result.RemoveEdge(gSource, gTarget)
				if err != nil && !errors.Is(err, graph.ErrEdgeNotFound) {
					return nil, fmt.Errorf("failed to remove edge (%v, %v): %w", gSource, gTarget, err)
				}
			}
		}
	}

	// Process edges from h
	for hSource, hTargets := range hAdj {
		// Ensure the source vertex exists in the result graph
		if exists, _ := result.HasVertex(hSource); !exists {
			var hVertex graph.Vertex[K, T]
			hVertex, err = h.Vertex(hSource)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch vertex %v from h: %w", hSource, err)
			}
			if err = result.AddVertex(hVertex); err != nil {
				return nil, fmt.Errorf("failed to add vertex %v to result: %w", hSource, err)
			}
		}
		for hTarget := range hTargets {
			// Ensure the target vertex exists in the result graph
			if exists, _ := result.HasVertex(hTarget); !exists {
				var hVertex graph.Vertex[K, T]
				hVertex, err = h.Vertex(hTarget)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch vertex %v from h: %w", hTarget, err)
				}
				if err := result.AddVertex(hVertex); err != nil {
					return nil, fmt.Errorf("failed to add vertex %v to result: %w", hTarget, err)
				}
			}

			if _, exists := gAdj[hSource][hTarget]; !exists {
				// Edge exists only in h, add it
				err = result.AddEdgeWithOptions(hSource, hTarget)
				if err != nil {
					return nil, fmt.Errorf("failed to add edge (%v, %v): %w", hSource, hTarget, err)
				}
			}
		}
	}

	return result, nil
}

// symmetricDifferenceUndirected computes the symmetric difference of two undirected graphs.
func symmetricDifferenceUndirected[K graph.Ordered, T any](g, h graph.Interface[K, T]) (graph.Interface[K, T], error) {
	// Clone graph g to start with the result
	result, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone g: %w", err)
	}

	// Get adjacency maps for both graphs
	gAdj, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of g: %w", err)
	}

	hAdj, err := h.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map of h: %w", err)
	}

	// Process edges from g
	for gSource, gTargets := range gAdj {
		for gTarget := range gTargets {
			if _, exists := hAdj[gSource][gTarget]; exists {
				// Edge exists in both graphs, remove it
				err = result.RemoveEdge(gSource, gTarget)
				if err != nil && !errors.Is(err, graph.ErrEdgeNotFound) {
					return nil, fmt.Errorf("failed to remove edge (%v, %v): %w", gSource, gTarget, err)
				}
			}
		}
	}

	// Process edges from h
	for hSource, hTargets := range hAdj {
		// Ensure the source vertex exists in the result graph
		if exists, _ := result.HasVertex(hSource); !exists {
			var hVertex graph.Vertex[K, T]
			hVertex, err = h.Vertex(hSource)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch vertex %v from h: %w", hSource, err)
			}
			if err := result.AddVertex(hVertex); err != nil {
				return nil, fmt.Errorf("failed to add vertex %v to result: %w", hSource, err)
			}
		}
		for hTarget := range hTargets {
			// Ensure the target vertex exists in the result graph
			if exists, _ := result.HasVertex(hTarget); !exists {
				var hVertex graph.Vertex[K, T]
				hVertex, err = h.Vertex(hTarget)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch vertex %v from h: %w", hTarget, err)
				}
				if err := result.AddVertex(hVertex); err != nil {
					return nil, fmt.Errorf("failed to add vertex %v to result: %w", hTarget, err)
				}
			}

			// Check if the edge exists in the result to avoid duplicates
			if exists, _ := result.HasEdge(hSource, hTarget); !exists {
				// Edge exists only in h, add it
				err = result.AddEdgeWithOptions(hSource, hTarget)
				if err != nil {
					return nil, fmt.Errorf("failed to add edge (%v, %v): %w", hSource, hTarget, err)
				}
			}
		}
	}

	return result, nil
}
