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

// PageRank calculates the PageRank for each vertex in the given graph.
// It returns a map where each key corresponds to a vertex in the graph,
// and the value represents its PageRank score.
func PageRank[K graph.Ordered, T any](g graph.Interface[K, T], dampingFactor float64, maxIterations int, tol float64) (map[K]float64, error) {
	// Validate input graph
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	// Validate dampingFactor
	if dampingFactor <= 0.0 || dampingFactor >= 1.0 {
		return nil, errors.New("pagerank: dampingFactor must be between 0 and 1")
	}

	// Validate maxIterations
	if maxIterations <= 0 {
		return nil, errors.New("pagerank: maxIterations must be greater than 0")
	}

	// Validate tolerance
	if tol <= 0.0 {
		return nil, errors.New("pagerank: tolerance must be greater than 0")
	}

	// Retrieve graph traits to determine directionality
	traits := g.Traits()
	if !traits.IsDirected {
		return nil, errors.New("pagerank: PageRank is typically applied to directed graphs")
	}

	// Retrieve all vertices in the graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("pagerank: failed to retrieve vertices: %w", err)
	}

	N := len(vertices)
	if N == 0 {
		return nil, errors.New("pagerank: graph has no vertices")
	}

	// Initialize PageRank scores to 1/N
	pr := make([]float64, N)
	for i := 0; i < N; i++ {
		pr[i] = 1.0 / float64(N)
	}

	// Retrieve the hash function
	hashFunc := g.Hash()

	// Prepare a list of vertex IDs and a mapping from ID to index
	vertexIDs := make([]K, 0, N)
	idToIndex := make(map[K]int, N)
	for idx, vertex := range vertices {
		k := hashFunc(vertex.Value()) // Extract K using the hash function
		vertexIDs = append(vertexIDs, k)
		idToIndex[k] = idx
	}

	// Identify dangling nodes (vertices with no outgoing edges)
	danglingNodes := make([]int, 0)
	for i, id := range vertexIDs {
		outDegree, err := g.OutDegree(id)
		if err != nil {
			return nil, fmt.Errorf("pagerank: failed to get OutDegree for vertex %v: %w", id, err)
		}
		if outDegree == 0 {
			danglingNodes = append(danglingNodes, i)
		}
	}

	// Precompute the teleportation factor
	teleport := (1.0 - dampingFactor) / float64(N)

	// Retrieve PredecessorMap once
	predecessorsMap, err := g.PredecessorMap()
	if err != nil {
		return nil, fmt.Errorf("pagerank: failed to retrieve PredecessorMap: %w", err)
	}

	// Power Iteration Loop
	for iter := 0; iter < maxIterations; iter++ {
		newPr := make([]float64, N)

		// Calculate the total PageRank from dangling nodes
		danglingSum := 0.0
		for _, i := range danglingNodes {
			danglingSum += pr[i]
		}
		danglingContribution := dampingFactor * danglingSum / float64(N)

		// Initialize all PageRank scores with teleportation and dangling contribution
		for i := 0; i < N; i++ {
			newPr[i] = teleport + danglingContribution
		}

		// Compute contributions from incoming links
		for i, id := range vertexIDs {
			inNeighbors, exists := predecessorsMap[id]
			if !exists {
				continue
			}
			for src := range inNeighbors {
				srcIdx, exists := idToIndex[src]
				if !exists {
					continue
				}
				outDegree, err := g.OutDegree(src)
				if err != nil || outDegree == 0 {
					continue
				}
				newPr[i] += dampingFactor * pr[srcIdx] / float64(outDegree)
			}
		}

		// Compute the total difference for convergence
		diff := 0.0
		for i := 0; i < N; i++ {
			diff += math.Abs(newPr[i] - pr[i])
		}

		// Update PageRank scores
		pr = newPr

		// Check for convergence
		if diff < tol {
			break
		}
	}

	// Assign PageRank scores to the result map
	result := make(map[K]float64, N)
	for i, id := range vertexIDs {
		result[id] = pr[i]
	}

	return result, nil
}
