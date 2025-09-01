// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"github.com/sixafter/graph"
)

// DegreeCentrality computes the degree centrality for all vertices in the provided graph g.
//
// Degree centrality is a measure of the number of edges connected to a vertex.
// It is normalized by dividing by the maximum possible degree (n-1), where n is the total number of vertices in the graph.
// The resulting centrality values range from 0 to 1, where 1 indicates a vertex with the highest possible degree.
//
// This function supports both directed and undirected graphs. For directed graphs, it uses the total degree
// (sum of in-degree and out-degree) for normalization.
//
// If the input graph g is nil, contains fewer than two vertices, or if an error occurs while retrieving
// vertices or degrees, the function returns an error.
//
// Example:
//
//	centrality, err := DegreeCentrality(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute degree centrality: %v", err)
//	}
//	for vertex, score := range centrality {
//	    fmt.Printf("Vertex %v has degree centrality %.4f\n", vertex, score)
//	}
//
// Parameters:
//   - g: The graph for which to compute degree centrality. It must implement the graph.Interface[K, T] interface.
//
// Returns:
//   - A map[K]float64 where each key corresponds to a vertex in the graph, and the value represents its degree centrality score.
//   - An error if the input graph is nil, contains fewer than two vertices, or if any internal operation fails.
//
// Errors:
//   - graph.ErrNilInputGraph: Returned if the input graph g is nil.
//   - Errors propagated from g.Vertices() or g.Degree(K): Returned if there is a failure in retrieving vertices or their degrees.
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
//	    {1, 5},
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
//	// Compute degree centrality
//	centrality, err := DegreeCentrality(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute degree centrality: %v", err)
//	}
//
//	// Display centrality scores
//	for vertex, score := range centrality {
//	    fmt.Printf("Vertex %v has degree centrality %.4f\n", vertex, score)
//	}
func DegreeCentrality[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	vertices, err := g.Vertices()
	if err != nil {
		return nil, err
	}

	n := len(vertices)
	if n < 2 {
		return map[K]float64{}, nil // No meaningful centrality for graphs with fewer than 2 vertices
	}

	centrality := make(map[K]float64, n)
	hashFunc := g.Hash()
	for _, vertex := range vertices {
		hash := hashFunc(vertex.Value())
		degree, err := g.Degree(hash)
		if err != nil {
			return nil, err
		}
		// Ensure correct calculation of degree centrality
		centrality[hash] = float64(degree) / float64(n-1)
	}

	return centrality, nil
}
