// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"fmt"
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestDifference(t *testing.T) {
	t.Parallel()

	t.Run("Basic Difference of two graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create graph h
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create graph h
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to h should not fail")

		// Perform the difference h - h
		differenceGraph, err := Difference(g, h)
		is.NoError(err, "Difference operation should not fail")

		// Debugging: Print adjacency maps for clarity
		adjMapG, err := g.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of g should not fail")
		adjMapH, err := h.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of h should not fail")
		adjMapDifference, err := differenceGraph.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of differenceGraph should not fail")

		fmt.Printf("GraphG adjacency map: %v\n", adjMapG)
		fmt.Printf("GraphH adjacency map: %v\n", adjMapH)
		fmt.Printf("Difference graph adjacency map: %v\n", adjMapDifference)

		// Validate the difference graph's order (number of vertices)
		order, err := differenceGraph.Order()
		is.NoError(err, "Getting order of differenceGraph should not fail")
		is.Equal(2, order, "Difference graph should contain 2 vertices")

		// Validate the difference graph's size (number of edges)
		size, err := differenceGraph.Size()
		is.NoError(err, "Getting size of differenceGraph should not fail")
		is.Equal(1, size, "Difference graph should contain 1 edge")

		// Check existence of edges
		hasEdge, err := differenceGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.True(hasEdge, "Edge (1,2) should exist in the difference graph")

		hasEdge, err = differenceGraph.HasEdge(2, 3)
		is.NoError(err, "Checking existence of edge (2,3) should not fail")
		is.False(hasEdge, "Edge (2,3) should not exist in the difference graph")
	})

	t.Run("Difference with overlapping edges", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create directed graph h with overlapping edges
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")
		err = h.AddEdgeWithOptions(2, 1)
		is.NoError(err, "Adding edge (2,1) to h should not fail")

		// Perform the difference g - h
		differenceGraph, err := Difference(g, h)
		is.NoError(err, "Difference operation should not fail")

		// Debugging: Print adjacency maps for clarity
		adjMapG, err := g.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of g should not fail")
		adjMapH, err := h.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of h should not fail")
		adjMapDifference, err := differenceGraph.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of differenceGraph should not fail")

		fmt.Printf("GraphG adjacency map: %v\n", adjMapG)
		fmt.Printf("GraphH adjacency map: %v\n", adjMapH)
		fmt.Printf("Difference graph adjacency map: %v\n", adjMapDifference)

		// Validate the difference graph's order (number of vertices)
		order, err := differenceGraph.Order()
		is.NoError(err, "Getting order of differenceGraph should not fail")
		is.Equal(2, order, "Difference graph should contain 2 vertices")

		// Validate the difference graph's size (number of edges)
		size, err := differenceGraph.Size()
		is.NoError(err, "Getting size of differenceGraph should not fail")
		is.Equal(0, size, "Difference graph should contain 0 edges")

		// Check existence of edges
		hasEdge, err := differenceGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.False(hasEdge, "Edge (1,2) should not exist in the difference graph")

		hasEdge, err = differenceGraph.HasEdge(2, 1)
		is.NoError(err, "Checking existence of edge (2,1) should not fail")
		is.False(hasEdge, "Edge (2,1) should not exist in the difference graph")
	})

	t.Run("Difference with no common vertices", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first undirected graph g
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create the second undirected graph h
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to h should not fail")
		err = h.AddEdgeWithOptions(3, 4)
		is.NoError(err, "Adding edge (3,4) to h should not fail")

		// Perform the difference g - h
		differenceGraph, err := Difference(g, h)
		is.NoError(err, "Difference operation should not fail")

		// Debugging: Print adjacency maps for clarity
		adjMapG, err := g.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of g should not fail")
		adjMapH, err := h.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of h should not fail")
		adjMapDifference, err := differenceGraph.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of differenceGraph should not fail")

		fmt.Printf("GraphG adjacency map: %v\n", adjMapG)
		fmt.Printf("GraphH adjacency map: %v\n", adjMapH)
		fmt.Printf("Difference graph adjacency map: %v\n", adjMapDifference)

		// Validate the difference graph's order (number of vertices)
		order, err := differenceGraph.Order()
		is.NoError(err, "Getting order of differenceGraph should not fail")
		is.Equal(2, order, "Difference graph should contain 2 vertices")

		// Validate the difference graph's size (number of edges)
		size, err := differenceGraph.Size()
		is.NoError(err, "Getting size of differenceGraph should not fail")
		is.Equal(1, size, "Difference graph should contain 1 edge")

		// Check existence of edges
		hasEdge, err := differenceGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.True(hasEdge, "Edge (1,2) should exist in the difference graph")

		// For undirected graphs, (2,1) is equivalent to (1,2)
		hasEdge, err = differenceGraph.HasEdge(2, 1)
		is.NoError(err, "Checking existence of edge (2,1) should not fail")
		is.True(hasEdge, "Edge (2,1) should exist in the difference graph as it is undirected")
	})

	t.Run("Difference with trait mismatch", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create the second undirected graph h
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Perform the difference g - h
		differenceGraph, err := Difference(g, h)
		is.ErrorIs(err, graph.ErrGraphTypeMismatch, "Difference should fail due to trait mismatch")
		is.Nil(differenceGraph, "Difference graph should be nil due to trait mismatch")
	})

	t.Run("Difference where h is a subset of g", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the larger directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to g should not fail")

		// Create the subset directed graph h
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Perform the difference g - h
		differenceGraph, err := Difference(g, h)
		is.NoError(err, "Difference operation should not fail")

		// Debugging: Print adjacency maps for clarity
		adjMapG, err := g.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of g should not fail")
		adjMapH, err := h.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of h should not fail")
		adjMapDifference, err := differenceGraph.AdjacencyMap()
		is.NoError(err, "Getting adjacency map of differenceGraph should not fail")

		fmt.Printf("GraphG adjacency map: %v\n", adjMapG)
		fmt.Printf("GraphH adjacency map: %v\n", adjMapH)
		fmt.Printf("Difference graph adjacency map: %v\n", adjMapDifference)

		// Validate the difference graph's order (number of vertices)
		order, err := differenceGraph.Order()
		is.NoError(err, "Getting order of differenceGraph should not fail")
		is.Equal(3, order, "Difference graph should contain 3 vertices")

		// Validate the difference graph's size (number of edges)
		size, err := differenceGraph.Size()
		is.NoError(err, "Getting size of differenceGraph should not fail")
		is.Equal(1, size, "Difference graph should contain 1 edge")

		// Check existence of edges
		hasEdge, err := differenceGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.False(hasEdge, "Edge (1,2) should not exist in the difference graph")

		hasEdge, err = differenceGraph.HasEdge(2, 3)
		is.NoError(err, "Checking existence of edge (2,3) should not fail")
		is.True(hasEdge, "Edge (2,3) should exist in the difference graph")

		// Ensure all vertices from g are present
		_, err = differenceGraph.Vertex(1)
		is.NoError(err, "Vertex 1 should still exist in the difference graph")
		_, err = differenceGraph.Vertex(2)
		is.NoError(err, "Vertex 2 should still exist in the difference graph")
		_, err = differenceGraph.Vertex(3)
		is.NoError(err, "Vertex 3 should still exist in the difference graph")
	})
}

