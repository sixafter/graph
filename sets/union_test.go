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

func TestUnionDirectedWithOverlappingVertices(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(2, 3))

	unionGraph, err := Union(g, h)
	is.NoError(err)
	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(3, order, "Union graph should have 3 vertices")
	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(2, size, "Union graph should have 2 edges")
	exists, err := unionGraph.HasEdge(1, 2)
	is.True(exists, "Edge (1,2) should exist")
	exists, err = unionGraph.HasEdge(2, 3)
	is.True(exists, "Edge (2,3) should exist")
	exists, err = unionGraph.HasEdge(2, 1)
	is.False(exists, "Edge (2,1) should not exist")
}

func TestUnionDirectedWithDisjointGraphs(t *testing.T) {
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

	unionGraph, err := Union(g, h)
	is.NoError(err)
	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(4, order, "Union graph should have 4 vertices")
	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(2, size, "Union graph should have 2 edges")
	exists, err := unionGraph.HasEdge(1, 2)
	is.True(exists, "Edge (1,2) should exist")

	exists, err = unionGraph.HasEdge(3, 4)
	is.True(exists, "Edge (3,4) should exist")
}

func TestUnionDirectedWithEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddEdgeWithOptions(1, 1)) // Self-loop

	unionGraph, err := Union(g, h)
	is.NoError(err)
	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(1, order, "Union graph should have 1 vertex")
	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(1, size, "Union graph should have 1 edge")
	exists, err := unionGraph.HasEdge(1, 1)
	is.True(exists, "Self-loop (1,1) should exist")
}
