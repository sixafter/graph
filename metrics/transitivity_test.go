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

func TestTransitivityCompleteGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a complete graph with 4 vertices
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	for i := 1; i <= 4; i++ {
		for j := i + 1; j <= 4; j++ {
			is.NoError(g.AddEdgeWithOptions(i, j))
		}
	}

	// Compute transitivity
	transitivity, err := Transitivity(g)
	is.NoError(err)

	// In a complete graph, transitivity = 1
	is.Equal(1.0, transitivity, "Transitivity of a complete graph should be 1.0")
}

func TestTransitivityStarGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a star graph with 1 center and 3 leaves
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(1, 4))

	// Compute transitivity
	transitivity, err := Transitivity(g)
	is.NoError(err)

	// In a star graph, there are no closed triplets, so transitivity = 0
	is.Equal(0.0, transitivity, "Transitivity of a star graph should be 0.0")
}

func TestTransitivityTriangleGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a triangle graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Compute transitivity
	transitivity, err := Transitivity(g)
	is.NoError(err)

	// In a triangle graph, transitivity = 1
	is.Equal(1.0, transitivity, "Transitivity of a triangle graph should be 1.0")
}

func TestTransitivityDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a disconnected graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	// Component 1: Triangle
	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Component 2: Single edge
	is.NoError(g.AddVertexWithOptions(4))
	is.NoError(g.AddVertexWithOptions(5))
	is.NoError(g.AddEdgeWithOptions(4, 5))

	// Compute transitivity
	transitivity, err := Transitivity(g)
	is.NoError(err)

	// Transitivity should only account for component 1
	is.Equal(1.0, transitivity, "Transitivity of the disconnected graph should be 1.0")
}
