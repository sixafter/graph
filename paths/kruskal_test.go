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

func TestUnionFind(t *testing.T) {
	t.Parallel()

	t.Run("Initialize and find roots", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		uf := newUnionFind(1, 2, 3, 4)
		is.Equal(1, uf.Find(1), "Initial root of 1 should be itself")
		is.Equal(2, uf.Find(2), "Initial root of 2 should be itself")
	})

	t.Run("Add new vertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		uf := newUnionFind(1, 2)
		uf.Add(3)
		is.Equal(3, uf.Find(3), "Newly added vertex should be its own root")
	})

	t.Run("Union merges sets", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		uf := newUnionFind(1, 2)
		uf.Union(1, 2)
		is.Equal(uf.Find(1), uf.Find(2), "1 and 2 should have the same root after union")
	})

	t.Run("Find with path compression", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		uf := newUnionFind(1, 2, 3)
		uf.Union(1, 2)
		uf.Union(2, 3)

		root := uf.Find(3)
		is.Equal(uf.Find(1), root, "All vertices should have the same root after union")
	})
}

type testEdge struct {
	Source, Target string
	Weight         float64
}

// verifyMST checks that the given MST spans all vertices in the graph and
// includes the correct total weight based on expected edges.
func verifyMST(t *testing.T, g graph.Interface[string, string], mst graph.Interface[string, string], expectedEdges []testEdge) {
	is := assert.New(t)

	// Verify that all vertices from the original graph are in the MST
	vertices, _ := g.Order()
	mstVertices, _ := mst.Order()
	is.Equal(vertices, mstVertices, "MST should contain all vertices")

	// Verify the edges in the MST
	mstEdges, _ := mst.Edges()
	is.Equal(len(expectedEdges), len(mstEdges), "MST should contain the correct number of edges")

	// Verify total weight and edge structure
	expectedWeight := float64(0)
	for _, edge := range expectedEdges {
		expectedWeight += edge.Weight
		found := false
		for _, mstEdge := range mstEdges {
			if (mstEdge.Source() == edge.Source && mstEdge.Target() == edge.Target) ||
				(mstEdge.Source() == edge.Target && mstEdge.Target() == edge.Source) {
				is.Equal(edge.Weight, mstEdge.Properties().Weight(), "Edge weight mismatch")
				found = true
				break
			}
		}
		is.True(found, "Edge %v -> %v not found in MST", edge.Source, edge.Target)
	}
}

func TestMinimumSpanningTree(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a weighted graph
	g, _ := simple.New[string, string](graph.StringHash, graph.Weighted())
	vertices := []string{"A", "B", "C", "D", "E"}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}
	edges := []testEdge{
		{"A", "B", 1},
		{"A", "C", 4},
		{"B", "C", 2},
		{"B", "D", 6},
		{"C", "D", 3},
		{"C", "E", 5},
		{"D", "E", 7},
	}
	for _, edge := range edges {
		is.NoError(g.AddEdgeWithOptions(edge.Source, edge.Target, simple.EdgeWeight(float64(edge.Weight))))
	}

	// Calculate the MST
	mst, err := MinimumSpanningTree(g)
	is.NoError(err)

	// Verify the MST
	expectedEdges := []testEdge{
		{"A", "B", 1},
		{"B", "C", 2},
		{"C", "D", 3},
		{"C", "E", 5},
	}
	verifyMST(t, g, mst, expectedEdges)
}

func TestMaximumSpanningTree(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a weighted graph
	g, _ := simple.New[string, string](graph.StringHash, graph.Weighted())
	vertices := []string{"A", "B", "C", "D", "E"}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}
	edges := []testEdge{
		{"A", "B", 1},
		{"A", "C", 4},
		{"B", "C", 2},
		{"B", "D", 6},
		{"C", "D", 3},
		{"C", "E", 5},
		{"D", "E", 7},
	}
	for _, edge := range edges {
		is.NoError(g.AddEdgeWithOptions(edge.Source, edge.Target, simple.EdgeWeight(float64(edge.Weight))))
	}

	// Calculate the MST
	mst, err := MaximumSpanningTree(g)
	is.NoError(err)

	// Verify the MST
	expectedEdges := []testEdge{
		{"D", "E", 7},
		{"B", "D", 6},
		{"C", "E", 5},
		{"A", "C", 4},
	}
	verifyMST(t, g, mst, expectedEdges)
}

func TestSpanningTreeNoEdges(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a graph with no edges
	g, _ := simple.New[string, string](graph.StringHash, graph.Weighted())
	vertices := []string{"A", "B", "C"}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Calculate the MST
	mst, err := MinimumSpanningTree(g)
	is.NoError(err)

	// Verify the MST
	mstEdges, _ := mst.Edges()
	is.Equal(0, len(mstEdges), "MST should have no edges for a graph with no edges")
}

func TestSpanningTreeSingleVertex(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a graph with a single vertex
	g, _ := simple.New[string, string](graph.StringHash, graph.Weighted())
	is.NoError(g.AddVertexWithOptions("A"))

	// Calculate the MST
	mst, err := MinimumSpanningTree(g)
	is.NoError(err)

	// Verify the MST
	mstEdges, _ := mst.Edges()
	is.Equal(0, len(mstEdges), "MST should have no edges for a single-vertex graph")
}

func TestSpanningTreeDisconnectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a disconnected graph
	g, _ := simple.New[string, string](graph.StringHash, graph.Weighted())
	vertices := []string{"A", "B", "C", "D"}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}
	is.NoError(g.AddEdgeWithOptions("A", "B", simple.EdgeWeight(3)))
	is.NoError(g.AddEdgeWithOptions("C", "D", simple.EdgeWeight(5)))

	// Calculate the MST
	mst, err := MinimumSpanningTree(g)
	is.NoError(err)

	// Verify the MST
	mstEdges, _ := mst.Edges()
	is.LessOrEqual(len(mstEdges), len(vertices)-1, "MST should span connected components")
}
