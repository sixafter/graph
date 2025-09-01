// Copyright (c) 2024-2025 Six After, Inc
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

func TestIsSupersetBasic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.True(superset, "Graph g should be a superset of graph h")
}

func TestIsSupersetMissingVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.False(superset, "Graph g should not be a superset of graph h due to missing vertex 3")
}

func TestIsSupersetMissingEdge(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.True(superset, "Graph g should be a superset of graph h as h has no edges")
}

func TestIsSupersetEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.True(superset, "Any graph g should be a superset of an empty graph h")
}

func TestIsSupersetIdenticalGraphs(t *testing.T) {
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

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.True(superset, "Identical graphs g and h should satisfy IsSuperset")
}

func TestIsSupersetDifferentTraits(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	superset, err := IsSuperset(g, h)
	is.ErrorIs(err, graph.ErrGraphTypeMismatch)
	is.False(superset, "Graphs with different traits should not satisfy IsSuperset")
}

func TestIsSupersetNoCommonVertices(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddVertexWithOptions(4))
	is.NoError(h.AddEdgeWithOptions(3, 4))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.False(superset, "Graphs with no common vertices should not satisfy IsSuperset")
}

func TestIsSupersetWithExtraVerticesAndEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddVertexWithOptions(4))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.True(superset, "Graph g should be a superset of graph h")
}

func TestIsSupersetDirectedGraphs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	superset, err := IsSuperset(g, h)
	is.NoError(err)
	is.True(superset, "Directed graph g should be a superset of graph h")
}
