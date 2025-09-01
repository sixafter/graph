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

// DegreeCentrality tests
func TestDegreeCentralityBasic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Add vertices and edges
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	// Calculate degree centrality
	centrality, err := DegreeCentrality(g)
	is.NoError(err)

	// Verify the centrality values
	is.Equal(0.5, centrality[1], "Vertex 1 should have degree centrality 0.5")
	is.Equal(1.0, centrality[2], "Vertex 2 should have degree centrality 1.0")
	is.Equal(0.5, centrality[3], "Vertex 3 should have degree centrality 0.5")
}

func TestDegreeCentralitySingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Add a single vertex
	is.NoError(g.AddVertexWithOptions(1))

	// Calculate degree centrality
	centrality, err := DegreeCentrality(g)
	is.NoError(err)

	// Verify the centrality value
	is.Equal(0.0, centrality[1], "Single vertex should have degree centrality 0")
}

func TestDegreeCentralityEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Calculate degree centrality for an empty graph
	centrality, err := DegreeCentrality(g)
	is.NoError(err)

	// Verify that the result is an empty map
	is.Empty(centrality, "Empty graph should return an empty degree centrality map")
}

func TestDegreeCentralityUndirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)

	// Add vertices and edges
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	// Calculate degree centrality
	centrality, err := DegreeCentrality(g)
	is.NoError(err)

	// Verify the centrality values
	is.Equal(1.0/2.0, centrality[1], "Vertex 1 should have degree centrality 1/2")
	is.Equal(2.0/2.0, centrality[2], "Vertex 2 should have degree centrality 2/2")
	is.Equal(1.0/2.0, centrality[3], "Vertex 3 should have degree centrality 1/2")
}

func TestDegreeCentralityDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Add disconnected vertices
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	// Calculate degree centrality
	centrality, err := DegreeCentrality(g)
	is.NoError(err)

	// Verify the centrality values
	is.Equal(0.0, centrality[1], "Vertex 1 should have degree centrality 0")
	is.Equal(0.0, centrality[2], "Vertex 2 should have degree centrality 0")
	is.Equal(0.0, centrality[3], "Vertex 3 should have degree centrality 0")
}
