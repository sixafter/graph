package metrics

import (
	"fmt"

	"github.com/sixafter/graph"
)

// Modularity calculates the modularity of the given graph based on a provided community structure.
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
	fmt.Printf("Edge count: %d, Total edges (2m): %f\n", edgeCount, totalEdges)

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
			fmt.Printf("Pair (%v, %v): Degrees: (%d, %d), Aij: %f, Expected: %f\n",
				v1, v2, degreeV1, degreeV2, aij, expected)
			fmt.Printf("Expected calculation: (%d * %d) / %f = %f\n", degreeV1, degreeV2, totalEdges, expected)

			contribution := aij - expected
			modularitySum += contribution
		}
	}

	// Debugging total before normalization
	fmt.Printf("Modularity sum before normalization: %f\n", modularitySum)

	// Normalize by total edges (2m)
	modularity := modularitySum / totalEdges
	fmt.Printf("Normalized modularity: %f\n", modularity)

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
			modularitySum += (aij - expected)
		}
	}

	// Normalize by total edges (m)
	return modularitySum / totalEdges, nil
}
