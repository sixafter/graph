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

func TestIntersectionBasic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash)
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddEdgeWithOptions(2, 3))

	intersectionGraph, err := Intersection(g, h)
	is.NoError(err)

	order, err := intersectionGraph.Order()
	is.NoError(err)
	is.Equal(1, order, "Intersection graph should contain 1 vertex")

	size, err := intersectionGraph.Size()
	is.NoError(err)
	is.Equal(0, size, "Intersection graph should contain 0 edges")

	_, err = intersectionGraph.Vertex(2)
	is.NoError(err)

	_, err = intersectionGraph.Vertex(1)
	is.ErrorIs(err, graph.ErrVertexNotFound)

	_, err = intersectionGraph.Vertex(3)
	is.ErrorIs(err, graph.ErrVertexNotFound)
}

func TestIntersectionOverlappingEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2, simple.EdgeWeight(10), simple.EdgeItem("color", "blue")))

	h, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2, simple.EdgeWeight(10), simple.EdgeItem("color", "blue")))
	is.NoError(h.AddEdgeWithOptions(2, 1, simple.EdgeWeight(5), simple.EdgeItem("color", "red")))

	intersectionGraph, err := Intersection(g, h)
	is.NoError(err)

	order, err := intersectionGraph.Order()
	is.NoError(err)
	is.Equal(2, order)

	size, err := intersectionGraph.Size()
	is.NoError(err)
	is.Equal(1, size)

	hasEdge, err := intersectionGraph.HasEdge(1, 2)
	is.NoError(err)
	is.True(hasEdge)

	hasEdge, err = intersectionGraph.HasEdge(2, 1)
	is.NoError(err)
	is.False(hasEdge)
}

func TestIntersectionNoCommonVertices(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash)
	is.NoError(h.AddVertexWithOptions(3))
	is.NoError(h.AddVertexWithOptions(4))
	is.NoError(h.AddEdgeWithOptions(3, 4))

	intersectionGraph, err := Intersection(g, h)
	is.NoError(err)

	order, err := intersectionGraph.Order()
	is.NoError(err)
	is.Equal(0, order)

	size, err := intersectionGraph.Size()
	is.NoError(err)
	is.Equal(0, size)
}

func TestIntersectionTraitMismatch(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddEdgeWithOptions(1, 2))

	h, _ := simple.New(graph.IntHash)
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	intersectionGraph, err := Intersection(g, h)
	is.ErrorIs(err, graph.ErrGraphTypeMismatch)
	is.Nil(intersectionGraph)
}

func TestIntersectionSubset(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	h, _ := simple.New(graph.IntHash, graph.Directed())
	is.NoError(h.AddVertexWithOptions(1))
	is.NoError(h.AddVertexWithOptions(2))
	is.NoError(h.AddEdgeWithOptions(1, 2))

	intersectionGraph, err := Intersection(g, h)
	is.NoError(err)

	order, err := intersectionGraph.Order()
	is.NoError(err)
	is.Equal(2, order)

	size, err := intersectionGraph.Size()
	is.NoError(err)
	is.Equal(1, size)

	hasEdge, err := intersectionGraph.HasEdge(1, 2)
	is.NoError(err)
	is.True(hasEdge)

	hasEdge, err = intersectionGraph.HasEdge(2, 3)
	is.NoError(err)
	is.False(hasEdge)
}
