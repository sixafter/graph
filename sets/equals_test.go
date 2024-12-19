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

// TestEquals tests the Equals function across various scenarios.
func TestEquals(t *testing.T) {
	t.Parallel()

	t.Run("Identical Graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create undirected graph g with vertices {1, 2, 3} and edges (1,2), (2,3), (3,1)
		g, _ := simple.New(graph.IntHash) // undirected by default
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
		err = g.AddEdgeWithOptions(3, 1)
		is.NoError(err, "Adding edge (3,1) to g should not fail")

		// Create identical undirected graph h
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to h should not fail")
		err = h.AddEdgeWithOptions(3, 1)
		is.NoError(err, "Adding edge (3,1) to h should not fail")

		// Check if g isEqual h
		isEqual, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.True(isEqual, "Identical graphs g and h should satisfy Equals")
	})

	t.Run("Different Traits", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create undirected graph g with vertices {1, 2} and edge (1,2)
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create directed graph h with vertices {1, 2} and edge (1,2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.ErrorIs(err, graph.ErrGraphTypeMismatch, "Equals should return ErrGraphTypeMismatch due to different traits")
		is.False(equals, "Graphs with different traits should not be equal")
	})

	t.Run("Different Vertices", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create undirected graph g with vertices {1, 2, 3}
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to g should not fail")

		// Create undirected graph h with vertices {1, 2, 4}
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.False(equals, "Graphs g and h should not be equal due to different vertices")
	})

	t.Run("Different Edges", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create undirected graph g with vertices {1, 2, 3} and edges (1,2), (2,3)
		g, _ := simple.New(graph.IntHash) // undirected by default
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

		// Create undirected graph h with vertices {1, 2, 3} and edges (1,2), (3,1)
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")
		err = h.AddEdgeWithOptions(3, 1)
		is.NoError(err, "Adding edge (3,1) to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.False(equals, "Graphs g and h should not be equal due to different edges")
	})

	t.Run("Both Empty Graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create empty directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())

		// Create empty directed graph h
		h, _ := simple.New(graph.IntHash, graph.Directed())

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.True(equals, "Both empty graphs g and h should be equal")
	})

	t.Run("One Empty, One Non-Empty", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create empty directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())

		// Create directed graph h with vertices and edges
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err := h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddEdgeWithOptions(1, 1) // Self-loop
		is.NoError(err, "Adding edge (1,1) to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.False(equals, "Empty graph g should not equal non-empty graph h")
	})

	t.Run("Same Vertices, Different Edge Directions (Directed Graphs)", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g with vertices {1, 2} and edge (1,2)
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create directed graph h with vertices {1, 2} and edge (2,1)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(2, 1)
		is.NoError(err, "Adding edge (2,1) to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.False(equals, "Graphs g and h should not be equal due to different edge directions")
	})

	t.Run("Same Vertices, Different Edge Directions (undirected Graphs)", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create undirected graph g with vertices {1, 2} and edge (1,2)
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create undirected graph h with vertices {1, 2} and edge (2,1)
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(2, 1)
		is.NoError(err, "Adding edge (2,1) to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.True(equals, "undirected graphs g and h should be equal as edge direction is irrelevant")
	})

	t.Run("Graphs with Self-Loops", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g with self-loop on vertex 1
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddEdgeWithOptions(1, 1)
		is.NoError(err, "Adding self-loop (1,1) to g should not fail")

		// Create directed graph h with self-loop on vertex 1
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddEdgeWithOptions(1, 1)
		is.NoError(err, "Adding self-loop (1,1) to h should not fail")

		// Check if g equals h
		equals, err := Equals(g, h)
		is.NoError(err, "Equals should not return an error")
		is.True(equals, "Graphs g and h should be equal as they have identical self-loops")
	})
}
