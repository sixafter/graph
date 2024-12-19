// Copyright (c) 2024 Six After, Inc
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

func TestFindAllPaths(t *testing.T) {
	t.Parallel()

	t.Run("Finds all paths in directedGraph graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := simple.New(graph.IntHash, graph.Directed())
		is.NoError(g.AddVertexWithOptions(1))
		is.NoError(g.AddVertexWithOptions(2))
		is.NoError(g.AddVertexWithOptions(3))
		is.NoError(g.AddVertexWithOptions(4))
		is.NoError(g.AddEdgeWithOptions(1, 2))
		is.NoError(g.AddEdgeWithOptions(1, 3))
		is.NoError(g.AddEdgeWithOptions(2, 4))
		is.NoError(g.AddEdgeWithOptions(3, 4))

		paths, err := FindAllPaths(g, 1, 4)
		is.NoError(err)
		is.Equal(2, len(paths), "There should be two paths from 1 to 4")
		is.ElementsMatch([][]int{{1, 2, 4}, {1, 3, 4}}, paths, "Paths should match expected")
	})

	t.Run("No paths if start or end not in graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := simple.New(graph.IntHash, graph.Directed())
		is.NoError(g.AddVertexWithOptions(1))

		paths, err := FindAllPaths(g, 1, 2)
		is.NoError(err)
		is.Equal(0, len(paths), "There should be no paths if end vertex is missing")
	})
}
