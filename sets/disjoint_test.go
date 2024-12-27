// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestIsDisjointUndirectedDisjointGraphs(t *testing.T) {
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

	disjoint, err := IsDisjoint(g, h)
	is.NoError(err)
	is.True(disjoint, "Graphs g and h should be disjoint")
}

func TestIsDisjointUndirectedNonDisjointGraphs(t *testing.T) {
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

	disjoint, err := IsDisjoint(g, h)
	is.NoError(err)
	is.False(disjoint, "Graphs g and h should not be disjoint as they share vertex 2")
}

func TestIsDisjointDirectedDisjointGraphs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddVertexWithOptions(4))
	is.NoError(h.AddEdgeWithOptions(3, 4))

	disjoint, err := IsDisjoint(g, h)
	is.NoError(err)
	is.True(disjoint, "Graphs g and h should be disjoint")
}

func TestIsDisjointDirectedNonDisjointGraphs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(2, 3))

	disjoint, err := IsDisjoint(g, h)
	is.NoError(err)
	is.False(disjoint, "Graphs g and h should not be disjoint as they share vertex 2")
}
