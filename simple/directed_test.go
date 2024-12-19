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

func TestDirectedGraph(t *testing.T) {
	t.Parallel()

	t.Run("AddVertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A", VertexWeight(5), VertexItem("label", "VertexA"))
		is.NoError(err, "Adding vertex should not fail")

		vertex, err := g.Vertex("A")
		is.NoError(err, "Fetching vertex should not fail")
		is.Equal("A", vertex.Value(), "Vertex should match the added value")
	})

	t.Run("AddVerticesFrom", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		source, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := source.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = source.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")

		target, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err = target.AddVerticesFrom(source)
		is.NoError(err, "Adding vertices from another graph should not fail")

		order, err := target.Order()
		is.NoError(err, "Fetching graph order should not fail")
		is.Equal(2, order, "Interface should contain the same number of vertices as the source")
	})

	t.Run("AddEdge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A")
		err = g.AddVertexWithOptions("B")

		e := NewEdgeWithOptions("A", "B", EdgeWeight(3))
		err = g.AddEdge(e)
		is.NoError(err, "Adding edge should not fail")

		edge, err := g.Edge("A", "B")
		is.NoError(err, "Fetching edge should not fail")
		is.Equal("A", edge.Source(), "Edge source should match")
		is.Equal("B", edge.Target(), "Edge target should match")
		is.Equal(float64(3), edge.Properties().Weight(), "Edge weight should match")
	})

	t.Run("AddEdgesFrom", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		source, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := source.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = source.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")
		e := NewEdgeWithOptions("A", "B")
		err = source.AddEdge(e)
		is.NoError(err, "Adding edge A -> B should not fail")

		target, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err = target.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = target.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")

		err = target.AddEdgesFrom(source)
		is.NoError(err, "Adding edges from another graph should not fail")

		size, err := target.Size()
		is.NoError(err, "Fetching graph size should not fail")
		is.Equal(1, size, "Interface should contain the same number of edges as the source")
	})

	t.Run("RemoveVertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex should not fail")

		err = g.RemoveVertex("A")
		is.NoError(err, "Removing vertex should not fail")

		_, err = g.Vertex("A")
		is.Error(err, "Fetching removed vertex should fail")
	})

	t.Run("RemoveEdge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = g.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")
		e := NewEdgeWithOptions("A", "B")
		err = g.AddEdge(e)
		is.NoError(err, "Adding edge A -> B should not fail")

		err = g.RemoveEdge("A", "B")
		is.NoError(err, "Removing edge should not fail")

		_, err = g.Edge("A", "B")
		is.Error(err, "Fetching removed edge should fail")
	})

	t.Run("AdjacencyMap", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = g.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")
		e := NewEdgeWithOptions("A", "B")
		err = g.AddEdge(e)
		is.NoError(err, "Adding edge A -> B should not fail")

		adjMap, err := g.AdjacencyMap()
		is.NoError(err, "Fetching adjacency map should not fail")
		is.Contains(adjMap["A"], "B", "Adjacency map should contain the edge A -> B")
	})

	t.Run("PredecessorMap", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = g.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")
		e := NewEdgeWithOptions("A", "B")
		err = g.AddEdge(e)
		is.NoError(err, "Adding edge A -> B should not fail")

		predMap, err := g.PredecessorMap()
		is.NoError(err, "Fetching predecessor map should not fail")
		is.Contains(predMap["B"], "A", "Predecessor map should contain the edge A -> B")
	})

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		err := g.AddVertexWithOptions("A")
		is.NoError(err, "Adding vertex A should not fail")
		err = g.AddVertexWithOptions("B")
		is.NoError(err, "Adding vertex B should not fail")
		e := NewEdgeWithOptions("A", "B")
		err = g.AddEdge(e)
		is.NoError(err, "Adding edge A -> B should not fail")

		clone, err := g.Clone()
		is.NoError(err, "Cloning graph should not fail")

		adjMap, err := clone.AdjacencyMap()
		is.NoError(err, "Fetching adjacency map of clone should not fail")
		is.Contains(adjMap["A"], "B", "Cloned graph should contain the edge A -> B")
	})

	t.Run("HasEdge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		// Create a directed graph
		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())

		// Add vertices and edges
		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))
		is.NoError(g.AddVertexWithOptions("C"))
		e := NewEdgeWithOptions("A", "B")
		is.NoError(g.AddEdge(e)) // A -> B
		e = NewEdgeWithOptions("B", "C")
		is.NoError(g.AddEdge(e)) // B -> C

		// Check for existing edges
		hasEdge, err := g.HasEdge("A", "B")
		is.NoError(err)
		is.True(hasEdge, "Edge A -> B should exist")

		hasEdge, err = g.HasEdge("B", "C")
		is.NoError(err)
		is.True(hasEdge, "Edge B -> C should exist")

		// Check for non-existent edges
		hasEdge, err = g.HasEdge("B", "A")
		is.NoError(err)
		is.False(hasEdge, "Edge B -> A should not exist")

		hasEdge, err = g.HasEdge("C", "A")
		is.NoError(err)
		is.False(hasEdge, "Edge C -> A should not exist")

		// Check for non-existent vertices
		hasEdge, err = g.HasEdge("D", "E")
		is.NoError(err)
		is.False(hasEdge, "Edge D -> E should not exist")
	})

	t.Run("HasVertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New[string, string](graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())

		// Add some vertices
		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))

		exists, err := g.HasVertex("A")
		is.NoError(err)
		is.True(exists, "Vertex A should exist in the graph")

		exists, err = g.HasVertex("C")
		is.NoError(err)
		is.False(exists, "Vertex C should not exist in the graph")
	})

	t.Run("Degree", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())

		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))
		is.NoError(g.AddVertexWithOptions("C"))
		e := NewEdgeWithOptions("A", "B")
		is.NoError(g.AddEdge(e)) // A -> B
		e = NewEdgeWithOptions("C", "A")
		is.NoError(g.AddEdge(e)) // C -> A

		degree, err := g.Degree("A")
		is.NoError(err, "Degree should not fail")
		is.Equal(2, degree, "Degree of A should be 2 (1 in + 1 out)")
	})

	t.Run("InDegree", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))
		is.NoError(g.AddVertexWithOptions("C"))
		e := NewEdgeWithOptions("B", "A")
		is.NoError(g.AddEdge(e)) // B -> A
		e = NewEdgeWithOptions("C", "A")
		is.NoError(g.AddEdge(e)) // C -> A

		inDegree, err := g.InDegree("A")
		is.NoError(err, "InDegree should not fail")
		is.Equal(2, inDegree, "InDegree of A should be 2")
	})

	t.Run("OutDegree", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))
		is.NoError(g.AddVertexWithOptions("C"))
		e := NewEdgeWithOptions("A", "B")
		is.NoError(g.AddEdge(e)) // A -> B
		e = NewEdgeWithOptions("A", "C")
		is.NoError(g.AddEdge(e)) // A -> C

		outDegree, err := g.OutDegree("A")
		is.NoError(err, "OutDegree should not fail")
		is.Equal(2, outDegree, "OutDegree of A should be 2")
	})

	t.Run("Neighbors", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))
		is.NoError(g.AddVertexWithOptions("C"))
		e := NewEdgeWithOptions("A", "B")
		is.NoError(g.AddEdge(e)) // A -> B
		e = NewEdgeWithOptions("A", "C")
		is.NoError(g.AddEdge(e)) // A -> C

		neighbors, err := g.Neighbors("A")
		is.NoError(err, "Neighbors should not fail")
		is.ElementsMatch([]graph.Vertex[string, string]{
			NewVertexWithOptions(graph.StringHash("B"), "B"),
			NewVertexWithOptions(graph.StringHash("C"), "C"),
		}, neighbors, "Neighbors of A should be B and C")
	})

	t.Run("Vertices", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())
		is.NoError(g.AddVertexWithOptions("A"))
		is.NoError(g.AddVertexWithOptions("B"))
		is.NoError(g.AddVertexWithOptions("C"))

		vertices, err := g.Vertices()
		is.NoError(err, "Vertices should not fail")
		is.ElementsMatch([]graph.Vertex[string, string]{
			NewVertexWithOptions(graph.StringHash("A"), "A"),
			NewVertexWithOptions(graph.StringHash("B"), "B"),
			NewVertexWithOptions(graph.StringHash("C"), "C"),
		}, vertices, "Vertices should include A, B, and C")
	})
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
