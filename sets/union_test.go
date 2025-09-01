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

func TestUnionWithNilGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Union with a nil graph
	unionGraph, err := Union(g, nil)
	is.Error(err)
	is.Nil(unionGraph, "Union result should be nil")

	// Union with both graphs nil
	unionGraph, err = Union[int, int](nil, nil)
	is.Error(err)
	is.Nil(unionGraph, "Union of two nil graphs should return nil")
}

func TestUnionDirectedWithSelfLoops(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddEdgeWithOptions(1, 1)) // Self-loop in g

	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(2, 2)) // Self-loop in h

	unionGraph, err := Union(g, h)
	is.NoError(err)

	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(2, order, "Union graph should have 2 vertices")

	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(2, size, "Union graph should have 2 edges (two self-loops)")

	exists, err := unionGraph.HasEdge(1, 1)
	is.True(exists, "Self-loop (1,1) should exist")

	exists, err = unionGraph.HasEdge(2, 2)
	is.True(exists, "Self-loop (2,2) should exist")
}

func TestUnionUndirectedWithOverlappingEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash) // Undirected by default
	h, _ := simple.New(graph.IntHash)

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
	is.Equal(2, size, "Union graph should have 2 edges (undirected edges do not duplicate)")
}

func TestUnionDirectedWithEdgeAlreadyExists(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	h, _ := simple.New(graph.IntHash, graph.Directed())

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2)) // Same edge in h

	unionGraph, err := Union(g, h)
	is.NoError(err)

	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(2, order, "Union graph should have 2 vertices")

	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(1, size, "Union graph should have 1 edge (duplicate edge should not be added)")
}

func TestUnionUndirectedWithEdgeAlreadyExists(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash)

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddEdgeWithOptions(1, 2)) // Same edge in h

	unionGraph, err := Union(g, h)
	is.NoError(err)

	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(2, order, "Union graph should have 2 vertices")

	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(1, size, "Union graph should have 1 edge (duplicate edge should not be added)")
}

func TestUnionDirectedWithMultipleComponents(t *testing.T) {
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
	is.Equal(2, size, "Union graph should have 2 edges (disjoint components)")
}

func TestUnionUndirectedWithNoEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	h, _ := simple.New(graph.IntHash)

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))

	unionGraph, err := Union(g, h)
	is.NoError(err)

	order, err := unionGraph.Order()
	is.NoError(err)
	is.Equal(2, order, "Union graph should have 2 vertices")

	size, err := unionGraph.Size()
	is.NoError(err)
	is.Equal(0, size, "Union graph should have no edges")
}
