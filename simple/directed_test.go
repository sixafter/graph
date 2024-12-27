// Copyright (c) 2024 Six After, Inc
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

func TestAddVertex_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	err := g.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA"))
	is.NoError(err, "Adding vertex should not fail")

	vertex, err := g.Vertex("A")
	is.NoError(err, "Fetching vertex should not fail")
	is.Equal("A", vertex.Value(), "Vertex should match the added value")
}

func TestAddVerticesFrom_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	source, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(source.AddVertexWithOptions("A"))
	is.NoError(source.AddVertexWithOptions("B"))

	target, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(target.AddVerticesFrom(source))

	order, err := target.Order()
	is.NoError(err, "Fetching graph order should not fail")
	is.Equal(2, order, "Target graph should contain the same number of vertices as the source")
}

func TestAddEdge_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))

	e := NewEdgeWithOptions("A", "B", EdgeWeight(3))
	is.NoError(g.AddEdge(e), "Adding edge should not fail")

	edge, err := g.Edge("A", "B")
	is.NoError(err, "Fetching edge should not fail")
	is.Equal("A", edge.Source(), "Edge source should match")
	is.Equal("B", edge.Target(), "Edge target should match")
	is.Equal(float64(3), edge.Properties().Weight(), "Edge weight should match")
}

func TestAddEdgesFrom_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	source, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(source.AddVertexWithOptions("A"))
	is.NoError(source.AddVertexWithOptions("B"))
	is.NoError(source.AddEdge(NewEdgeWithOptions("A", "B")))

	target, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(target.AddVertexWithOptions("A"))
	is.NoError(target.AddVertexWithOptions("B"))
	is.NoError(target.AddEdgesFrom(source))

	size, err := target.Size()
	is.NoError(err, "Fetching graph size should not fail")
	is.Equal(1, size, "Target graph should contain the same number of edges as the source")
}

func TestRemoveVertex_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.RemoveVertex("A"))

	_, err := g.Vertex("A")
	is.Error(err, "Fetching removed vertex should fail")
}

func TestRemoveEdge_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))
	is.NoError(g.RemoveEdge("A", "B"))

	_, err := g.Edge("A", "B")
	is.Error(err, "Fetching removed edge should fail")
}

func TestAdjacencyMap_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))

	adjMap, err := g.AdjacencyMap()
	is.NoError(err, "Fetching adjacency map should not fail")
	is.Contains(adjMap["A"], "B", "Adjacency map should contain the edge A -> B")
}

func TestPredecessorMap_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdge(NewEdgeWithOptions("A", "B")))

	predMap, err := g.PredecessorMap()
	is.NoError(err, "Fetching predecessor map should not fail")
	is.Contains(predMap["B"], "A", "Predecessor map should contain the edge A -> B")
}

func TestStreamEdgesWithContext_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	d, _ := New(graph.StringHash, graph.Directed()) // Create a directed graph
	is.NoError(d.AddVertexWithOptions("A", VertexWeight(5)))
	is.NoError(d.AddVertexWithOptions("B", VertexWeight(5)))
	is.NoError(d.AddVertexWithOptions("C", VertexWeight(5)))
	is.NoError(d.AddVertexWithOptions("D", VertexWeight(5)))
	is.NoError(d.AddEdgeWithOptions("A", "B", EdgeWeight(1)))
	is.NoError(d.AddEdgeWithOptions("C", "D", EdgeWeight(2)))

	ctx := context.Background()
	cursor := &Cursor{position: 0}
	ch := make(chan graph.Edge[string])

	go func() {
		_, err := d.StreamEdgesWithContext(ctx, cursor, 2, ch)
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

	is.ElementsMatch([]string{"A", "C"}, sources, "Streamed edges should have correct sources")
}

func TestStreamVerticesWithContext_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	d, _ := New(graph.StringHash, graph.Directed()) // Create a directed graph
	is.NoError(d.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA")))
	is.NoError(d.AddVertexWithOptions("B", VertexWeight(10), VertexItem("label", "VertexB")))

	ctx := context.Background()
	cursor := &Cursor{position: 0}
	ch := make(chan []graph.Vertex[string, string])

	go func() {
		_, err := d.StreamVerticesWithContext(ctx, cursor, 1, ch)
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

func TestSetVertexWithOptions_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA")))

	exists, err := g.HasVertex("A")
	is.NoError(err)
	is.True(exists)

	is.NoError(g.SetVertexWithOptions("A", VertexWeight(10), VertexItem("label", "UpdatedVertexA")))
	vertex, err := g.Vertex("A")
	is.NoError(err)

	is.Equal(float64(10), vertex.Properties().Weight(), "Vertex weight should be updated")
	is.Equal("UpdatedVertexA", vertex.Properties().Items()["label"], "Vertex label should be updated")
}

func TestSetEdgeWithOptions_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed())
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B", EdgeWeight(5)))

	is.NoError(g.SetEdgeWithOptions("A", "B", EdgeWeight(10)))
	edge, err := g.Edge("A", "B")
	is.NoError(err)

	is.Equal(float64(10), edge.Properties().Weight(), "Edge weight should be updated")
}

func TestNeighbors_Directed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed graph
	g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
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
