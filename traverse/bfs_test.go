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

func TestBFS(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices
	vertices := []int{1, 2, 3}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges to form a connected graph
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))

	// Test BFS traversal
	visited := map[int]bool{}
	err := BFS(g, 1, func(vertex int) bool {
		visited[vertex] = true
		return false // Continue traversal
	})

	is.NoError(err)

	// Verify all expected vertices are visited
	expectedVertices := []int{1, 2, 3}
	for _, v := range expectedVertices {
		is.True(visited[v], "Vertex %d should have been visited", v)
	}

	// Verify the number of visited vertices
	is.Equal(len(expectedVertices), len(visited), "Unexpected number of visited vertices")
}

func TestBFSWithDepth(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices
	vertices := []int{1, 2, 3, 4, 5}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 4))
	is.NoError(g.AddEdgeWithOptions(3, 5))

	// Expected depths for each vertex
	expectedDepths := map[int]int{
		1: 1,
		2: 2,
		3: 2,
		4: 3,
		5: 3,
	}

	// Capture visited vertices with depth
	visited := make(map[int]int)
	err := BFSWithDepth(g, 1, func(vertex, depth int) bool {
		visited[vertex] = depth
		return false // Continue traversal
	})

	is.NoError(err)

	// Verify that all expected vertices were visited
	is.Equal(len(expectedDepths), len(visited), "Number of visited vertices mismatch")

	// Verify depths
	for vertex, expectedDepth := range expectedDepths {
		actualDepth, ok := visited[vertex]
		is.True(ok, "Vertex %d was not visited", vertex)
		is.Equal(expectedDepth, actualDepth, "Depth mismatch for vertex %d", vertex)
	}
}
