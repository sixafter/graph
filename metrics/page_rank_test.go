// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"
	"math"
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func floatApproxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestPageRankDirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 1->3, 2->3, 3->1, 3->4
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Compute PageRank
	dampingFactor := 0.85
	maxIterations := 100
	tolerance := 1e-6
	pr, err := PageRank(g, dampingFactor, maxIterations, tolerance)
	is.NoError(err)

	// Retrieve PageRank scores
	c1 := pr[1]
	c2 := pr[2]
	c3 := pr[3]
	c4 := pr[4]

	// Assertions
	epsilon := 1e-4
	is.True(c3 > c1, "Vertex 3 should have a higher PageRank than Vertex 1")
	is.True(c3 > c4, "Vertex 3 should have a higher PageRank than Vertex 4")
	is.True(math.Abs(c1-c4) < epsilon, "Vertex 1 and Vertex 4 should have approximately equal PageRank")
	is.True(c1 > c2, "Vertex 1 should have a higher PageRank than Vertex 2")
	is.True(c4 > c2, "Vertex 4 should have a higher PageRank than Vertex 2")

	// Optional: Check total PageRank sum
	totalPR := c1 + c2 + c3 + c4
	is.True(floatApproxEqual(totalPR, 1.0, 1e-4), "Total PageRank should sum to approximately 1.0")
}

func TestPageRankWithDanglingNodes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 2->3, 3->1
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))
	// Vertex 4 is a dangling node (no outgoing edges)

	// Compute PageRank
	dampingFactor := 0.85
	maxIterations := 100
	tolerance := 1e-6
	pr, err := PageRank(g, dampingFactor, maxIterations, tolerance)
	is.NoError(err)

	// Retrieve PageRank scores
	c1 := pr[1]
	c2 := pr[2]
	c3 := pr[3]
	c4 := pr[4]

	epsilon := 1e-4
	// Assert that Vertex 1, Vertex 2, and Vertex 3 have approximately equal PageRank
	is.True(floatApproxEqual(c1, c2, epsilon), "Vertex 1 and Vertex 2 should have approximately equal PageRank")
	is.True(floatApproxEqual(c2, c3, epsilon), "Vertex 2 and Vertex 3 should have approximately equal PageRank")

	// Assert that Vertex 4 has a significantly lower PageRank than Vertex 1
	is.True(c4 < c1, "Vertex 4 should have a lower PageRank than Vertex 1")

	// Assert that Vertex 4 has some contribution due to teleportation (optional)
	is.True(c4 > 0.0, "Vertex 4 should have a positive PageRank due to teleportation")

	// Assert the total PageRank sums to approximately 1.0
	totalPR := c1 + c2 + c3 + c4
	is.True(floatApproxEqual(totalPR, 1.0, epsilon), "Total PageRank should sum to approximately 1.0")
}

func TestPageRankFullyConnectedDirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 2->1, 2->3, 3->2, 1->4, 4->3
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 1))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 2))
	is.NoError(g.AddEdgeWithOptions(1, 4))
	is.NoError(g.AddEdgeWithOptions(4, 3))

	// Compute PageRank
	dampingFactor := 0.85
	maxIterations := 100
	tolerance := 1e-6
	pr, err := PageRank(g, dampingFactor, maxIterations, tolerance)
	is.NoError(err)

	// Retrieve PageRank scores
	c1 := pr[1]
	c2 := pr[2]
	c3 := pr[3]
	c4 := pr[4]

	// Assertions
	is.True(c2 > c3, "Vertex 2 should have a higher PageRank than Vertex 3")
	is.True(c3 > c1, "Vertex 3 should have a higher PageRank than Vertex 1")
	is.True(c1 > c4, "Vertex 1 should have a higher PageRank than Vertex 4")

	// Optional: Check total PageRank sum
	totalPR := c1 + c2 + c3 + c4
	is.True(floatApproxEqual(totalPR, 1.0, 1e-4), "Total PageRank should sum to approximately 1.0")
}

func TestPageRankGraphWithSelfLoops(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges: 1->2, 2->3, 3->1, 3->4
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	// Compute PageRank
	dampingFactor := 0.85
	maxIterations := 100
	tolerance := 1e-6
	pr, err := PageRank(g, dampingFactor, maxIterations, tolerance)
	is.NoError(err)

	// Retrieve PageRank scores
	c1 := pr[1]
	c2 := pr[2]
	c3 := pr[3]
	c4 := pr[4]

	// Assertions
	epsilon := 1e-4

	// Assert that Vertex 2 has a higher PageRank than Vertex 1
	is.True(c2 > c1, "Vertex 2 should have a higher PageRank than Vertex 1")

	// Assert that Vertex 3 has the highest PageRank
	is.True(c3 > c2, "Vertex 3 should have a higher PageRank than Vertex 2")
	is.True(c3 > c1, "Vertex 3 should have a higher PageRank than Vertex 1")
	is.True(c3 > c4, "Vertex 3 should have a higher PageRank than Vertex 4")

	// Assert that Vertex 4 and Vertex 1 have approximately equal PageRank
	is.True(floatApproxEqual(c4, c1, epsilon), "Vertex 4 and Vertex 1 should have approximately equal PageRank")

	// Assert that the total PageRank sums to approximately 1.0
	totalPR := c1 + c2 + c3 + c4
	is.True(floatApproxEqual(totalPR, 1.0, epsilon), "Total PageRank should sum to approximately 1.0")

}

func TestPageRankNoEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph with no edges
	g, err := simple.New[int, int](graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Compute PageRank
	dampingFactor := 0.85
	maxIterations := 100
	tolerance := 1e-6
	pr, err := PageRank(g, dampingFactor, maxIterations, tolerance)
	is.NoError(err)

	// Expected centrality: All vertices have equal PageRank
	expectedPR := 1.0 / 3.0
	for i := 1; i <= 3; i++ {
		is.True(floatApproxEqual(pr[i], expectedPR, 1e-4), fmt.Sprintf("Vertex %d should have a PageRank of %.4f", i, expectedPR))
	}

	// Optional: Check total PageRank sum
	totalPR := pr[1] + pr[2] + pr[3]
	is.True(floatApproxEqual(totalPR, 1.0, 1e-4), "Total PageRank should sum to approximately 1.0")
}
