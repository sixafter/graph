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

// TestIsSuperset tests the IsSuperset function across various scenarios.
func TestIsSuperset(t *testing.T) {
	t.Parallel()

	t.Run("Basic Superset", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g with vertices {1, 2, 3} and edges (1,2), (2,3)
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

		// Create directed graph h with vertices {1, 2} and edge (1,2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "Interface g should be a superset of graph h")
	})

	t.Run("Not a Superset - Missing Vertex", func(t *testing.T) {
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

		// Create directed graph h with vertices {1, 2, 3} and edge (1,2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.False(superset, "Interface g should not be a superset of graph h due to missing vertex 3")
	})

	t.Run("Not a Superset - Missing Edge", func(t *testing.T) {
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

		// Create directed graph h with vertices {1, 2} but no edges
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "Interface g should be a superset of graph h as h has no edges")
	})

	t.Run("Empty Interface h", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create any directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create an empty directed graph h
		h, _ := simple.New(graph.IntHash, graph.Directed())

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "Any graph g should be a superset of an empty graph h")
	})

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

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "Identical graphs g and h should satisfy IsSuperset")
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

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.ErrorIs(err, graph.ErrGraphTypeMismatch, "IsSuperset should return ErrGraphTypeMismatch due to different traits")
		is.False(superset, "Interface g should not be a superset of graph h due to different traits")
	})

	t.Run("Directed Graphs Superset", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g with vertices {1, 2, 3} and edges (1,2), (2,3)
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

		// Create directed graph h with vertices {1, 2} and edge (1,2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "Interface g should be a superset of graph h")
	})

	t.Run("undirected Graphs Superset", func(t *testing.T) {
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

		// Create undirected graph h with vertices {1, 2} and edge (1,2)
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "undirected graph g should be a superset of graph h")
	})

	t.Run("No Common Vertices", func(t *testing.T) {
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

		// Create directed graph h with vertices {3, 4} and edge (3,4)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to h should not fail")
		err = h.AddEdgeWithOptions(3, 4)
		is.NoError(err, "Adding edge (3,4) to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.False(superset, "Interface g should not be a superset of graph h as there are no common vertices")
	})

	t.Run("Superset with Extra Vertices and Edges", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g with vertices {1, 2, 3, 4} and edges (1,2), (2,3), (3,4)
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to g should not fail")
		err = g.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")
		err = g.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to g should not fail")
		err = g.AddEdgeWithOptions(3, 4)
		is.NoError(err, "Adding edge (3,4) to g should not fail")

		// Create directed graph h with vertices {1, 2} and edge (1,2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g is a superset of h
		superset, err := IsSuperset(g, h)
		is.NoError(err, "IsSuperset should not return an error")
		is.True(superset, "Interface g should be a superset of graph h")
	})
}
