package traverse

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestDFSConnectedGraph(t *testing.T) {
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
	is.Equal([]int{1, 2, 3, 4, 5}, visited, "DFS should visit all vertices in order for a connected graph")
}

func TestDFSEarlyTermination(t *testing.T) {
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

	// Test early termination
	var visited []int
	err := DFS(g, 1, func(vertex int) bool {
		visited = append(visited, vertex)
		return vertex == 3 // Stop at vertex 3
	})

	is.NoError(err)
	is.Equal([]int{1, 2, 3}, visited, "DFS should stop traversal when early termination condition is met")
}

func TestDFSDisconnectedVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices, including a disconnected vertex
	vertices := []int{1, 2, 3, 4, 5, 6} // Vertex 6 is disconnected
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges to form a connected subgraph
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 4))
	is.NoError(g.AddEdgeWithOptions(4, 5))

	// Start DFS from the disconnected vertex
	var visited []int
	err := DFS(g, 6, func(vertex int) bool {
		visited = append(visited, vertex)
		return false
	})

	is.NoError(err)
	is.Equal([]int{6}, visited, "DFS starting from a disconnected vertex should only visit that vertex")
}

func TestDFSUnreachableVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := simple.New[int, int](graph.IntHash)

	// Add vertices and edges
	vertices := []int{1, 2, 3, 4, 5, 6, 7}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Create two disconnected subgraphs
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(4, 5))
	is.NoError(g.AddEdgeWithOptions(5, 6))

	// Test DFS starting from a vertex in the first subgraph
	var visited []int
	err := DFS(g, 1, func(vertex int) bool {
		visited = append(visited, vertex)
		return false
	})

	is.NoError(err)
	is.Equal([]int{1, 2, 3}, visited, "DFS should only traverse the connected component of the starting vertex")
}

func TestDFSWithCycle(t *testing.T) {
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

	// Perform DFS traversal
	var visited []int
	err := DFS(g, 1, func(vertex int) bool {
		visited = append(visited, vertex)
		return false
	})

	is.NoError(err)
	is.ElementsMatch([]int{1, 2, 3, 4}, visited, "DFS should visit all vertices even in the presence of a cycle")
}
