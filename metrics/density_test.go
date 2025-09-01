// Copyright (c) 2024-2025 Six After, Inc
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

func TestDensityCompleteGraph(t *testing.T) {
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

	// Compute density
	density, err := Density(g)
	is.NoError(err)

	// In a complete graph, density should be 1.0
	is.True(floatEquals(1.0, density), "Density of a complete graph should be 1.0")
}

func TestDensityStarGraph(t *testing.T) {
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

	// Compute density
	density, err := Density(g)
	is.NoError(err)

	// For a star graph with 4 vertices, density = 2|E| / |V|(|V|-1)
	expectedDensity := float64(2*3) / float64(4*3)
	is.True(floatEquals(expectedDensity, density), "Density of a star graph should match expected value")
}

func TestDensityTriangle(t *testing.T) {
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

	// Compute density
	density, err := Density(g)
	is.NoError(err)

	// In a triangle graph, density should be 1.0
	is.True(floatEquals(1.0, density), "Density of a triangle graph should be 1.0")
}

func TestDensityDisconnectedGraph(t *testing.T) {
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

	// Compute density
	density, err := Density(g)
	is.NoError(err)

	// Density for this graph is calculated as:
	// Total edges = 3 (triangle) + 1 (edge) = 4
	// Total vertices = 5
	expectedDensity := float64(2*4) / float64(5*4)
	is.True(floatEquals(expectedDensity, density), "Density of the disconnected graph should match expected value")
}

func TestDensitySingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph with a single vertex
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))

	// Compute density
	density, err := Density(g)
	is.NoError(err)

	// Density for a single vertex graph should be 0.0
	is.True(floatEquals(0.0, density), "Density of a single vertex graph should be 0.0")
}
