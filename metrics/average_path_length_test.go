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

func TestAveragePathLengthCompleteGraph(t *testing.T) {
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

	// Compute APL
	apl, err := AveragePathLength(g)
	is.NoError(err)

	// In a complete graph, APL = 1
	is.Equal(1.0, apl, "APL of a complete graph should be 1.0")
}

func TestAveragePathLengthStarGraph(t *testing.T) {
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

	// Compute APL
	apl, err := AveragePathLength(g)
	is.NoError(err)

	// In a star graph, APL = (1+1+1+2+2+2) / 6 = 1.5
	is.Equal(1.5, apl, "APL of a star graph should be 1.5")
}

func TestAveragePathLengthDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a graph with two disconnected components
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

	// Compute APL
	_, err = AveragePathLength(g)
	is.Error(err, "Disconnected graph should return an error for APL calculation")
}

func TestAveragePathLengthSingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a graph with a single vertex
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))

	// Compute APL
	apl, err := AveragePathLength(g)
	is.NoError(err)

	// APL of a single vertex graph is 0
	is.Equal(0.0, apl, "APL of a single vertex graph should be 0.0")
}

func TestAveragePathLengthEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an empty graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	// Compute APL
	apl, err := AveragePathLength(g)
	is.NoError(err)

	// APL of an empty graph is 0
	is.Equal(0.0, apl, "APL of an empty graph should be 0.0")
}
