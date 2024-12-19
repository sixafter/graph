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

// TestIsSubset tests the IsSubset function across various scenarios.
func TestIsSubset(t *testing.T) {
	t.Parallel()

	t.Run("Basic Subset", func(t *testing.T) {
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

		// Create directed graph h with vertices {1, 2, 3} and edges (1,2), (2,3)
		h, _ := simple.New(graph.IntHash, graph.Directed())
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.True(subset, "Interface g should be a subset of graph h")
	})

	t.Run("Not a Subset - Missing Vertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create directed graph g with vertices {1, 2, 4} and edge (1,2)
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create directed graph h with vertices {1, 2, 3} and edges (1,2), (2,3)
		h, _ := simple.New(graph.IntHash, graph.Directed())
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.False(subset, "Interface g should not be a subset of graph h due to missing vertex 4")
	})

	t.Run("Not a Subset - Missing Edge", func(t *testing.T) {
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.False(subset, "Interface g should not be a subset of graph h due to missing edge (2,3)")
	})

	t.Run("Empty Interface as Subset", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create an empty directed graph g
		g, _ := simple.New(graph.IntHash, graph.Directed())

		// Create any directed graph h with vertices and edges
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err := h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if empty g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.True(subset, "Empty graph g should be a subset of any graph h")
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.True(subset, "Interface g should be a subset of graph h as they are identical")
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.ErrorIs(err, graph.ErrGraphTypeMismatch, "IsSubset should return ErrGraphTypeMismatch due to different traits")
		is.False(subset, "Interface g should not be a subset of graph h due to different traits")
	})

	t.Run("Directed Graphs Subset", func(t *testing.T) {
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

		// Create directed graph h with vertices {1, 2, 3} and edges (1,2), (2,3)
		h, _ := simple.New(graph.IntHash, graph.Directed())
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.True(subset, "Directed graph g should be a subset of graph h")
	})

	t.Run("undirected Graphs Subset", func(t *testing.T) {
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

		// Create undirected graph h with vertices {1, 2, 3} and edges (1,2), (2,3)
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.True(subset, "undirected graph g should be a subset of graph h")
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

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.False(subset, "Interface g should not be a subset of graph h as there are no common vertices")
	})

	t.Run("g is not a subset - Extra Edge in h", func(t *testing.T) {
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

		// Create directed graph h with vertices {1, 2} and edges (1,2), (2,1)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")
		err = h.AddEdgeWithOptions(2, 1)
		is.NoError(err, "Adding edge (2,1) to h should not fail")

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.True(subset, "Interface g should be a subset of graph h even if h has extra edges")
	})

	t.Run("g has an edge not in h", func(t *testing.T) {
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

		// Create directed graph h with vertices {1, 2, 3} and only edge (1,2)
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to h should not fail")
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to h should not fail")

		// Check if g is a subset of h
		subset, err := IsSubset(g, h)
		is.NoError(err, "IsSubset should not return an error")
		is.False(subset, "Interface g should not be a subset of graph h because h is missing edge (2,3)")
	})
}
