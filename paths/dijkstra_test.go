// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package paths

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestDijkstraFrom(t *testing.T) {
	t.Parallel()

	t.Run("Finds shortest path in weighted graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := simple.New(graph.IntHash, graph.Directed(), graph.Weighted())
		is.NoError(g.AddVertexWithOptions(1))
		is.NoError(g.AddVertexWithOptions(2))
		is.NoError(g.AddVertexWithOptions(3))
		is.NoError(g.AddEdgeWithOptions(1, 2, simple.EdgeWeight(2)))
		is.NoError(g.AddEdgeWithOptions(2, 3, simple.EdgeWeight(1)))

		p, err := DijkstraFrom(g, 1, 3)
		is.NoError(err)
		is.Equal([]int{1, 2, 3}, p, "Shortest path from 1 to 3 should be [1, 2, 3]")
	})

	t.Run("Returns error for unreachable target", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := simple.New(graph.IntHash, graph.Directed())
		is.NoError(g.AddVertexWithOptions(1))
		is.NoError(g.AddVertexWithOptions(2))

		p, err := DijkstraFrom(g, 1, 2)
		is.Error(err)
		is.Nil(p, "Path should be nil if target is unreachable")
		is.ErrorIs(err, graph.ErrTargetNotReachable, "Error should be ErrTargetNotReachable")
	})
}
