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

func TestStronglyConnectedComponents(t *testing.T) {
	t.Parallel()

	t.Run("Finds SCCs in directedGraph graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := simple.New(graph.IntHash, graph.Directed())
		is.NoError(g.AddVertexWithOptions(1))
		is.NoError(g.AddVertexWithOptions(2))
		is.NoError(g.AddVertexWithOptions(3))
		is.NoError(g.AddEdgeWithOptions(1, 2))
		is.NoError(g.AddEdgeWithOptions(2, 3))
		is.NoError(g.AddEdgeWithOptions(3, 1))

		components, err := TarjanFrom(g)
		is.NoError(err)
		is.Equal(1, len(components), "Interface should have 1 SCC")
		is.ElementsMatch([]int{1, 2, 3}, components[0], "SCC should contain all vertices")
	})

	t.Run("Returns error for undirected graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := simple.New(graph.IntHash) // undirected graph
		is.NoError(g.AddVertexWithOptions(1))

		components, err := TarjanFrom(g)
		is.Error(err)
		is.Nil(components, "Components should be nil for undirected graph")
		is.ErrorIs(err, graph.ErrSCCDetectionNotDirected, "Error should be ErrSCCDetectionNotDirected")
	})
}
