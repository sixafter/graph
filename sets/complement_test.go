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

func TestComplementBasicUndirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	complementGraph, err := Complement(g)
	is.NoError(err)
	is.NotNil(complementGraph)

	order, err := complementGraph.Order()
	is.NoError(err)
	is.Equal(3, order)

	size, err := complementGraph.Size()
	is.NoError(err)
	is.Equal(1, size)

	hasEdge, err := complementGraph.HasEdge(1, 3)
	is.NoError(err)
	is.True(hasEdge)

	hasEdge, err = complementGraph.HasEdge(1, 2)
	is.NoError(err)
	is.False(hasEdge)

	hasEdge, err = complementGraph.HasEdge(2, 3)
	is.NoError(err)
	is.False(hasEdge)
}

func TestComplementCompleteUndirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	completeGraph, err := simple.New(graph.StringHash)
	is.NoError(err)

	is.NoError(completeGraph.AddVertexWithOptions("A"))
	is.NoError(completeGraph.AddVertexWithOptions("B"))
	is.NoError(completeGraph.AddVertexWithOptions("C"))

	is.NoError(completeGraph.AddEdgeWithOptions("A", "B"))
	is.NoError(completeGraph.AddEdgeWithOptions("A", "C"))
	is.NoError(completeGraph.AddEdgeWithOptions("B", "C"))

	complementGraph, err := Complement(completeGraph)
	is.NoError(err)
	is.NotNil(complementGraph)

	order, err := complementGraph.Order()
	is.NoError(err)
	is.Equal(3, order)

	size, err := complementGraph.Size()
	is.NoError(err)
	is.Equal(0, size)
}

func TestComplementEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	emptyGraph, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(emptyGraph.AddVertexWithOptions(1))
	is.NoError(emptyGraph.AddVertexWithOptions(2))
	is.NoError(emptyGraph.AddVertexWithOptions(3))

	complementGraph, err := Complement(emptyGraph)
	is.NoError(err)
	is.NotNil(complementGraph)

	order, err := complementGraph.Order()
	is.NoError(err)
	is.Equal(3, order)

	size, err := complementGraph.Size()
	is.NoError(err)
	is.Equal(3, size)

	expectedEdges := map[[2]int]bool{
		{1, 2}: true,
		{1, 3}: true,
		{2, 3}: true,
	}

	for edge := range expectedEdges {
		hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
		is.NoError(err)
		is.True(hasEdge)
	}
}

func TestComplementBasicDirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	directedGraph, err := simple.New(graph.StringHash, graph.Directed())
	is.NoError(err)

	is.NoError(directedGraph.AddVertexWithOptions("A"))
	is.NoError(directedGraph.AddVertexWithOptions("B"))
	is.NoError(directedGraph.AddVertexWithOptions("C"))

	is.NoError(directedGraph.AddEdgeWithOptions("A", "B"))
	is.NoError(directedGraph.AddEdgeWithOptions("B", "C"))

	complementGraph, err := Complement(directedGraph)
	is.NoError(err)
	is.NotNil(complementGraph)

	order, err := complementGraph.Order()
	is.NoError(err)
	is.Equal(3, order)

	size, err := complementGraph.Size()
	is.NoError(err)
	is.Equal(4, size)

	expectedEdges := map[[2]string]bool{
		{"A", "C"}: true,
		{"B", "A"}: true,
		{"C", "A"}: true,
		{"C", "B"}: true,
	}

	for edge := range expectedEdges {
		hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
		is.NoError(err)
		is.True(hasEdge)
	}

	originalEdges := [][2]string{
		{"A", "B"},
		{"B", "C"},
	}

	for _, edge := range originalEdges {
		hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
		is.NoError(err)
		is.False(hasEdge)
	}
}

func TestComplementWithSelfLoops(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	graphWithSelfLoops, err := simple.New(graph.IntHash, graph.Directed())
	is.NoError(err)

	is.NoError(graphWithSelfLoops.AddVertexWithOptions(1))
	is.NoError(graphWithSelfLoops.AddVertexWithOptions(2))
	is.NoError(graphWithSelfLoops.AddVertexWithOptions(3))

	is.NoError(graphWithSelfLoops.AddEdgeWithOptions(1, 1))
	is.NoError(graphWithSelfLoops.AddEdgeWithOptions(1, 2))
	is.NoError(graphWithSelfLoops.AddEdgeWithOptions(2, 3))

	complementGraph, err := Complement(graphWithSelfLoops)
	is.NoError(err)
	is.NotNil(complementGraph)

	order, err := complementGraph.Order()
	is.NoError(err)
	is.Equal(3, order)

	size, err := complementGraph.Size()
	is.NoError(err)
	is.Equal(4, size)

	expectedEdges := map[[2]int]bool{
		{1, 3}: true,
		{2, 1}: true,
		{3, 1}: true,
		{3, 2}: true,
	}

	for edge := range expectedEdges {
		hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
		is.NoError(err)
		is.True(hasEdge)
	}

	originalEdges := [][2]int{
		{1, 1}, // self-loop
		{1, 2},
		{2, 3},
	}

	for _, edge := range originalEdges {
		hasEdge, err := complementGraph.HasEdge(edge[0], edge[1])
		is.NoError(err)
		is.False(hasEdge)
	}
}

func TestComplementErrorOnNilGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	complementGraph, err := Complement[int, string](nil)
	is.Error(err)
	is.Nil(complementGraph)
}
