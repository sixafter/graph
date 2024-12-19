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

func TestTopologicalSort(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("Valid DAG", func(t *testing.T) {
		t.Parallel()
		g, _ := simple.New[int, int](graph.IntHash, graph.Directed(), graph.Acyclic())

		// Add vertices first
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 should not fail")

		// Add edges
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge 1->2 should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge 2->3 should not fail")

		order, err := TopologicalSort(g)
		is.NoError(err)
		is.Equal([]int{1, 2, 3}, order, "Order should be topological")
	})

	t.Run("Interface with Cycles", func(t *testing.T) {
		t.Parallel()
		g, _ := simple.New(graph.IntHash, graph.Directed())

		// Add vertices first
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 should not fail")

		// Add edges to form a cycle
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge 1->2 should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge 2->3 should not fail")
		err = g.AddEdgeWithOptions(3, 1)
		is.NoError(err, "Adding edge 3->1 should not fail")

		_, err = TopologicalSort(g)
		is.ErrorIs(err, graph.ErrCyclicGraph, "Should return ErrCyclicGraph for graphs with cycles")
	})

	t.Run("undirected Interface", func(t *testing.T) {
		t.Parallel()
		g, _ := simple.New(graph.IntHash)

		// Add vertices
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 should not fail")

		// Add edge
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge 1->2 should not fail")

		_, err = TopologicalSort(g)
		is.ErrorIs(err, graph.ErrUndirectedGraph, "Should return ErrUndirectedGraph for undirected graphs")
	})
}

func TestTopologicalSortDeterministic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("Valid DAG with Deterministic Order", func(t *testing.T) {
		t.Parallel()
		g, _ := simple.New(graph.IntHash, graph.Directed(), graph.Acyclic())

		// Add vertices first
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 should not fail")

		// Add edges
		err = g.AddEdgeWithOptions(1, 3)
		is.NoError(err, "Adding edge 1->3 should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge 2->3 should not fail")

		less := func(a, b int) bool {
			return a < b
		}

		order, err := TopologicalSortDeterministic(g, less)
		is.NoError(err)
		// Possible valid orders: [1,2,3] or [2,1,3]
		// With the 'less' function, [1,2,3] should be returned
		is.Equal([]int{1, 2, 3}, order, "Order should be deterministic")
	})

	t.Run("Interface with Cycles", func(t *testing.T) {
		t.Parallel()
		g, _ := simple.New(graph.IntHash, graph.Directed())

		// Add vertices first
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 should not fail")

		// Add edges to form a cycle
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge 1->2 should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge 2->3 should not fail")
		err = g.AddEdgeWithOptions(3, 1)
		is.NoError(err, "Adding edge 3->1 should not fail")

		less := func(a, b int) bool {
			return a < b
		}

		_, err = TopologicalSortDeterministic(g, less)
		is.ErrorIs(err, graph.ErrCyclicGraph, "Should return ErrCyclicGraph for graphs with cycles")
	})
}
