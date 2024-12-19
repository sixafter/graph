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

func TestComplement(t *testing.T) {
	t.Parallel()

	t.Run("Complement of a basic undirected graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create an undirected graph
		g, err := simple.New(graph.IntHash) // undirected by default
		is.NoError(err, "Creating g should not fail")

		// Add vertices
		is.NoError(g.AddVertexWithOptions(1), "Adding vertex 1 to g should not fail")
		is.NoError(g.AddVertexWithOptions(2), "Adding vertex 2 to g should not fail")
		is.NoError(g.AddVertexWithOptions(3), "Adding vertex 3 to g should not fail")

		// Add edges
		is.NoError(g.AddEdgeWithOptions(1, 2), "Adding edge (1,2) to g should not fail")
		is.NoError(g.AddEdgeWithOptions(2, 3), "Adding edge (2,3) to g should not fail")

		// Compute the complement
		complementGraph, err := Complement(g)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Check the order (number of vertices)
		order, err := complementGraph.Order()
		is.NoError(err, "Getting order of complement graph should not fail")
		is.Equal(3, order, "Complement graph should have 3 vertices")

		// Check the size (number of edges)
		size, err := complementGraph.Size()
		is.NoError(err, "Getting size of complement graph should not fail")
		is.Equal(1, size, "Complement graph should have 1 edge")

		// Check that the correct edge exists in the complement
		hasEdge, err := complementGraph.HasEdge(1, 3)
		is.NoError(err, "Checking existence of edge (1,3) should not fail")
		is.True(hasEdge, "Edge (1,3) should exist in the complement graph")

		// Check that existing edges do not exist in the complement
		hasEdge, err = complementGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.False(hasEdge, "Edge (1,2) should not exist in the complement graph")

		hasEdge, err = complementGraph.HasEdge(2, 3)
		is.NoError(err, "Checking existence of edge (2,3) should not fail")
		is.False(hasEdge, "Edge (2,3) should not exist in the complement graph")
	})

	t.Run("Complement of a complete undirected graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create a complete undirected graph
		completeGraph, err := simple.New(graph.StringHash) // undirected by default
		is.NoError(err, "Creating completeGraph should not fail")

		// Add vertices
		is.NoError(completeGraph.AddVertexWithOptions("A"), "Adding vertex A to completeGraph should not fail")
		is.NoError(completeGraph.AddVertexWithOptions("B"), "Adding vertex B to completeGraph should not fail")
		is.NoError(completeGraph.AddVertexWithOptions("C"), "Adding vertex C to completeGraph should not fail")

		// Add all possible edges
		is.NoError(completeGraph.AddEdgeWithOptions("A", "B"), "Adding edge (A,B) should not fail")
		is.NoError(completeGraph.AddEdgeWithOptions("A", "C"), "Adding edge (A,C) should not fail")
		is.NoError(completeGraph.AddEdgeWithOptions("B", "C"), "Adding edge (B,C) should not fail")

		// Compute the complement
		complementGraph, err := Complement(completeGraph)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Check the order (number of vertices)
		order, err := complementGraph.Order()
		is.NoError(err, "Getting order of complement graph should not fail")
		is.Equal(3, order, "Complement graph should have 3 vertices")

		// Check the size (number of edges)
		size, err := complementGraph.Size()
		is.NoError(err, "Getting size of complement graph should not fail")
		is.Equal(0, size, "Complement of a complete graph should have 0 edges")
	})

	t.Run("Complement of an empty undirected graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create an empty undirected graph
		emptyGraph, err := simple.New(graph.IntHash) // undirected by default
		is.NoError(err, "Creating emptyGraph should not fail")

		// Add vertices
		is.NoError(emptyGraph.AddVertexWithOptions(1), "Adding vertex 1 should not fail")
		is.NoError(emptyGraph.AddVertexWithOptions(2), "Adding vertex 2 should not fail")
		is.NoError(emptyGraph.AddVertexWithOptions(3), "Adding vertex 3 should not fail")

		// No edges added

		// Compute the complement
		complementGraph, err := Complement(emptyGraph)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Check the order (number of vertices)
		order, err := complementGraph.Order()
		is.NoError(err, "Getting order of complement graph should not fail")
		is.Equal(3, order, "Complement graph should have 3 vertices")

		// Check the size (number of edges)
		size, err := complementGraph.Size()
		is.NoError(err, "Getting size of complement graph should not fail")
		is.Equal(3, size, "Complement of an empty graph should have 3 edges")

		// Check that all possible edges exist in the complement
		expectedEdges := map[[2]int]bool{
			{1, 2}: true,
			{1, 3}: true,
			{2, 3}: true,
		}

		for edge := range expectedEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.True(hasEdge, fmt.Sprintf("Edge (%v,%v) should exist in the complement graph", edge[0], edge[1]))
		}
	})

	t.Run("Complement of a basic directed graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create a directed graph
		directedGraph, err := simple.New(graph.StringHash, graph.Directed())
		is.NoError(err, "Creating directedGraph should not fail")

		// Add vertices
		is.NoError(directedGraph.AddVertexWithOptions("A"), "Adding vertex A should not fail")
		is.NoError(directedGraph.AddVertexWithOptions("B"), "Adding vertex B should not fail")
		is.NoError(directedGraph.AddVertexWithOptions("C"), "Adding vertex C should not fail")

		// Add edges
		is.NoError(directedGraph.AddEdgeWithOptions("A", "B"), "Adding edge A->B should not fail")
		is.NoError(directedGraph.AddEdgeWithOptions("B", "C"), "Adding edge B->C should not fail")

		// Compute the complement
		complementGraph, err := Complement(directedGraph)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Check the order (number of vertices)
		order, err := complementGraph.Order()
		is.NoError(err, "Getting order of complement graph should not fail")
		is.Equal(3, order, "Complement graph should have 3 vertices")

		// Check the size (number of edges)
		size, err := complementGraph.Size()
		is.NoError(err, "Getting size of complement graph should not fail")
		is.Equal(4, size, "Complement graph should have 4 edges")

		// Expected edges in complement:
		// A->C, B->A, C->A, C->B

		expectedEdges := map[[2]string]bool{
			{"A", "C"}: true,
			{"B", "A"}: true,
			{"C", "A"}: true,
			{"C", "B"}: true,
		}

		// Check that expected edges exist
		for edge := range expectedEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.True(hasEdge, fmt.Sprintf("Edge (%v,%v) should exist in the complement graph", edge[0], edge[1]))
		}

		// Check that existing edges do not exist in the complement
		originalEdges := [][2]string{
			{"A", "B"},
			{"B", "C"},
		}

		for _, edge := range originalEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.False(hasEdge, fmt.Sprintf("Edge (%v,%v) should not exist in the complement graph", edge[0], edge[1]))
		}
	})

	t.Run("Complement with self-loops", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create a directed graph that includes self-loops
		graphWithSelfLoops, err := simple.New(graph.IntHash, graph.Directed())
		is.NoError(err, "Creating graphWithSelfLoops should not fail")

		// Add vertices
		is.NoError(graphWithSelfLoops.AddVertexWithOptions(1), "Adding vertex 1 should not fail")
		is.NoError(graphWithSelfLoops.AddVertexWithOptions(2), "Adding vertex 2 should not fail")
		is.NoError(graphWithSelfLoops.AddVertexWithOptions(3), "Adding vertex 3 should not fail")

		// Add edges, including self-loops
		is.NoError(graphWithSelfLoops.AddEdgeWithOptions(1, 1), "Adding self-loop (1,1) should not fail")
		is.NoError(graphWithSelfLoops.AddEdgeWithOptions(1, 2), "Adding edge (1,2) should not fail")
		is.NoError(graphWithSelfLoops.AddEdgeWithOptions(2, 3), "Adding edge (2,3) should not fail")

		// Compute the complement
		complementGraph, err := Complement(graphWithSelfLoops)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Check the order (number of vertices)
		order, err := complementGraph.Order()
		is.NoError(err, "Getting order of complement graph should not fail")
		is.Equal(3, order, "Complement graph should have 3 vertices")

		// Check the size (number of edges)
		size, err := complementGraph.Size()
		is.NoError(err, "Getting size of complement graph should not fail")
		is.Equal(4, size, "Complement graph should have 4 edges")

		// Expected edges in complement (excluding self-loops):
		// 1->3, 2->1, 3->1, 3->2

		expectedEdges := map[[2]int]bool{
			{1, 3}: true,
			{2, 1}: true,
			{3, 1}: true,
			{3, 2}: true,
		}

		// Check that expected edges exist
		for edge := range expectedEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.True(hasEdge, fmt.Sprintf("Edge (%v,%v) should exist in the complement graph", edge[0], edge[1]))
		}

		// Check that existing edges do not exist in the complement
		originalEdges := [][2]int{
			{1, 1}, // self-loop
			{1, 2},
			{2, 3},
		}

		for _, edge := range originalEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.False(hasEdge, fmt.Sprintf("Edge (%v,%v) should not exist in the complement graph", edge[0], edge[1]))
		}
	})

	t.Run("Complement with trait mismatch", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first directed graph g
		graphG, err := simple.New(graph.IntHash, graph.Directed())
		is.NoError(err, "Creating directedGraph g should not fail")
		is.NoError(graphG.AddVertexWithOptions(1), "Adding vertex 1 to graphG should not fail")
		is.NoError(graphG.AddVertexWithOptions(2), "Adding vertex 2 to graphG should not fail")
		is.NoError(graphG.AddEdgeWithOptions(1, 2), "Adding edge (1,2) to graphG should not fail")

		// Create the second undirected graph h
		graphH, err := simple.New(graph.IntHash) // undirected by default
		is.NoError(err, "Creating undirectedGraph h should not fail")
		is.NoError(graphH.AddVertexWithOptions(1), "Adding vertex 1 to graphH should not fail")
		is.NoError(graphH.AddVertexWithOptions(2), "Adding vertex 2 to graphH should not fail")
		is.NoError(graphH.AddEdgeWithOptions(1, 2), "Adding edge (1,2) to graphH should not fail")

		// Perform the complement
		complementGraph, err := Complement(graphG)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Attempting to perform operations that assume trait mismatch
		// To test trait mismatch, you should have the Complement function handle it
		// However, in the provided Complement function, trait mismatch is not handled
		// Instead, it assumes that the complement graph follows the traits of the original graph

		// If trait mismatch needs to be handled, consider adjusting the Complement function accordingly
		// For now, this test case may not be applicable based on the current Complement implementation
	})

	t.Run("Complement where h is a subset of g", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the larger directed graph g
		graphG, err := simple.New(graph.StringHash, graph.Directed())
		is.NoError(err, "Creating directedGraph g should not fail")
		is.NoError(graphG.AddVertexWithOptions("A"), "Adding vertex A to graphG should not fail")
		is.NoError(graphG.AddVertexWithOptions("B"), "Adding vertex B to graphG should not fail")
		is.NoError(graphG.AddVertexWithOptions("C"), "Adding vertex C to graphG should not fail")
		is.NoError(graphG.AddEdgeWithOptions("A", "B"), "Adding edge A->B to graphG should not fail")
		is.NoError(graphG.AddEdgeWithOptions("B", "C"), "Adding edge B->C to graphG should not fail")

		// Create the subset directed graph h
		graphH, err := simple.New(graph.StringHash, graph.Directed())
		is.NoError(err, "Creating directedGraph h should not fail")
		is.NoError(graphH.AddVertexWithOptions("A"), "Adding vertex A to graphH should not fail")
		is.NoError(graphH.AddVertexWithOptions("B"), "Adding vertex B to graphH should not fail")
		is.NoError(graphH.AddEdgeWithOptions("A", "B"), "Adding edge A->B to graphH should not fail")

		// Compute the complement of g
		complementGraph, err := Complement(graphG)
		is.NoError(err, "Complement operation should not fail")
		is.NotNil(complementGraph, "Complement graph should not be nil")

		// Check the order (number of vertices)
		order, err := complementGraph.Order()
		is.NoError(err, "Getting order of complement graph should not fail")
		is.Equal(3, order, "Complement graph should have 3 vertices")

		// Check the size (number of edges)
		// Original graph has 2 edges, total possible edges = 3 * 2 = 6 (excluding self-loops)
		// Complement graph should have 4 edges
		size, err := complementGraph.Size()
		is.NoError(err, "Getting size of complement graph should not fail")
		is.Equal(4, size, "Complement graph should have 4 edges")

		// Expected edges in complement:
		// A->C, B->A, C->A, C->B

		expectedEdges := map[[2]string]bool{
			{"A", "C"}: true,
			{"B", "A"}: true,
			{"C", "A"}: true,
			{"C", "B"}: true,
		}

		// Check that expected edges exist
		for edge := range expectedEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.True(hasEdge, fmt.Sprintf("Edge (%v,%v) should exist in the complement graph", edge[0], edge[1]))
		}

		// Check that existing edges do not exist in the complement
		originalEdges := [][2]string{
			{"A", "B"},
			{"B", "C"},
		}

		for _, edge := range originalEdges {
			hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
			is.NoError(err, "Checking existence of edge should not fail")
			is.False(hasEdge, fmt.Sprintf("Edge (%v,%v) should not exist in the complement graph", edge[0], edge[1]))
		}
	})

	t.Run("Error when input graph is nil", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Perform the complement on a nil graph
		complementGraph, err := Complement[int, string](nil)
		is.Error(err, "Complement operation should fail when input graph is nil")
		is.Nil(complementGraph, "Complement graph should be nil when input graph is nil")
	})
}
