// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/sixafter/graph"
)

// EigenvectorCentrality calculates the eigenvector centrality for each vertex in the given graph.
// It returns a map where each key corresponds to a vertex in the graph,
// and the value represents its eigenvector centrality.
//
// Parameters:
//   - g: The graph for which to compute eigenvector centrality. It must implement the graph.Interface[K, T] interface.
//
// Returns:
//   - A map[K]float64 where each key corresponds to a vertex in the graph,
//     and the value represents its eigenvector centrality.
//   - An error if the input graph is nil, or if there is a failure in retrieving vertices or adjacency maps.
//
// Errors:
//   - graph.ErrNilInputGraph: Returned if the input graph g is nil.
//   - Errors propagated from graph methods.
//
// Usage:
//
//	centrality, err := metrics.EigenvectorCentrality(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute eigenvector centrality: %v", err)
//	}
//	for vertex, score := range centrality {
//	    fmt.Printf("Vertex %v has an eigenvector centrality of %.4f\n", vertex, score)
//	}
//
// Complexity: O(iterations * (V + E)), where iterations is the number of power iterations,
// V is the number of vertices, and E is the number of edges.
func EigenvectorCentrality[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Validate input graph
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	// Retrieve graph traits to determine directionality
	traits := g.Traits()
	isDirected := traits.IsDirected

	// Retrieve all vertices in the graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("eigenvector_centrality: failed to retrieve vertices: %w", err)
	}

	// Initialize the eigenvector centrality map
	centrality := make(map[K]float64, len(vertices))

	// Retrieve adjacency map
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("eigenvector_centrality: failed to retrieve adjacency map: %w", err)
	}

	// Prepare a list of vertex IDs for easy indexing
	vertexIDs := make([]K, 0, len(vertices))
	idToIndex := make(map[K]int, len(vertices))
	for idx, vertex := range vertices {
		vertexIDs = append(vertexIDs, vertex.ID())
		idToIndex[vertex.ID()] = idx
	}

	V := len(vertexIDs)
	if V == 0 {
		return centrality, nil // Empty graph
	}

	// Initialize centrality scores to 1.0
	x := make([]float64, V)
	for i := range x {
		x[i] = 1.0
	}

	// Power Iteration parameters
	maxIterations := 100
	tolerance := 1e-6

	var incomingAdjacencyMap map[K]map[K]struct{}
	if isDirected {
		// Preprocess incoming adjacency map
		incomingAdjacencyMap = make(map[K]map[K]struct{}, V)
		for _, vertex := range vertexIDs {
			incomingAdjacencyMap[vertex] = make(map[K]struct{})
		}

		for _, vertex := range vertexIDs {
			for neighbor := range adjacencyMap[vertex] {
				incomingAdjacencyMap[neighbor][vertex] = struct{}{}
			}
		}
	}

	// Power Iteration Loop
	for iter := 0; iter < maxIterations; iter++ {
		y := make([]float64, V)

		// Parallel computation of y = A * x
		var wg sync.WaitGroup
		numWorkers := 8 // Adjust based on your system's capabilities
		workChan := make(chan int, V)

		// Worker function
		worker := func() {
			defer wg.Done()
			for i := range workChan {
				vertex := vertexIDs[i]
				if isDirected {
					// Sum incoming neighbors' centrality scores
					for inNeighbor := range incomingAdjacencyMap[vertex] {
						neighborIdx := idToIndex[inNeighbor]
						y[i] += x[neighborIdx]
					}
				} else {
					// Sum all neighbors' centrality scores
					for neighbor := range adjacencyMap[vertex] {
						neighborIdx := idToIndex[neighbor]
						y[i] += x[neighborIdx]
					}
				}
			}
		}

		// Start workers
		wg.Add(numWorkers)
		for w := 0; w < numWorkers; w++ {
			go worker()
		}

		// Dispatch work
		for i := 0; i < V; i++ {
			workChan <- i
		}
		close(workChan)

		// Wait for all workers to finish
		wg.Wait()

		// Normalize y (Euclidean norm)
		var norm float64
		for _, val := range y {
			norm += val * val
		}
		norm = math.Sqrt(norm)
		if norm == 0 {
			return nil, errors.New("eigenvector_centrality: zero norm encountered during normalization")
		}
		for i := range y {
			y[i] /= norm
		}

		// Check for convergence
		var diff float64
		for i := 0; i < V; i++ {
			diff += math.Abs(y[i] - x[i])
		}
		if diff < tolerance {
			// Converged
			break
		}

		// Update x for next iteration
		x = y
	}

	// Assign centrality scores
	for i, vertex := range vertexIDs {
		centrality[vertex] = x[i]
	}

	return centrality, nil
}
