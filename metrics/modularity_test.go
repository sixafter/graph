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

func TestModularityCompleteGraph(t *testing.T) {
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

	// Define communities: All vertices in the same community
	communities := map[int]int{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
	}

	// Compute modularity
	modularity, err := Modularity(g, communities)
	is.NoError(err)

	// Modularity for a complete graph where all vertices are in the same community should be 0.125
	is.Equal(0.125, modularity, "Modularity of a complete graph with one community should be 0.125")
}

func TestModularityTwoCommunities(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a graph with 4 vertices and two communities
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Define communities
	communities := map[int]int{
		1: 0,
		2: 0,
		3: 1,
		4: 1,
	}

	// Compute modularity
	modularity, err := Modularity(g, communities)
	is.NoError(err)

	// Modularity for this graph is non-zero because of two distinct communities
	is.True(modularity > 0, "Modularity of a two-community graph should be greater than 0")
}

func TestModularityDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a disconnected graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddVertexWithOptions(4))

	// Define communities
	communities := map[int]int{
		1: 0,
		2: 0,
		3: 1,
		4: 1,
	}

	// Compute modularity
	_, err = Modularity(g, communities)
	is.Error(err, "Modularity should be undefined for disconnected graphs with no edges")
}
