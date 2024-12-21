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

func TestUnionUndirected(t *testing.T) {
	t.Parallel()

	t.Run("Union of two graphs with overlapping vertices", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first graph
		g, _ := simple.New(graph.IntHash)
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create the second graph with an overlapping vertex (2)
		h, _ := simple.New(graph.IntHash)
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to h should not fail")

		// Perform the union
		unionGraph, err := Union(g, h)
		is.NoError(err, "Union operation should not fail")

		// Validate the union graph's order (number of vertices)
		order, err := unionGraph.Order()
		is.NoError(err, "Getting order of unionGraph should not fail")
		is.Equal(3, order, "Union graph should contain 3 vertices")

		// Validate the union graph's size (number of edges)
		size, err := unionGraph.Size()
		is.NoError(err, "Getting size of unionGraph should not fail")
		is.Equal(2, size, "Union graph should contain 2 edges")

		// Check existence of edges
		edgeExists, err := unionGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.True(edgeExists, "Edge (1,2) should exist in the union graph")

		edgeExists, err = unionGraph.HasEdge(2, 3)
		is.NoError(err, "Checking existence of edge (2,3) should not fail")
		is.True(edgeExists, "Edge (2,3) should exist in the union graph")

		// Check existence of vertices
		_, err = unionGraph.Vertex(1)
		is.NoError(err, "Vertex 1 should exist in the union graph")
		_, err = unionGraph.Vertex(2)
		is.NoError(err, "Vertex 2 should exist in the union graph")
		_, err = unionGraph.Vertex(3)
		is.NoError(err, "Vertex 3 should exist in the union graph")
	})
}

func TestUnionDirected(t *testing.T) {
	t.Parallel()

	t.Run("Union of two directed graphs with overlapping vertices", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first directed graph
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding directed edge (1,2) to g should not fail")

		// Create the second directed graph with an overlapping vertex (2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding directed edge (2,3) to h should not fail")

		// Perform the union
		unionGraph, err := Union(g, h)
		is.NoError(err, "Union operation should not fail")

		// Validate the union graph's order (number of vertices)
		order, err := unionGraph.Order()
		is.NoError(err, "Getting order of unionGraph should not fail")
		is.Equal(3, order, "Union graph should contain 3 vertices")

		// Validate the union graph's size (number of edges)
		size, err := unionGraph.Size()
		is.NoError(err, "Getting size of unionGraph should not fail")
		is.Equal(2, size, "Union graph should contain 2 edges")

		// Check existence of edges
		edgeExists, err := unionGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.True(edgeExists, "Directed edge (1,2) should exist in the union graph")

		edgeExists, err = unionGraph.HasEdge(2, 3)
		is.NoError(err, "Checking existence of edge (2,3) should not fail")
		is.True(edgeExists, "Directed edge (2,3) should exist in the union graph")

		// Ensure directionality is respected
		edgeExists, err = unionGraph.HasEdge(2, 1)
		is.NoError(err, "Checking existence of edge (2,1) should not fail")
		is.False(edgeExists, "Directed edge (2,1) should not exist in the union graph")

		// Check existence of vertices
		_, err = unionGraph.Vertex(1)
		is.NoError(err, "Vertex 1 should exist in the union graph")
		_, err = unionGraph.Vertex(2)
		is.NoError(err, "Vertex 2 should exist in the union graph")
		_, err = unionGraph.Vertex(3)
		is.NoError(err, "Vertex 3 should exist in the union graph")
	})
}
