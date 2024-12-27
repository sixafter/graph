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

func TestIsSubsetBasic(t *testing.T) {
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
	is.NoError(h.AddEdgeWithOptions(2, 3))

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.True(subset, "Graph g should be a subset of graph h")
}

func TestIsSubsetMissingVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(4))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(2, 3))

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.False(subset, "Graph g should not be a subset of graph h due to missing vertex 4")
}

func TestIsSubsetMissingEdge(t *testing.T) {
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
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.False(subset, "Graph g should not be a subset of graph h due to missing edge (2,3)")
}

func TestIsSubsetEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.True(subset, "An empty graph should be a subset of any graph")
}

func TestIsSubsetIdenticalGraphs(t *testing.T) {
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

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.True(subset, "Identical graphs should be subsets of each other")
}

func TestIsSubsetDifferentTraits(t *testing.T) {
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

	subset, err := IsSubset(g, h)
	is.ErrorIs(err, graph.ErrGraphTypeMismatch)
	is.False(subset, "Graphs with different traits should not be subsets")
}

func TestIsSubsetNoCommonVertices(t *testing.T) {
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

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.False(subset, "Graphs with no common vertices should not be subsets")
}

func TestIsSubsetExtraEdgesInH(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))
	is.NoError(h.AddEdgeWithOptions(2, 1))

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.True(subset, "Graph g should be a subset of graph h even if h has extra edges")
}

func TestIsSubsetEdgeNotInH(t *testing.T) {
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
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	subset, err := IsSubset(g, h)
	is.NoError(err)
	is.False(subset, "Graph g should not be a subset of graph h because h is missing edge (2,3)")
}
