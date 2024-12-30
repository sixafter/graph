package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// Modularity calculates the modularity of the given graph based on a provided community structure.
//
// Modularity is a measure of the structure of a graph, indicating the strength of division
// of the graph into communities. A higher modularity value implies that the graph has
// dense connections within communities but sparse connections between them.
//
// Parameters:
//   - g: A graph.Interface representing the graph. The graph must provide the necessary methods
//     to calculate edge weights and community-based statistics.
//   - communities: A map[K]int representing the community structure, where each key is a vertex
//     identifier, and the corresponding value is the community ID to which the vertex belongs.
//
// Returns:
//   - A float64 representing the modularity of the graph, a value between -0.5 and 1.
//   - An error if the computation cannot be performed (e.g., invalid community structure
//     or graph properties).
//
// Type Parameters:
// - K: The type of the graph's vertex keys, which must implement the graph.Ordered interface.
// - T: The type of the graph's vertex data, which can be any type.
func Modularity[K graph.Ordered, T any](g graph.Interface[K, T], communities map[K]int) (float64, error) {
	if g == nil {
		return 0, graph.ErrNilInputGraph
	}

	// Determine if the graph is directed or undirected
	traits := g.Traits()
	if traits.IsDirected {
		return calculateDirectedModularity(g, communities)
	}
	return calculateUndirectedModularity(g, communities)
}

// calculateUndirectedModularity computes the modularity for an undirected graph.
func calculateUndirectedModularity[K graph.Ordered, T any](g graph.Interface[K, T], communities map[K]int) (float64, error) {
	// Retrieve the total number of edges
	edgeCount, err := g.Size()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve edge count: %w", err)
	}
	if edgeCount == 0 {
		return 0, fmt.Errorf("graph has no edges, modularity is undefined")
	}

	totalEdges := float64(edgeCount * 2) // Normalize by 2m

	modularitySum := 0.0

	// Retrieve vertices
	vertices, err := g.Vertices()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve vertices: %w", err)
	}

	hashFunc := g.Hash()

	// Iterate over all unique vertex pairs (i, j)
	for i := 0; i < len(vertices); i++ {
		v1 := hashFunc(vertices[i].Value())
		degreeV1, err := g.Degree(v1)
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve degree for vertex %v: %w", v1, err)
		}

		for j := i + 1; j < len(vertices); j++ { // Ensure only unique pairs (i, j) where i < j
			v2 := hashFunc(vertices[j].Value())
			degreeV2, err := g.Degree(v2)
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve degree for vertex %v: %w", v2, err)
			}

			// Only consider pairs in the same community
			if communities[v1] != communities[v2] {
				continue
			}

			// Calculate modularity contribution
			aij := 0.0
			if hasEdge, _ := g.HasEdge(v1, v2); hasEdge {
				aij = 1.0
			}

			// Correct expected term using totalEdges
			expected := float64(degreeV1*degreeV2) / totalEdges

			contribution := aij - expected
			modularitySum += contribution
		}
	}

	// Normalize by total edges (2m)
	modularity := modularitySum / totalEdges
	return modularity, nil
}

// calculateDirectedModularity computes the modularity for a directed graph.
func calculateDirectedModularity[K graph.Ordered, T any](g graph.Interface[K, T], communities map[K]int) (float64, error) {
	// Retrieve the total number of edges
	edgeCount, err := g.Size()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve edge count: %w", err)
	}
	if edgeCount == 0 {
		return 0, fmt.Errorf("graph has no edges, modularity is undefined")
	}

	totalEdges := float64(edgeCount) // For directed graphs, normalize by m
	modularitySum := 0.0

	// Retrieve vertices
	vertices, err := g.Vertices()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve vertices: %w", err)
	}

	hashFunc := g.Hash()

	// Iterate over all vertex pairs (i, j)
	for i := 0; i < len(vertices); i++ {
		v1 := hashFunc(vertices[i].Value())
		outDegreeV1, err := g.OutDegree(v1)
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve out-degree for vertex %v: %w", v1, err)
		}

		for j := 0; j < len(vertices); j++ { // Directed graphs consider all pairs, including self-loops
			v2 := hashFunc(vertices[j].Value())
			inDegreeV2, err := g.InDegree(v2)
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve in-degree for vertex %v: %w", v2, err)
			}

			// Only consider pairs in the same community
			if communities[v1] != communities[v2] {
				continue
			}

			// Calculate modularity contribution
			aij := 0.0
			if hasEdge, _ := g.HasEdge(v1, v2); hasEdge {
				aij = 1.0
			}

			expected := float64(outDegreeV1*inDegreeV2) / totalEdges
			modularitySum += aij - expected
		}
	}

	// Normalize by total edges (m)
	return modularitySum / totalEdges, nil
}
