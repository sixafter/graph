// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"errors"
	"fmt"
	"math"

	"github.com/sixafter/graph"
)

// PageRank computes the PageRank for each vertex in a directed graph.
//
// PageRank is a graph algorithm that assigns a score to each vertex in a graph,
// representing its relative importance based on its connections. It considers the
// direction and optionally the weights of edges to calculate the scores.
//
// This implementation supports:
// - Directed graphs (required).
// - Weighted and unweighted edges.
//
// The algorithm uses an iterative process to compute the PageRank scores until
// convergence or a maximum number of iterations is reached.
//
// Parameters:
//   - g: The graph.Interface representing the directed graph. The graph must have
//     directed traits, and may optionally support weighted edges.
//   - dampingFactor: A float64 in the range (0, 1) representing the probability of
//     continuing to follow links in a random walk (default is typically 0.85).
//   - maxIterations: An integer representing the maximum number of iterations to
//     perform before stopping.
//   - tol: A float64 representing the tolerance for convergence. The algorithm stops
//     iterating when the total difference in PageRank values between iterations is
//     less than this value.
//
// Returns:
//   - A map[K]float64 where each key is a vertex ID and each value is its PageRank score.
//     The scores are normalized to sum to 1.0.
//   - An error if the computation fails (e.g., invalid graph or parameters).
//
// Example Usage:
//   g, _ := simple.New[int, int](graph.IntHash, graph.Directed())
//   // Add vertices and edges to `g`
//   pr, err := PageRank(g, 0.85, 100, 1e-6)
//   if err != nil {
//       fmt.Println("Error:", err)
//   } else {
//       fmt.Println("PageRank Scores:", pr)
//   }
func PageRank[K graph.Ordered, T any](
	g graph.Interface[K, T],
	dampingFactor float64,
	maxIterations int,
	tol float64,
) (map[K]float64, error) {
	// Validate input parameters
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	if !g.Traits().IsDirected {
		return nil, graph.ErrUndirectedGraph
	}

	if dampingFactor <= 0.0 || dampingFactor >= 1.0 {
		return nil, errors.New("pagerank: dampingFactor must be between 0 and 1")
	}
	if maxIterations <= 0 {
		return nil, errors.New("pagerank: maxIterations must be greater than 0")
	}
	if tol <= 0.0 {
		return nil, errors.New("pagerank: tolerance must be greater than 0")
	}

	// Verify that the graph is directed
	traits := g.Traits()
	if !traits.IsDirected {
		return nil, errors.New("pagerank: graph must be directed")
	}

	// Retrieve vertices and initialize PageRank scores
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("pagerank: failed to retrieve vertices: %w", err)
	}

	N := len(vertices)
	if N == 0 {
		return nil, errors.New("pagerank: graph has no vertices")
	}

	pr := make([]float64, N)        // Current PageRank values
	prNew := make([]float64, N)     // Temporary storage for updates
	vertexIDs := make([]K, N)       // Vertex IDs
	idToIndex := make(map[K]int, N) // Map of vertex ID to index

	// Initialize PageRank values uniformly
	for i, vertex := range vertices {
		id := g.Hash()(vertex.Value())
		vertexIDs[i] = id
		idToIndex[id] = i
		pr[i] = 1.0 / float64(N)
	}

	// Precompute teleportation factor
	teleport := (1.0 - dampingFactor) / float64(N)

	// Retrieve the predecessor map for directed graphs
	predecessors, err := g.PredecessorMap()
	if err != nil {
		return nil, fmt.Errorf("pagerank: failed to retrieve predecessor map: %w", err)
	}

	// Main iteration loop
	for iter := 0; iter < maxIterations; iter++ {
		// Reset new PageRank values to the teleportation contribution
		for i := range prNew {
			prNew[i] = teleport
		}

		// Calculate contributions from predecessors
		for i, id := range vertexIDs {
			neighbors, exists := predecessors[id]
			if !exists {
				continue
			}

			for neighborID, edge := range neighbors {
				neighborIdx, exists := idToIndex[neighborID]
				if !exists {
					continue
				}

				// Calculate the weight of the edge
				weight := 1.0
				if traits.IsWeighted {
					weight = edge.Properties().Weight()
				}

				// Get the out-degree of the neighbor
				outDegree, err := g.OutDegree(neighborID)
				if err != nil || outDegree == 0 {
					continue
				}

				// Distribute the neighbor's contribution
				contribution := (pr[neighborIdx] * weight) / float64(outDegree)
				prNew[i] += dampingFactor * contribution
			}
		}

		// Check for convergence
		diff := 0.0
		for i := range pr {
			diff += math.Abs(prNew[i] - pr[i])
		}

		copy(pr, prNew) // Update current PageRank values

		if diff < tol {
			break
		}
	}

	// Normalize PageRank values to sum to 1
	totalPR := 0.0
	for _, score := range pr {
		totalPR += score
	}
	result := make(map[K]float64, N)
	for i, id := range vertexIDs {
		result[id] = pr[i] / totalPR
	}

	return result, nil
}
