// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/traverse"
)

// ClosenessCentrality calculates the closeness centrality for each vertex in the given graph.
//
// The closeness centrality is a measure of how close a vertex is to all other reachable vertices
// in the graph. It is defined as the reciprocal of the sum of the shortest path distances from the
// vertex to all other reachable vertices. A higher closeness centrality indicates that a vertex
// can quickly interact with all other vertices in the graph.
//
// The closeness centrality for a vertex v is defined as:
//
//	C(v) = (number of reachable vertices - 1) / sum of shortest path distances from v to all reachable vertices
//
// If a vertex has no reachable vertices (isolated vertex), its closeness centrality is defined as 0.
//
// The function returns a map where each key corresponds to a vertex in the graph,
// and the value represents its closeness centrality.
//
// Example:
//
//	centrality, err := ClosenessCentrality(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute closeness centrality: %v", err)
//	}
//	for vertex, score := range centrality {
//	    fmt.Printf("Vertex %v has a closeness centrality of %.4f\n", vertex, score)
//	}
//
// Parameters:
//   - g: The graph for which to compute closeness centrality. It must implement the graph.Interface[K, T] interface.
//
// Returns:
//   - A map[K]float64 where each key corresponds to a vertex in the graph,
//     and the value represents its closeness centrality.
//   - An error if the input graph is nil, or if there is a failure in retrieving vertices or adjacency maps.
//
// Errors:
//   - graph.ErrNilInputGraph: Returned if the input graph g is nil.
//   - Errors propagated from traverse.BFSWithDepthTracking or other graph methods.
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
//	// Compute closeness centrality
//	centrality, err := ClosenessCentrality(g)
//	if err != nil {
//	    log.Fatalf("Failed to compute closeness centrality: %v", err)
//	}
//
//	// Display closeness centrality
//	for vertex, score := range centrality {
//	    fmt.Printf("Vertex %v has a closeness centrality of %.4f\n", vertex, score)
//	}
func ClosenessCentrality[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Validate input graph
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	// Retrieve all vertices in the graph
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("closeness_centrality: failed to retrieve vertices: %w", err)
	}

	// Initialize the closeness centrality map
	centrality := make(map[K]float64, len(vertices))

	// Iterate over each vertex to compute its closeness centrality
	for _, vertex := range vertices {
		vKey := vertex.ID()

		// Initialize variables to track sum of distances and number of reachable vertices
		var sumDistances int
		var numReachable int

		// Define the visit function to accumulate distances
		visit := func(uKey K, depth int) bool {
			if uKey == vKey {
				// Skip the source vertex itself
				return false // Continue BFS
			}
			sumDistances += depth
			numReachable++
			return false // Continue BFS
		}

		// Perform BFS with depth tracking from the current vertex
		err := traverse.BFSWithDepthTracking(g, vKey, visit)
		if err != nil {
			return nil, fmt.Errorf("closeness_centrality: BFS failed for vertex %v: %w", vKey, err)
		}

		// Compute closeness centrality
		if sumDistances > 0 && numReachable > 0 {
			centrality[vKey] = float64(numReachable) / float64(sumDistances)
		} else {
			centrality[vKey] = 0.0
		}
	}

	return centrality, nil
}
