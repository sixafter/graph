// Copyright (c) 2024-2025 Six After, Inc
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

func TestDifferenceBasic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash) // undirected by default
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash) // undirected by default
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(2, 3))

	differenceGraph, err := Difference(g, h)
	is.NoError(err)
	is.NotNil(differenceGraph)

	order, err := differenceGraph.Order()
	is.NoError(err)
	is.Equal(2, order)

	size, err := differenceGraph.Size()
	is.NoError(err)
	is.Equal(1, size)

	hasEdge, err := differenceGraph.HasEdge(1, 2)
	is.NoError(err)
	is.True(hasEdge)
}

func TestDifferenceOverlappingEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(2, 1))

	differenceGraph, err := Difference(g, h)
	is.NoError(err)
	is.NotNil(differenceGraph)

	size, err := differenceGraph.Size()
	is.NoError(err)
	is.Equal(0, size)
}

func TestDifferenceNoCommonVertices(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash) // undirected by default
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash) // undirected by default
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddVertexWithOptions(4))
	is.NoError(h.AddEdgeWithOptions(3, 4))

	differenceGraph, err := Difference(g, h)
	is.NoError(err)
	is.NotNil(differenceGraph)

	order, err := differenceGraph.Order()
	is.NoError(err)
	is.Equal(2, order)

	size, err := differenceGraph.Size()
	is.NoError(err)
	is.Equal(1, size)

	hasEdge, err := differenceGraph.HasEdge(1, 2)
	is.NoError(err)
	is.True(hasEdge)
}

func TestDifferenceTraitMismatch(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash) // undirected by default
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	differenceGraph, err := Difference(g, h)
	is.Error(err)
	is.Nil(differenceGraph)
}

func TestDifferenceSubset(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	h, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	differenceGraph, err := Difference(g, h)
	is.NoError(err)
	is.NotNil(differenceGraph)

	order, err := differenceGraph.Order()
	is.NoError(err)
	is.Equal(3, order)

	size, err := differenceGraph.Size()
	is.NoError(err)
	is.Equal(1, size)

	hasEdge, err := differenceGraph.HasEdge(2, 3)
	is.NoError(err)
	is.True(hasEdge)
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
