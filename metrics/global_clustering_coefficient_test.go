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

func TestGlobalClusteringCoefficientCompleteGraph(t *testing.T) {
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

	// Compute global clustering coefficient
	coeff, err := GlobalClusteringCoefficient(g)
	is.NoError(err)

	// In a complete graph, the global clustering coefficient should be 1.0
	is.True(floatEquals(1.0, coeff), "Global clustering coefficient of a complete graph should be 1.0")
}

func TestGlobalClusteringCoefficientStarGraph(t *testing.T) {
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

	// Compute global clustering coefficient
	coeff, err := GlobalClusteringCoefficient(g)
	is.NoError(err)

	// In a star graph, the global clustering coefficient should be 0.0
	is.True(floatEquals(0.0, coeff), "Global clustering coefficient of a star graph should be 0.0")
}

func TestGlobalClusteringCoefficientTriangle(t *testing.T) {
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

	// Compute global clustering coefficient
	coeff, err := GlobalClusteringCoefficient(g)
	is.NoError(err)

	// In a triangle graph, the global clustering coefficient should be 1.0
	is.True(floatEquals(1.0, coeff), "Global clustering coefficient of a triangle graph should be 1.0")
}

func TestGlobalClusteringCoefficientDisconnectedGraph(t *testing.T) {
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

	// Compute global clustering coefficient
	coeff, err := GlobalClusteringCoefficient(g)
	is.NoError(err)

	// Global clustering coefficient should be 1.0 due to the triangle
	is.True(floatEquals(1.0, coeff), "Global clustering coefficient of the graph should be 1.0")
}

func TestGlobalClusteringCoefficientSingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph with a single vertex
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))

	// Compute global clustering coefficient
	coeff, err := GlobalClusteringCoefficient(g)
	is.NoError(err)

	// Global clustering coefficient for a single vertex should be 0.0
	is.True(floatEquals(0.0, coeff), "Global clustering coefficient of a single vertex graph should be 0.0")
}
