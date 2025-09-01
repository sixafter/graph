// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"context"
	"sort"
	"testing"

	"github.com/sixafter/graph"
	"github.com/stretchr/testify/assert"
)

func TestAddVertex_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash) // Default to undirected
	err := g.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA"))
	is.NoError(err, "Adding vertex should not fail")

	vertex, err := g.Vertex("A")
	is.NoError(err, "Fetching vertex should not fail")
	is.Equal("A", vertex.ID(), "Vertex should match the added value")
}

func TestAddVerticesFrom_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	source, _ := New[string, string](graph.StringHash)
	is.NoError(source.AddVertexWithOptions("A"))
	is.NoError(source.AddVertexWithOptions("B"))

	target, _ := New[string, string](graph.StringHash)
	err := target.AddVerticesFrom(source)
	is.NoError(err, "Adding vertices from another graph should not fail")

	order, err := target.Order()
	is.NoError(err, "Fetching graph order should not fail")
	is.Equal(2, order, "Interface should contain the same number of vertices as the source")
}

func TestAddEdge_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))

	err := g.AddEdgeWithOptions("A", "B", EdgeWeight(3))
	is.NoError(err, "Adding edge should not fail")

	edge, err := g.Edge("A", "B")
	is.NoError(err, "Fetching edge should not fail")
	is.Equal("A", edge.Source(), "Edge source should match")
	is.Equal("B", edge.Target(), "Edge target should match")
	is.Equal(float64(3), edge.Properties().Weight(), "Edge weight should match")

	reverseEdge, err := g.Edge("B", "A")
	is.NoError(err, "Fetching reverse edge in undirected graph should not fail")
	is.Equal("A", reverseEdge.Target(), "Reverse edge target should match original source")
	is.Equal("B", reverseEdge.Source(), "Reverse edge source should match original target")
}

func TestAddEdgesFrom_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	source, _ := New[string, string](graph.StringHash)
	is.NoError(source.AddVertexWithOptions("A"))
	is.NoError(source.AddVertexWithOptions("B"))
	is.NoError(source.AddEdgeWithOptions("A", "B"))

	target, _ := New[string, string](graph.StringHash)
	is.NoError(target.AddVertexWithOptions("A"))
	is.NoError(target.AddVertexWithOptions("B"))

	err := target.AddEdgesFrom(source)
	is.NoError(err, "Adding edges from another graph should not fail")

	size, err := target.Size()
	is.NoError(err, "Fetching graph size should not fail")
	is.Equal(1, size, "Interface should contain the same number of unique edges as the source")
}

func TestRemoveVertex_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A"))

	err := g.RemoveVertex("A")
	is.NoError(err, "Removing vertex should not fail")

	_, err = g.Vertex("A")
	is.Error(err, "Fetching removed vertex should fail")
}

func TestRemoveEdge_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash) // undirected by default
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B"))

	err := g.RemoveEdge("A", "B")
	is.NoError(err, "Removing edge should not fail")

	hasEdge, err := g.HasEdge("A", "B")
	is.NoError(err, "Checking edge existence should not fail")
	is.False(hasEdge, "Edge A -- B should no longer exist")
}

func TestAdjacencyMap_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B"))

	adjMap, err := g.AdjacencyMap()
	is.NoError(err, "Fetching adjacency map should not fail")
	is.Contains(adjMap["A"], "B", "Adjacency map should contain the edge A -> B")
	is.Contains(adjMap["B"], "A", "Adjacency map should contain the edge B -> A for undirected graph")
}

func TestPredecessorMap_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B"))

	predMap, err := g.PredecessorMap()
	is.NoError(err, "Fetching predecessor map should not fail")
	is.Contains(predMap["B"], "A", "Predecessor map should contain the edge A -> B")
	is.Contains(predMap["A"], "B", "Predecessor map should contain the edge B -> A for undirected graph")
}

func TestCloneGraph_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B"))

	clone, err := g.Clone()
	is.NoError(err, "Cloning graph should not fail")

	adjMap, err := clone.AdjacencyMap()
	is.NoError(err, "Fetching adjacency map of clone should not fail")
	is.Contains(adjMap["A"], "B", "Cloned graph should contain the edge A -> B")
	is.Contains(adjMap["B"], "A", "Cloned graph should contain the reverse edge B -> A for undirected graph")
}

func TestHasEdges_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)

	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddVertexWithOptions("C"))

	is.NoError(g.AddEdgeWithOptions("A", "B"))
	is.NoError(g.AddEdgeWithOptions("B", "C"))

	hasEdge, err := g.HasEdge("A", "B")
	is.NoError(err)
	is.True(hasEdge, "Edge A -- B should exist")
}

func TestHasVertex_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)

	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))

	exists, err := g.HasVertex("A")
	is.NoError(err)
	is.True(exists, "Vertex A should exist in the graph")
}

