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

func TestEqualsIdenticalGraphs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash)

	for _, vertex := range []int{1, 2, 3} {
		is.NoError(g.AddVertexWithOptions(vertex))
		is.NoError(h.AddVertexWithOptions(vertex))
	}

	for _, edge := range [][2]int{{1, 2}, {2, 3}, {3, 1}} {
		is.NoError(g.AddEdgeWithOptions(edge[0], edge[1]))
		is.NoError(h.AddEdgeWithOptions(edge[0], edge[1]))
	}

	equals, err := Equals(g, h)
	is.NoError(err)
	is.True(equals, "Identical graphs g and h should be equal")
}

func TestEqualsDifferentTraits(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)                   // undirected
	h, _ := simple.New(graph.IntHash, graph.Directed()) // directed

	for _, vertex := range []int{1, 2} {
		is.NoError(g.AddVertexWithOptions(vertex))
		is.NoError(h.AddVertexWithOptions(vertex))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	equals, err := Equals(g, h)
	is.ErrorIs(err, graph.ErrGraphTypeMismatch)
	is.False(equals, "Graphs with different traits should not be equal")
}

func TestEqualsDifferentVertices(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash)

	for _, vertex := range []int{1, 2, 3} {
		is.NoError(g.AddVertexWithOptions(vertex))
	}
	for _, vertex := range []int{1, 2, 4} {
		is.NoError(h.AddVertexWithOptions(vertex))
	}

	equals, err := Equals(g, h)
	is.NoError(err)
	is.False(equals, "Graphs with different vertices should not be equal")
}

func TestEqualsDifferentEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash)

	for _, vertex := range []int{1, 2, 3} {
		is.NoError(g.AddVertexWithOptions(vertex))
		is.NoError(h.AddVertexWithOptions(vertex))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	is.NoError(h.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(3, 1))

	equals, err := Equals(g, h)
	is.NoError(err)
	is.False(equals, "Graphs with different edges should not be equal")
}

func TestEqualsBothEmptyGraphs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	equals, err := Equals(g, h)
	is.NoError(err)
	is.True(equals, "Both empty graphs should be equal")
}

func TestEqualsOneEmptyOneNonEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddEdgeWithOptions(1, 1)) // self-loop

	equals, err := Equals(g, h)
	is.NoError(err)
	is.False(equals, "An empty graph should not equal a non-empty graph")
}

func TestEqualsDirectedDifferentEdgeDirections(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	for _, vertex := range []int{1, 2} {
		is.NoError(g.AddVertexWithOptions(vertex))
		is.NoError(h.AddVertexWithOptions(vertex))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(2, 1))

	equals, err := Equals(g, h)
	is.NoError(err)
	is.False(equals, "Directed graphs with different edge directions should not be equal")
}

func TestEqualsUndirectedDifferentEdgeDirections(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash)

	for _, vertex := range []int{1, 2} {
		is.NoError(g.AddVertexWithOptions(vertex))
		is.NoError(h.AddVertexWithOptions(vertex))
	}

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(2, 1))

	equals, err := Equals(g, h)
	is.NoError(err)
	is.True(equals, "Undirected graphs should be equal regardless of edge direction")
}

func TestEqualsGraphsWithSelfLoops(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddEdgeWithOptions(1, 1)) // self-loop

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddEdgeWithOptions(1, 1)) // self-loop

	equals, err := Equals(g, h)
	is.NoError(err)
	is.True(equals, "Graphs with identical self-loops should be equal")
}
