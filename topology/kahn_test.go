// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package topology

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestTopologicalSortValidDAG(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash, graph.Directed(), graph.Acyclic())

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	// Add edges
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	order, err := TopologicalSort(g)
	is.NoError(err)
	is.Equal([]int{1, 2, 3}, order, "Order should be topological")
}

func TestTopologicalSortGraphWithCycles(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	// Add edges to form a cycle
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	_, err := TopologicalSort(g)
	is.ErrorIs(err, graph.ErrCyclicGraph, "Should return ErrCyclicGraph for graphs with cycles")
}

func TestTopologicalSortUndirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash)

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))

	// Add edge
	is.NoError(g.AddEdgeWithOptions(1, 2))

	_, err := TopologicalSort(g)
	is.ErrorIs(err, graph.ErrUndirectedGraph, "Should return ErrUndirectedGraph for undirected graphs")
}

func TestTopologicalSortDeterministicValidDAG(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed(), graph.Acyclic())

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	// Add edges
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	less := func(a, b int) bool {
		return a < b
	}

	order, err := TopologicalSortDeterministic(g, less)
	is.NoError(err)
	is.Equal([]int{1, 2, 3}, order, "Order should be deterministic")
}

func TestTopologicalSortDeterministicGraphWithCycles(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	// Add edges to form a cycle
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	less := func(a, b int) bool {
		return a < b
	}

	_, err := TopologicalSortDeterministic(g, less)
	is.ErrorIs(err, graph.ErrCyclicGraph, "Should return ErrCyclicGraph for graphs with cycles")
}
