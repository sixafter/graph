// File: metrics/clustering_coefficient_test.go

// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestClusteringCoefficientCompleteGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected complete graph with 4 vertices
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges between every pair of vertices
	for i := 1; i <= 4; i++ {
		for j := i + 1; j <= 4; j++ {
			is.NoError(g.AddEdgeWithOptions(i, j))
		}
	}

	// Compute clustering coefficients
	clustering, err := ClusteringCoefficient(g)
	is.NoError(err)

	// In a complete graph, all vertices should have a clustering coefficient of 1.0
	for i := 1; i <= 4; i++ {
		k := g.Hash()(i)
		is.True(floatEquals(1.0, clustering[k]), fmt.Sprintf("Vertex %d should have a clustering coefficient of 1.0", i))
	}
}

func TestClusteringCoefficientStarGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected star graph with center vertex 1 and leaves 2, 3, 4
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges (1-2, 1-3, 1-4)
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(1, 4))

	// Compute clustering coefficients
	clustering, err := ClusteringCoefficient(g)
	is.NoError(err)

	// Clustering coefficients:
	// - Center vertex 1: No edges between leaves, so C(1) = 0
	// - Leaves: Each has only one neighbor, so C(v) = 0
	for i := 1; i <= 4; i++ {
		k := g.Hash()(i)
		is.True(floatEquals(0.0, clustering[k]), fmt.Sprintf("Vertex %d should have a clustering coefficient of 0.0", i))
	}
}

func TestClusteringCoefficientTriangle(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected triangle graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges to form a triangle
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Compute clustering coefficients
	clustering, err := ClusteringCoefficient(g)
	is.NoError(err)

	// All vertices in a triangle should have C(v) = 1.0
	for i := 1; i <= 3; i++ {
		k := g.Hash()(i)
		is.True(floatEquals(1.0, clustering[k]), fmt.Sprintf("Vertex %d should have a clustering coefficient of 1.0", i))
	}
}

func TestClusteringCoefficientDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph with two disconnected components
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	// Component 1: Vertices 1, 2, 3 forming a triangle
	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Component 2: Vertices 4, 5 forming a single edge
	for i := 4; i <= 5; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(4, 5))

	// Compute clustering coefficients
	clustering, err := ClusteringCoefficient(g)
	is.NoError(err)

	// Clustering coefficients:
	// - Vertices 1,2,3: Each in a triangle, C(v) = 1.0
	// - Vertices 4,5: Each has only one neighbor, C(v) = 0.0
	for i := 1; i <= 3; i++ {
		k := g.Hash()(i)
		is.True(floatEquals(1.0, clustering[k]), fmt.Sprintf("Vertex %d should have a clustering coefficient of 1.0", i))
	}
	for i := 4; i <= 5; i++ {
		k := g.Hash()(i)
		is.True(floatEquals(0.0, clustering[k]), fmt.Sprintf("Vertex %d should have a clustering coefficient of 0.0", i))
	}
}

func TestClusteringCoefficientSingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph with a single vertex
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))

	// Compute clustering coefficients
	clustering, err := ClusteringCoefficient(g)
	is.NoError(err)

	// Clustering coefficient for a single vertex should be 0.0
	k := g.Hash()(1)
	is.True(floatEquals(0.0, clustering[k]), "Single vertex should have a clustering coefficient of 0.0")
}

func TestClusteringCoefficientDirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New(graph.IntHash, graph.Directed())
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 1->3, 2->3, 3->1, 3->4
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Debug adjacency map
	adjacency, _ := g.AdjacencyMap()
	fmt.Printf("Adjacency map: %v\n", adjacency)

	// Compute clustering coefficients
	clustering, err := ClusteringCoefficient(g)
	is.NoError(err)

	// Debug clustering coefficients
	for vertex, coefficient := range clustering {
		fmt.Printf("Vertex %v: Clustering coefficient %.4f\n", vertex, coefficient)
	}

	// Define expected coefficients
	expected := map[int]float64{
		1: 0.5,
		2: 0.0,
		3: 0.0,
		4: 0.0,
	}

	epsilon := 1e-4

	for k, expectedVal := range expected {
		actualVal, exists := clustering[k]
		is.True(exists, fmt.Sprintf("Vertex %v should exist in the clustering map", k))
		is.True(floatApproxEqual(expectedVal, actualVal, epsilon),
			fmt.Sprintf("Vertex %v should have a clustering coefficient of %.4f, got %.4f", k, expectedVal, actualVal))
	}
}
