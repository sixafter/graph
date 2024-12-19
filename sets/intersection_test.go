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

func TestIntersection(t *testing.T) {
	t.Parallel()

	t.Run("Basic Intersection of two graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first undirected graph
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create the second undirected graph
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to h should not fail")

		// Perform the intersection
		intersectionGraph, err := Intersection(g, h)
		is.NoError(err, "Intersection operation should not fail")

		// Validate the intersection graph's order (number of vertices)
		order, err := intersectionGraph.Order()
		is.NoError(err, "Getting order of intersectionGraph should not fail")
		is.Equal(1, order, "Intersection graph should contain 1 vertex")

		// Validate the intersection graph's size (number of edges)
		size, err := intersectionGraph.Size()
		is.NoError(err, "Getting size of intersectionGraph should not fail")
		is.Equal(0, size, "Intersection graph should contain 0 edges")

		// Check existence of vertices
		_, err = intersectionGraph.Vertex(2)
		is.NoError(err, "Vertex 2 should exist in the intersection graph")

		_, err = intersectionGraph.Vertex(1)
		is.ErrorIs(err, graph.ErrVertexNotFound, "Vertex 1 should not exist in the intersection graph")

		_, err = intersectionGraph.Vertex(3)
		is.ErrorIs(err, graph.ErrVertexNotFound, "Vertex 3 should not exist in the intersection graph")
	})

	t.Run("Intersection with overlapping edges", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first directed graph
		g, err := simple.New(graph.IntHash, graph.Directed())
		is.NoError(err, "Creating g should not fail")
		err = g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2, simple.EdgeWeight(10), simple.EdgeItem("color", "blue"))
		is.NoError(err, "Adding edge (1->2) to g should not fail")

		// Create the second directed graph
		h, err := simple.New(graph.IntHash, graph.Directed())
		is.NoError(err, "Creating h should not fail")
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2, simple.EdgeWeight(10), simple.EdgeItem("color", "blue"))
		is.NoError(err, "Adding edge (1->2) to h should not fail")
		err = h.AddEdgeWithOptions(2, 1, simple.EdgeWeight(5), simple.EdgeItem("color", "red")) // Adding reverse edge
		is.NoError(err, "Adding edge (2->1) to h should not fail")

		// Perform the intersection
		intersectionGraph, err := Intersection(g, h)
		is.NoError(err, "Intersection operation should not fail")

		// Validate the intersection graph's order (number of vertices)
		order, err := intersectionGraph.Order()
		is.NoError(err, "Getting order of intersectionGraph should not fail")
		is.Equal(2, order, "Intersection graph should contain 2 vertices")

		// Validate the intersection graph's size (number of edges)
		size, err := intersectionGraph.Size()
		is.NoError(err, "Getting size of intersectionGraph should not fail")
		is.Equal(1, size, "Intersection graph should contain 1 edge")

		// Check existence of edges
		hasEdge, err := intersectionGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1->2) should not fail")
		is.True(hasEdge, "Edge (1->2) should exist in the intersection graph")

		hasEdge, err = intersectionGraph.HasEdge(2, 1)
		is.NoError(err, "Checking existence of edge (2->1) should not fail")
		is.False(hasEdge, "Edge (2->1) should not exist in the intersection graph")
	})

	t.Run("Intersection with no common vertices", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first undirected graph
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2) // Add missing vertex 2
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create the second undirected graph
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddVertexWithOptions(4) // Add missing vertex 4
		is.NoError(err, "Adding vertex 4 to h should not fail")
		err = h.AddEdgeWithOptions(3, 4)
		is.NoError(err, "Adding edge (3,4) to h should not fail")

		// Perform the intersection
		intersectionGraph, err := Intersection(g, h)
		is.NoError(err, "Intersection operation should not fail")

		// Validate the intersection graph's order (number of vertices)
		order, err := intersectionGraph.Order()
		is.NoError(err, "Getting order of intersectionGraph should not fail")
		is.Equal(0, order, "Intersection graph should contain 0 vertices")

		// Validate the intersection graph's size (number of edges)
		size, err := intersectionGraph.Size()
		is.NoError(err, "Getting size of intersectionGraph should not fail")
		is.Equal(0, size, "Intersection graph should contain 0 edges")
	})

	t.Run("Intersection with trait mismatch", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the first directed graph
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create the second undirected graph
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Perform the intersection
		intersectionGraph, err := Intersection(g, h)
		is.ErrorIs(err, graph.ErrGraphTypeMismatch, "Intersection should fail due to trait mismatch")
		is.Nil(intersectionGraph, "Intersection graph should be nil due to trait mismatch")
	})

	t.Run("Intersection where one graph is subset of the other", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create the larger directed graph
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to g should not fail")

		// Create the subset directed graph
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Perform the intersection
		intersectionGraph, err := Intersection(g, h)
		is.NoError(err, "Intersection operation should not fail")

		// Validate the intersection graph's order (number of vertices)
		order, err := intersectionGraph.Order()
		is.NoError(err, "Getting order of intersectionGraph should not fail")
		is.Equal(2, order, "Intersection graph should contain 2 vertices")

		// Validate the intersection graph's size (number of edges)
		size, err := intersectionGraph.Size()
		is.NoError(err, "Getting size of intersectionGraph should not fail")
		is.Equal(1, size, "Intersection graph should contain 1 edge")

		// Check existence of edges
		hasEdge, err := intersectionGraph.HasEdge(1, 2)
		is.NoError(err, "Checking existence of edge (1,2) should not fail")
		is.True(hasEdge, "Edge (1,2) should exist in the intersection graph")

		hasEdge, err = intersectionGraph.HasEdge(2, 3)
		is.NoError(err, "Checking existence of edge (2,3) should not fail")
		is.False(hasEdge, "Edge (2,3) should not exist in the intersection graph")
	})
}
