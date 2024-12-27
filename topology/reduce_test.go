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

func TestTransitiveReduction_RemoveRedundantEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed(), graph.Acyclic())

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1), "Adding vertex 1 should not fail")
	is.NoError(g.AddVertexWithOptions(2), "Adding vertex 2 should not fail")
	is.NoError(g.AddVertexWithOptions(3), "Adding vertex 3 should not fail")

	// Add edges
	is.NoError(g.AddEdgeWithOptions(1, 2), "Adding edge 1->2 should not fail")
	is.NoError(g.AddEdgeWithOptions(2, 3), "Adding edge 2->3 should not fail")
	is.NoError(g.AddEdgeWithOptions(1, 3), "Adding redundant edge 1->3 should not fail")

	// Perform transitive reduction
	reduced, err := TransitiveReduction(g)
	is.NoError(err, "Transitive reduction should not fail")

	// Check resulting edges
	edges, err := reduced.Edges()
	is.NoError(err, "Fetching edges should not fail")

	expectedEdges := []graph.Edge[int]{
		simple.NewEdgeWithOptions(1, 2),
		simple.NewEdgeWithOptions(2, 3),
	}
	is.ElementsMatch(expectedEdges, edges, "Redundant edge (1 -> 3) should be removed")
}

func TestTransitiveReduction_WithCycles(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash, graph.Directed())

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1), "Adding vertex 1 should not fail")
	is.NoError(g.AddVertexWithOptions(2), "Adding vertex 2 should not fail")
	is.NoError(g.AddVertexWithOptions(3), "Adding vertex 3 should not fail")

	// Add edges to form a cycle
	is.NoError(g.AddEdgeWithOptions(1, 2), "Adding edge 1->2 should not fail")
	is.NoError(g.AddEdgeWithOptions(2, 3), "Adding edge 2->3 should not fail")
	is.NoError(g.AddEdgeWithOptions(3, 1), "Adding edge 3->1 should not fail")

	// Perform transitive reduction and expect failure
	_, err := TransitiveReduction(g)
	is.ErrorIs(err, graph.ErrCyclicGraph, "Should return ErrCyclicGraph for graphs with cycles")
}

func TestTransitiveReduction_UndirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New(graph.IntHash) // undirected by default

	// Add vertices
	is.NoError(g.AddVertexWithOptions(1), "Adding vertex 1 should not fail")
	is.NoError(g.AddVertexWithOptions(2), "Adding vertex 2 should not fail")

	// Add edge
	is.NoError(g.AddEdgeWithOptions(1, 2), "Adding edge 1->2 should not fail")

	// Perform transitive reduction and expect failure
	_, err := TransitiveReduction(g)
	is.ErrorIs(err, graph.ErrUndirectedGraph, "Should return ErrUndirectedGraph for undirected graphs")
}
