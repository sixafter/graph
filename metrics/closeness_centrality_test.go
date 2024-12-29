// File: metrics/closeness_centrality_test.go

package metrics

import (
	"fmt"
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func TestClosenessCentralityEmptyGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an empty graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// Expect empty centrality map
	is.Empty(centrality, "Centrality map should be empty for an empty graph")
}

func TestClosenessCentralitySingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a graph with a single vertex
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// Closeness centrality for a single vertex should be 0.0
	is.True(floatEquals(0.0, centrality[1]), "Single vertex should have a closeness centrality of 0.0000")
}

func TestClosenessCentralityCompleteGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected complete graph with 4 vertices
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges between every pair of vertices
	for i := 1; i <= 4; i++ {
		for j := i + 1; j <= 4; j++ {
			is.NoError(g.AddEdgeWithOptions(i, j))
		}
	}

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// In a complete graph with 4 vertices:
	// Each vertex has a distance of 1 to the other 3 vertices.
	// Sum of distances = 3
	// C(v) = 3 / 3 = 1.0
	for i := 1; i <= 4; i++ {
		is.True(floatEquals(1.0, centrality[i]), fmt.Sprintf("Vertex %d should have a closeness centrality of 1.0000", i))
	}
}

func TestClosenessCentralityStarGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected star graph with center vertex 1 and leaves 2, 3, 4
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges (1-2, 1-3, 1-4)
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(1, 4))

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// Closeness centrality:
	// - Center vertex 1:
	//   Sum of distances = 1 + 1 + 1 = 3
	//   C(v1) = 3 / 3 = 1.0
	// - Leaf vertices (2, 3, 4):
	//   Sum of distances = 1 + 2 + 2 = 5
	//   C(v2, v3, v4) = 3 / 5 = 0.6
	expected := map[int]float64{
		1: 1.0,
		2: 0.6,
		3: 0.6,
		4: 0.6,
	}

	for i := 1; i <= 4; i++ {
		is.True(floatEquals(expected[i], centrality[i]),
			fmt.Sprintf("Vertex %d should have a closeness centrality of %.4f, got %.4f", i, expected[i], centrality[i]))
	}
}

func TestClosenessCentralityDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph with two disconnected components
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	// Component 1: Vertices 1, 2, 3 forming a triangle
	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Component 2: Vertices 4, 5 forming a single edge
	for i := 4; i <= 5; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}
	is.NoError(g.AddEdgeWithOptions(4, 5))

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// Closeness centrality:
	// - Vertices 1, 2, 3: Sum of distances = 2 each
	//   C(v) = 2 / 2 = 1.0
	// - Vertices 4, 5: Sum of distances = 1 each
	//   C(v4, v5) = 1 / 1 = 1.0
	expected := map[int]float64{
		1: 1.0,
		2: 1.0,
		3: 1.0,
		4: 1.0,
		5: 1.0,
	}

	for i := 1; i <= 5; i++ {
		is.True(floatEquals(expected[i], centrality[i]),
			fmt.Sprintf("Vertex %d should have a closeness centrality of %.4f, got %.4f", i, expected[i], centrality[i]))
	}
}

func TestClosenessCentralityDirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New(graph.IntHash, graph.Directed())
	is.NoError(err)

	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 1->3, 2->3, 3->1, 3->4
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// Expected closeness centrality based on standard formula
	expected := map[int]float64{
		1: 0.75, // 3 /4
		2: 0.6,  // 3 /5
		3: 0.75, // 3 /4
		4: 0.0,  // 0 /0
	}

	for i := 1; i <= 4; i++ {
		is.True(floatEquals(expected[i], centrality[i]),
			fmt.Sprintf("Vertex %d should have a closeness centrality of %.4f, got %.4f", i, expected[i], centrality[i]))
	}
}

func TestClosenessCentralityTriangle(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected triangle graph
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges to form a triangle
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	// Compute closeness centrality
	centrality, err := ClosenessCentrality(g)
	is.NoError(err)

	// In a triangle:
	// Each vertex has distances 1 to the other 2 vertices.
	// Sum of distances = 2
	// C(v) = 2 / 2 = 1.0
	expected := map[int]float64{
		1: 1.0,
		2: 1.0,
		3: 1.0,
	}

	for i := 1; i <= 3; i++ {
		is.True(floatEquals(expected[i], centrality[i]), fmt.Sprintf("Vertex %d should have a closeness centrality of %.4f", i, expected[i]))
	}
}
