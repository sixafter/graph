// Copyright (c) 2024-2025 Six After, Inc
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

func TestBFSConnectedGraph(t *testing.T) {
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
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 4))
	is.NoError(g.AddEdgeWithOptions(3, 5))

	// Test BFS traversal
	visited := map[int]bool{}
	err := BFS(g, 1, func(vertex int) bool {
		visited[vertex] = true
		return false
	})

	is.NoError(err)

	// Verify all expected vertices are visited
	expectedVertices := []int{1, 2, 3, 4, 5}
	for _, v := range expectedVertices {
		is.True(visited[v], "Vertex %d should have been visited", v)
	}

	is.Equal(len(expectedVertices), len(visited), "Unexpected number of visited vertices")
}

func TestBFSDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices (disconnected)
	vertices := []int{1, 2, 3, 4}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges to form two disconnected components
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Test BFS traversal starting from vertex 1
	visited := map[int]bool{}
	err := BFS(g, 1, func(vertex int) bool {
		visited[vertex] = true
		return false
	})

	is.NoError(err)
	is.ElementsMatch([]int{1, 2}, visitedKeys(visited), "BFS should only visit the connected component of vertex 1")
}

func TestBFSWithCycle(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices
	vertices := []int{1, 2, 3, 4}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges to form a cycle
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 4))
	is.NoError(g.AddEdgeWithOptions(4, 1)) // Cycle back to 1

	// Test BFS traversal
	visited := map[int]bool{}
	err := BFS(g, 1, func(vertex int) bool {
		visited[vertex] = true
		return false
	})

	is.NoError(err)
	is.ElementsMatch([]int{1, 2, 3, 4}, visitedKeys(visited), "BFS should visit all vertices in a cyclic graph")
}

func TestBFSWithEarlyTermination(t *testing.T) {
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

	// Test BFS with early termination
	var visited []int
	err := BFS(g, 1, func(vertex int) bool {
		visited = append(visited, vertex)
		return vertex == 3 // Stop traversal when vertex 3 is visited
	})

	is.NoError(err)
	is.Equal([]int{1, 2, 3}, visited, "BFS should stop traversal at vertex 3 after processing all neighbors of vertex 1")
}
func TestBFSWithDepthTrackingDisconnectedGraph(t *testing.T) {
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
	is.NoError(g.AddEdgeWithOptions(4, 5)) // Disconnected component

	// Capture visited vertices with depth
	visited := make(map[int]int)
	err := BFSWithDepthTracking(g, 1, func(vertex, depth int) bool {
		visited[vertex] = depth
		return false
	})

	is.NoError(err)

	// Verify visited vertices are only in the connected component with correct depths
	expectedDepths := map[int]int{
		1: 0, // Start vertex at depth 0
		2: 1, // Neighbor at depth 1
		3: 1, // Neighbor at depth 1
	}
	is.Equal(expectedDepths, visited, "BFSWithDepthTracking should only visit the connected component of vertex 1 with correct depths")
}

func TestBFSWithNoEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add isolated vertices
	vertices := []int{1, 2, 3}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Perform BFS from one vertex
	visited := map[int]bool{}
	err := BFS(g, 1, func(vertex int) bool {
		visited[vertex] = true
		return false
	})

	is.NoError(err)
	is.ElementsMatch([]int{1}, visitedKeys(visited), "BFS should only visit the starting vertex when no edges exist")
}

// Utility function to extract keys from a map[int]bool
func visitedKeys(visited map[int]bool) []int {
	keys := make([]int, 0, len(visited))
	for k := range visited {
		keys = append(keys, k)
	}
	return keys
}
