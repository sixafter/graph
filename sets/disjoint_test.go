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

func TestIsDisjoint(t *testing.T) {
	t.Parallel()

	t.Run("Disjoint undirected Graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create first undirected graph `g`
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create second undirected graph `h`
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to h should not fail")
		err = h.AddEdgeWithOptions(3, 4)
		is.NoError(err, "Adding edge (3,4) to h should not fail")

		// Test if `g` and `h` are disjoint
		disjoint, err := IsDisjoint(g, h)
		is.NoError(err, "IsDisjoint should not return an error")
		is.True(disjoint, "Graphs g and h should be disjoint")
	})

	t.Run("Non-Disjoint undirected Graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create first undirected graph `g`
		g, _ := simple.New(graph.IntHash) // undirected by default
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create second undirected graph `h` with overlapping vertex
		h, _ := simple.New(graph.IntHash) // undirected by default
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to h should not fail")

		// Test if `g` and `h` are disjoint
		disjoint, err := IsDisjoint(g, h)
		is.NoError(err, "IsDisjoint should not return an error")
		is.False(disjoint, "Graphs g and h should not be disjoint as they share vertex 2")
	})

	t.Run("Disjoint Directed Graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create first directed graph `g`
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create second directed graph `h`
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddVertexWithOptions(4)
		is.NoError(err, "Adding vertex 4 to h should not fail")
		err = h.AddEdgeWithOptions(3, 4)
		is.NoError(err, "Adding edge (3,4) to h should not fail")

		// Test if `g` and `h` are disjoint
		disjoint, err := IsDisjoint(g, h)
		is.NoError(err, "IsDisjoint should not return an error")
		is.True(disjoint, "Graphs g and h should be disjoint")
	})

	t.Run("Non-Disjoint Directed Graphs", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create first directed graph `g`
		g, _ := simple.New(graph.IntHash, graph.Directed())
		err := g.AddVertexWithOptions(1)
		is.NoError(err, "Adding vertex 1 to g should not fail")
		err = g.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to g should not fail")
		err = g.AddEdgeWithOptions(1, 2)
		is.NoError(err, "Adding edge (1,2) to g should not fail")

		// Create second directed graph `h` with overlapping vertex
		h, _ := simple.New(graph.IntHash, graph.Directed())
		err = h.AddVertexWithOptions(2)
		is.NoError(err, "Adding vertex 2 to h should not fail")
		err = h.AddVertexWithOptions(3)
		is.NoError(err, "Adding vertex 3 to h should not fail")
		err = h.AddEdgeWithOptions(2, 3)
		is.NoError(err, "Adding edge (2,3) to h should not fail")

		// Test if `g` and `h` are disjoint
		disjoint, err := IsDisjoint(g, h)
		is.NoError(err, "IsDisjoint should not return an error")
		is.False(disjoint, "Graphs g and h should not be disjoint as they share vertex 2")
	})
}
