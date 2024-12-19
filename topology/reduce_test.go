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

func TestTransitiveReduction(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("Remove Redundant Edges", func(t *testing.T) {
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
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge 1->2 should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge 2->3 should not fail")
		err = g.AddEdgeWithOptions(1, 3) // Redundant edge
		is.NoError(err, "Adding edge 1->3 should not fail")

		reduced, err := TransitiveReduction(g)
		is.NoError(err)

		edges, err := reduced.Edges()
		is.NoError(err)

		expectedEdges := []graph.Edge[int]{
			simple.NewEdgeWithOptions(1, 2),
			simple.NewEdgeWithOptions(2, 3),
		}
		is.ElementsMatch(expectedEdges, edges, "Redundant edge (1 -> 3) should be removed")
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

		_, err = TransitiveReduction(g)
		is.ErrorIs(err, graph.ErrCyclicGraph, "Should return ErrCyclicGraph for graphs with cycles")
	})

	t.Run("undirected Interface", func(t *testing.T) {
		t.Parallel()
		g, _ := simple.New(graph.IntHash) // undirected by default

		// Add vertices
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 should not fail")

		// Add edge
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge 1->2 should not fail")

		_, err = TransitiveReduction(g)
		is.ErrorIs(err, graph.ErrUndirectedGraph, "Should return ErrUndirectedGraph for undirected graphs")
	})
}