func TestSymmetricDifferenceDirected(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Create two test graphs g and h
	g, _ := simple.New(graph.StringHash, graph.Directed())
	h, _ := simple.New(graph.StringHash, graph.Directed())

	// Add vertices and edges to graph g
	is.NoError(g.AddVertexWithOptions("A"), "Failed to add vertex A to graph g")
	is.NoError(g.AddVertexWithOptions("B"), "Failed to add vertex B to graph g")
	is.NoError(g.AddEdgeWithOptions("A", "B"), "Failed to add edge A->B to graph g")

	// Add vertices and edges to graph h
	is.NoError(h.AddVertexWithOptions("B"), "Failed to add vertex B to graph h")
	is.NoError(h.AddVertexWithOptions("C"), "Failed to add vertex C to graph h")
	is.NoError(h.AddEdgeWithOptions("B", "C"), "Failed to add edge B->C to graph h")

	// Compute the symmetric difference
	symmetricGraph, err := SymmetricDifference(g, h)

	// Assert no errors occurred
	is.NoError(err, "Failed to compute the symmetric difference between g and h")
	is.NotNil(symmetricGraph, "Symmetric difference graph is nil")

	// Debugging: print all vertices and edges of the symmetric graph
	vertices, vertexErr := symmetricGraph.Order()
	if vertexErr == nil {
		fmt.Printf("Symmetric Interface Vertex Count: %d\n", vertices)
	} else {
		fmt.Printf("Error fetching vertex cardinality: %v\n", vertexErr)
	}

	edges, edgeErr := symmetricGraph.Size()
	if edgeErr == nil {
		fmt.Printf("Symmetric Interface Edge Count: %d\n", edges)
	} else {
		fmt.Printf("Error fetching edge cardinality: %v\n", edgeErr)
	}

	// Verify the vertices in the symmetric difference
	exists, err := symmetricGraph.HasVertex("A")
	is.NoError(err, "Error checking existence of vertex A in symmetric difference graph")
	is.True(exists, "Vertex A should be included in the symmetric difference")

	exists, err = symmetricGraph.HasVertex("C")
	is.NoError(err, "Error checking existence of vertex C in symmetric difference graph")
	is.True(exists, "Vertex C should be included in the symmetric difference")

	// Corrected assertion for vertex "B"
	exists, err = symmetricGraph.HasVertex("B")
	is.NoError(err, "Error checking existence of vertex B in symmetric difference graph")
	is.True(exists, "Vertex B should be included in the symmetric difference")

	// Verify the edges in the symmetric difference
	hasEdge, err := symmetricGraph.HasEdge("A", "B")
	is.NoError(err, "Error checking existence of edge A->B in symmetric difference graph")
	is.True(hasEdge, "Edge A->B should be included in the symmetric difference")

	hasEdge, err = symmetricGraph.HasEdge("B", "C")
	is.NoError(err, "Error checking existence of edge B->C in symmetric difference graph")
	is.True(hasEdge, "Edge B->C should be included in the symmetric difference")
}

func TestSymmetricDifferenceUndirected(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Create two test graphs g and h
	g, _ := simple.New(graph.StringHash)
	h, _ := simple.New(graph.StringHash)

	// Add vertices and edges to graph g
	is.NoError(g.AddVertexWithOptions("A"), "Failed to add vertex A to graph g")
	is.NoError(g.AddVertexWithOptions("B"), "Failed to add vertex B to graph g")
	is.NoError(g.AddEdgeWithOptions("A", "B"), "Failed to add edge A-B to graph g")

	// Add vertices and edges to graph h
	is.NoError(h.AddVertexWithOptions("B"), "Failed to add vertex B to graph h")
	is.NoError(h.AddVertexWithOptions("C"), "Failed to add vertex C to graph h")
	is.NoError(h.AddEdgeWithOptions("B", "C"), "Failed to add edge B-C to graph h")

	// Compute the symmetric difference
	symmetricGraph, err := SymmetricDifference(g, h)

	// Assert no errors occurred
	is.NoError(err, "Failed to compute the symmetric difference between g and h")
	is.NotNil(symmetricGraph, "Symmetric difference graph is nil")
}
