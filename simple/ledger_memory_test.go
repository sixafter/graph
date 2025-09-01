// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"testing"

	"github.com/sixafter/graph"
	"github.com/stretchr/testify/assert"
)

func TestVertexOperations(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s, err := newMemoryStore[int, string]()
	is.NoError(err)

	v := NewVertexWithOptions(1, "A", VertexItems(map[string]any{
		"label":   "vertexA",
		"active":  true,
		"count":   10,
		"decimal": 3.14,
	}), VertexWeight(5))

	// Add vertices with various properties
	is.NoError(s.AddVertex(v.ID(), v.Value(), v.Properties()))

	v = NewVertexWithOptions(2, "B", VertexItems(map[string]any{
		"label":    "vertexB",
		"priority": 1,
		"active":   false,
	}), VertexWeight(10))

	is.NoError(s.AddVertex(v.ID(), v.Value(), v.Properties()))

	// Attempt to add a duplicate vertex
	err = s.AddVertex(1, "DuplicateA", &VertexProperties{})
	is.ErrorIs(err, graph.ErrVertexAlreadyExists)

	// Find vertices and verify properties
	vertex, props, err := s.FindVertex(1)
	is.NoError(err)
	is.Equal("A", vertex)
	is.Equal(float64(5), props.Weight())
	is.Equal("vertexA", props.Items()["label"])
	is.Equal(true, props.Items()["active"])
	is.Equal(10, props.Items()["count"])
	is.Equal(3.14, props.Items()["decimal"])

	// Modify vertex properties
	v = NewVertexWithOptions(1, "A", VertexItems(map[string]any{
		"label": "updatedA",
		"new":   true,
	}), VertexWeight(20))
	is.NoError(s.ModifyVertex(1, v.Properties()))

	_, updatedProps, err := s.FindVertex(1)
	is.NoError(err)
	is.Equal(float64(20), updatedProps.Weight())
	is.Equal("updatedA", updatedProps.Items()["label"])
	is.Equal(true, updatedProps.Items()["new"])

	// List vertices
	vertices, err := s.ListVertices()
	is.NoError(err)
	is.ElementsMatch([]int{1, 2}, vertices)

	// Delete a vertex
	is.NoError(s.RemoveVertex(1))

	_, _, err = s.FindVertex(1)
	is.ErrorIs(err, graph.ErrVertexNotFound)
}

func TestEdgeOperations(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s, err := newMemoryStore[int, string]()
	is.NoError(err)

	// Add vertices for edges
	is.NoError(s.AddVertex(1, "A", &VertexProperties{}))
	is.NoError(s.AddVertex(2, "B", &VertexProperties{}))
	is.NoError(s.AddVertex(3, "C", &VertexProperties{}))

	// Add edges with properties
	e := NewEdgeWithOptions(1, 2, EdgeItems(map[string]any{
		"relation":   "connects",
		"direction":  "one-way",
		"threshold":  7.89,
		"retryCount": 3,
	}), EdgeWeight(2))
	is.NoError(s.AddEdge(1, 2, e))

	e = NewEdgeWithOptions(2, 3, EdgeItems(map[string]any{
		"relation":      "extends",
		"bidirectional": true,
		"maxLoad":       100,
	}), EdgeWeight(4))
	is.NoError(s.AddEdge(2, 3, e))

	// Find edges and verify properties
	edge, err := s.FindEdge(1, 2)
	is.NoError(err)
	is.Equal(1, edge.Source())
	is.Equal(2, edge.Target())
	is.Equal(float64(2), edge.Properties().Weight())
	is.Equal("connects", edge.Properties().Items()["relation"])
	is.Equal("one-way", edge.Properties().Items()["direction"])
	is.Equal(7.89, edge.Properties().Items()["threshold"])
	is.Equal(3, edge.Properties().Items()["retryCount"])

	// Modify edge properties
	e = NewEdgeWithOptions(1, 2, EdgeItems(map[string]any{
		"relation": "updatedConnects",
		"urgent":   true,
	}), EdgeWeight(3.0))
	is.NoError(s.ModifyEdge(1, 2, e))

	updatedEdge, err := s.FindEdge(1, 2)
	is.NoError(err)
	is.Equal(float64(3), updatedEdge.Properties().Weight())
	is.Equal("updatedConnects", updatedEdge.Properties().Items()["relation"])
	is.Equal(true, updatedEdge.Properties().Items()["urgent"])

	// List all edges
	edges, err := s.ListEdges()
	is.NoError(err)
	is.Len(edges, 2)

	// Delete an edge
	is.NoError(s.RemoveEdge(1, 2))

	_, err = s.FindEdge(1, 2)
	is.ErrorIs(err, graph.ErrEdgeNotFound)
}

func TestCycleDetection(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s, err := newMemoryStore[int, string]()
	is.NoError(err)

	// Add vertices
	is.NoError(s.AddVertex(1, "A", &VertexProperties{}))
	is.NoError(s.AddVertex(2, "B", &VertexProperties{}))
	is.NoError(s.AddVertex(3, "C", &VertexProperties{}))

	// Add edges
	is.NoError(s.AddEdge(1, 2, NewEdgeWithOptions(1, 2)))
	is.NoError(s.AddEdge(2, 3, NewEdgeWithOptions(2, 3)))

	mem := s.(*memoryLedger[int, string])

	// Test for no cycle
	hasCycle, err := mem.WouldCreateCycle(1, 3)
	is.NoError(err)
	is.False(hasCycle)

	// Test for cycle
	hasCycle, err = mem.WouldCreateCycle(3, 1)
	is.NoError(err)
	is.True(hasCycle)
}
