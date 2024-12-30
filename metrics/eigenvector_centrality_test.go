// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestEigenvectorCentralityUndirectedStarGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph
	g, err := simple.New[int, int](graph.IntHash)
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges to form a star graph
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(1, 4))

	// Compute eigenvector centrality
	centrality, err := EigenvectorCentrality(g)
	is.NoError(err)

	// Expected centrality:
	// Vertex 1 should have the highest centrality
	// Vertices 2, 3, 4 should have equal lower centrality

	center := centrality[1]
	leaves := []float64{centrality[2], centrality[3], centrality[4]}

	// Center should be significantly higher than leaves
	is.True(center > 0.5, "Center vertex should have a high centrality score")
	for _, leaf := range leaves {
		is.True(leaf < center, "Leaf vertices should have lower centrality scores than the center")
	}
}

func TestEigenvectorCentralityDirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 1->3, 2->3, 3->1, 3->4
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Compute eigenvector centrality
	centrality, err := EigenvectorCentrality(g)
	is.NoError(err)

	// Expected centrality:
	// Vertex 3 should have the highest centrality
	// Vertex 1 should have the second highest
	// Vertex 4 should have higher centrality than Vertex 2

	c3 := centrality[3]
	c1 := centrality[1]
	c2 := centrality[2]
	c4 := centrality[4]

	is.True(c3 > c1, "Vertex 3 should have a higher centrality than Vertex 1")
	is.True(c1 > c2, "Vertex 1 should have a higher centrality than Vertex 2")
	is.True(c4 > c2, "Vertex 4 should have a higher centrality than Vertex 2")
}

func TestEigenvectorCentralityDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph with two disconnected components
	g, err := simple.New[int, int](graph.IntHash)
	is.NoError(err)

	// Component 1: Triangle (Vertices 1, 2, 3)
	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Component 2: Single edge (Vertices 4, 5)
	for i := 4; i <= 5; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(4, 5))

	// Compute eigenvector centrality
	centrality, err := EigenvectorCentrality(g)
	is.NoError(err)

	// Expected centrality:
	// Component 1: All vertices should have equal centrality
	// Component 2: Vertices 4 and 5 should have equal centrality
	c1 := centrality[1]
	c2 := centrality[2]
	c3 := centrality[3]
	c4 := centrality[4]
	c5 := centrality[5]

	// Verify that Component 1 vertices have equal centrality
	is.True(floatEquals(c1, c2), "Vertices 1 and 2 should have equal centrality")
	is.True(floatEquals(c2, c3), "Vertices 2 and 3 should have equal centrality")

	// Verify that Component 2 vertices have equal centrality
	is.True(floatEquals(c4, c5), "Vertices 4 and 5 should have equal centrality")

	// Verify that Component 1 centrality is higher than Component 2
	is.True(c1 > c4, "Centrality of Component 1 should be higher than Component 2")
}

func TestEigenvectorCentralitySelfLoopGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph with a self-loop
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 2; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 2->1, 2->2 (self-loop)
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 1))
	is.NoError(g.AddEdgeWithOptions(2, 2))

	// Compute eigenvector centrality
	centrality, err := EigenvectorCentrality(g)
	is.NoError(err)

	// Expected centrality:
	// Vertex 2 should have higher centrality due to the self-loop
	// Vertex 1 should have lower centrality
	c1 := centrality[1]
	c2 := centrality[2]

	is.True(c2 > c1, "Vertex 2 should have a higher centrality than Vertex 1 due to the self-loop")
}
