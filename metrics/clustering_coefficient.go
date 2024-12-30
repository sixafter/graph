// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// ClusteringCoefficient calculates the clustering coefficient for each vertex in the provided graph g.
//
// The clustering coefficient measures the degree to which vertices in a graph tend to cluster together.
// For a given vertex, it quantifies how close its neighbors are to forming a complete subgraph (clique).
//
// The clustering coefficient for a vertex v is defined as:
//
//	C(v) = (number of edges between neighbors of v) / (degree(v) * (degree(v) - 1) / 2)
//
// For undirected graphs:
//   - degree(v) is the number of neighbors.
//   - C(v) is normalized between 0 and 1.
//
// For directed graphs:
//   - degree(v) is the out-degree (number of outgoing edges).
//   - Neighbors are considered as out-neighbors only.
//   - C(v) is calculated without a normalization factor of 2.
//
// If a vertex has fewer than two neighbors, its clustering coefficient is defined as 0,
// since no meaningful clustering can occur.
//
// The function returns a map where each key corresponds to a vertex in the graph,
// and the value represents its clustering coefficient.
//
// Example:
//
//	centrality, err := ClusteringCoefficient(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute clustering coefficient: %v", err)
//	}
//	for vertex, score := range centrality {
//	    fmt.Printf("Vertex %v has a clustering coefficient of %.4f\n", vertex, score)
//	}
//
// Parameters:
//   - g: The graph for which to compute clustering coefficients. It must implement the graph.Interface[K, T] interface.
//
// Returns:
//   - A map[K]float64 where each key corresponds to a vertex in the graph,
//     and the value represents its clustering coefficient.
//   - An error if the input graph is nil, or if there is a failure in retrieving vertices or adjacency maps.
//
// Errors:
//   - graph.ErrNilInputGraph: Returned if the input graph g is nil.
//   - Errors propagated from g.Vertices(), g.AdjacencyMap(), or g.PredecessorMap(): Returned if there is a failure in retrieving vertices or adjacency maps.
//
// Usage:
//
//	g, err := simple.New(graph.IntHash, graph.Undirected())
//	if err != nil {
//	    log.Fatalf("Failed to create graph: %v", err)
//	}
//
//	// Add vertices
//	for i := 1; i <= 5; i++ {
//	    if err := g.AddVertexWithOptions(i); err != nil {
//	        log.Fatalf("Failed to add vertex %d: %v", i, err)
//	    }
//	}
//
//	// Add edges
//	edges := [][2]int{
//	    {1, 2},
//	    {1, 3},
//	    {1, 4},
//	    {2, 3},
//	    {3, 4},
//	    {4, 5},
//	}
//	for _, edge := range edges {
//	    if err := g.AddEdgeWithOptions(edge[0], edge[1]); err != nil {
//	        log.Fatalf("Failed to add edge %v->%v: %v", edge[0], edge[1], err)
//	    }
//	}
//
//	// Compute clustering coefficients
//	clustering, err := ClusteringCoefficient(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute clustering coefficients: %v", err)
//	}
//
//	// Display clustering coefficients
//	for vertex, score := range clustering {
//	    fmt.Printf("Vertex %v has a clustering coefficient of %.4f\n", vertex, score)
//	}
func ClusteringCoefficient[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Validate input graph
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	// Retrieve graph traits to determine directionality
	if g.Traits().IsDirected {
		return clusteringCoefficientDirected(g)
	}
	return clusteringCoefficientUndirected(g)
}

// clusteringCoefficientUndirected computes clustering coefficients for undirected graphs.
func clusteringCoefficientUndirected[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Retrieve all vertices in the graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("clustering_coefficient_undirected: failed to retrieve vertices: %w", err)
	}

	// Retrieve the adjacency map for quick neighbor lookups
	adjacency, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("clustering_coefficient_undirected: failed to retrieve adjacency map: %w", err)
	}

	// Initialize the clustering coefficient map
	clustering := make(map[K]float64, len(vertices))

	hashFunc := g.Hash()

	// Iterate over each vertex to compute its clustering coefficient
	for _, vertex := range vertices {
		vKey := hashFunc(vertex.Value())

		neighborsMap, exists := adjacency[vKey]
		if !exists {
			// If the vertex has no neighbors, its clustering coefficient is 0
			clustering[vKey] = 0.0
			continue
		}

		degree := len(neighborsMap)
		if degree < 2 {
			// Clustering coefficient is 0 for vertices with fewer than two neighbors
			clustering[vKey] = 0.0
			continue
		}

		// Collect all neighbor keys into a slice for easy iteration
		neighborKeys := make([]K, 0, degree)
		for neighborKey := range neighborsMap {
			neighborKeys = append(neighborKeys, neighborKey)
		}

		// Count the number of edges between neighbors
		edgeCount := 0
		for i := 0; i < len(neighborKeys); i++ {
			for j := i + 1; j < len(neighborKeys); j++ {
				u := neighborKeys[i]
				w := neighborKeys[j]

				if _, connected := adjacency[u][w]; connected {
					edgeCount++
				}
			}
		}

		// Calculate the clustering coefficient
		// C(v) = 2 * edgeCount / (degree * (degree - 1))
		clustering[vKey] = (2.0 * float64(edgeCount)) / (float64(degree) * float64(degree-1))
	}

	return clustering, nil
}

// clusteringCoefficientDirected computes clustering coefficients for directed graphs.
func clusteringCoefficientDirected[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Retrieve all vertices in the graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("clustering_coefficient_directed: failed to retrieve vertices: %w", err)
	}

	// Retrieve the adjacency map (outgoing edges)
	adjacency, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("clustering_coefficient_directed: failed to retrieve adjacency map: %w", err)
	}

	// Initialize the clustering coefficient map
	clustering := make(map[K]float64, len(vertices))

	hashFunc := g.Hash()

	// Iterate over each vertex to compute its clustering coefficient
	for _, vertex := range vertices {
		vKey := hashFunc(vertex.Value())

		// Retrieve out-neighbors
		outNeighborsMap, exists := adjacency[vKey]
		if !exists {
			// If the vertex has no out-neighbors, its clustering coefficient is 0
			clustering[vKey] = 0.0
			continue
		}

		degree := len(outNeighborsMap)
		if degree < 2 {
			// Clustering coefficient is 0 for vertices with fewer than two out-neighbors
			clustering[vKey] = 0.0
			continue
		}

		// Collect all out-neighbor keys into a slice for easy iteration
		neighborKeys := make([]K, 0, degree)
		for neighborKey := range outNeighborsMap {
			neighborKeys = append(neighborKeys, neighborKey)
		}

		// Count the number of edges between out-neighbors
		edgeCount := 0
		for i := 0; i < len(neighborKeys); i++ {
			for j := i + 1; j < len(neighborKeys); j++ {
				u := neighborKeys[i]
				w := neighborKeys[j]

				// Check if there's an edge from u to w
				if _, connected := adjacency[u][w]; connected {
					edgeCount++
				}
			}
		}

		// Calculate the clustering coefficient
		// C(v) = edgeCount / (degree * (degree - 1))
		clustering[vKey] = float64(edgeCount) / (float64(degree) * float64(degree-1))
	}

	return clustering, nil
}