func TestStreamEdgesWithContext_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash) // Create an undirected graph
	is.NoError(g.AddVertexWithOptions("A", VertexWeight(5)))
	is.NoError(g.AddVertexWithOptions("B", VertexWeight(5)))
	is.NoError(g.AddVertexWithOptions("C", VertexWeight(5)))
	is.NoError(g.AddEdgeWithOptions("A", "B", EdgeWeight(1)))
	is.NoError(g.AddEdgeWithOptions("B", "C", EdgeWeight(2)))

	ctx := context.Background()
	cursor := &Cursor{position: 0}
	ch := make(chan graph.Edge[string])

	go func() {
		_, err := g.StreamEdgesWithContext(ctx, cursor, 2, ch)
		is.NoError(err, "Streaming edges should not fail")
	}()

	var result []graph.Edge[string]
	for edge := range ch {
		result = append(result, edge)
	}

	// Extract and sort sources for comparison
	var sources []string
	for _, edge := range result {
		sources = append(sources, edge.Source())
	}
	sort.Strings(sources)

	is.ElementsMatch([]string{"A", "B"}, sources, "Streamed edges should have correct sources")
}

func TestStreamVerticesWithContext_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash) // Create an undirected graph
	is.NoError(g.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA")))
	is.NoError(g.AddVertexWithOptions("B", VertexWeight(10), VertexItem("label", "VertexB")))

	ctx := context.Background()
	cursor := &Cursor{position: 0}
	ch := make(chan []graph.Vertex[string, string])

	go func() {
		_, err := g.StreamVerticesWithContext(ctx, cursor, 1, ch)
		is.NoError(err, "Streaming vertices should not fail")
	}()

	var result [][]graph.Vertex[string, string]
	for batch := range ch {
		result = append(result, batch)
	}

	var ids []string
	for _, batch := range result {
		for _, vertex := range batch {
			ids = append(ids, vertex.ID())
		}
	}
	sort.Strings(ids)
	is.ElementsMatch([]string{"A", "B"}, ids, "Streamed vertices should have correct IDs")
}

func TestUndirected_SetVertexWithOptions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA")))

	is.NoError(g.SetVertexWithOptions("A", VertexWeight(10), VertexItem("label", "UpdatedVertexA")))
	vertex, err := g.Vertex("A")
	is.NoError(err)

	is.Equal(float64(10), vertex.Properties().Weight(), "Vertex weight should be updated")
	is.Equal("UpdatedVertexA", vertex.Properties().Items()["label"], "Vertex label should be updated")
}

func TestUndirected_SetEdgeWithOptions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash)
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B", EdgeWeight(5)))

	is.NoError(g.SetEdgeWithOptions("A", "B", EdgeWeight(10)))
	edge, err := g.Edge("A", "B")
	is.NoError(err)

	is.Equal(float64(10), edge.Properties().Weight(), "Edge weight should be updated")
}

func TestNeighbors_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, _ := New(graph.StringHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddVertexWithOptions("C"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "C")))

	// Fetch neighbors of "A"
	neighbors, err := g.Neighbors("A")
	is.NoError(err, "Fetching neighbors should not fail")

	// Extract IDs from the returned vertices
	var neighborIDs []string
	for _, neighbor := range neighbors {
		neighborIDs = append(neighborIDs, neighbor.ID())
	}

	// Assert the neighbor IDs
	is.ElementsMatch([]string{"B", "C"}, neighborIDs, "Neighbors should match")
}

func TestDegree_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph
	g, _ := New(graph.StringHash) // Default is undirected
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddVertexWithOptions("C"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "C")))

	// Test degree of vertex "A" (2 edges connected)
	degree, err := g.Degree("A")
	is.NoError(err, "Fetching degree should not fail")
	is.Equal(2, degree, "Degree of vertex A should be 2")

	// Test degree of vertex "B" (1 edge connected)
	degree, err = g.Degree("B")
	is.NoError(err, "Fetching degree should not fail")
	is.Equal(1, degree, "Degree of vertex B should be 1")
}

func TestOutDegree_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph
	g, _ := New(graph.StringHash) // Default is undirected
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddVertexWithOptions("C"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "C")))

	// Test out-degree of vertex "A" (treated as degree in undirected graph)
	outDegree, err := g.OutDegree("A")
	is.NoError(err, "Fetching out-degree should not fail for undirected graph")
	is.Equal(2, outDegree, "Out-degree of vertex A should be 2 in undirected graph")
}

func TestInDegree_Undirected(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected graph
	g, _ := New(graph.StringHash) // Default is undirected
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddVertexWithOptions("C"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "C")))

	// Test in-degree of vertex "A" (treated as degree in undirected graph)
	inDegree, err := g.InDegree("A")
	is.NoError(err, "Fetching in-degree should not fail for undirected graph")
	is.Equal(2, inDegree, "In-degree of vertex A should be 2 in undirected graph")
}
