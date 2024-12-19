// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package traverse

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestDFS(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices
	vertices := []int{1, 2, 3, 4, 5}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges to form a connected graph
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 4))
	is.NoError(g.AddEdgeWithOptions(4, 5))

	// Test DFS traversal
	var visited []int
	err := DFS(g, 1, func(vertex int) bool {
		visited = append(visited, vertex)
		return false // Continue traversal
	})

	is.NoError(err)
	is.Equal([]int{1, 2, 3, 4, 5}, visited)

	// Test early termination
	visited = []int{}
	err = DFS(g, 1, func(vertex int) bool {
		visited = append(visited, vertex)
		return vertex == 3 // Stop at vertex 3
	})

	is.NoError(err)
	is.Equal([]int{1, 2, 3}, visited)

	// Test starting from a disconnected vertex
	is.NoError(g.AddVertexWithOptions(6)) // Disconnected vertex
	visited = []int{}
	err = DFS(g, 6, func(vertex int) bool {
		visited = append(visited, vertex)
		return false
	})

	is.NoError(err)
	is.Equal([]int{6}, visited)
}
