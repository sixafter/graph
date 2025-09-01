// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"testing"

	"github.com/sixafter/graph"
	internal "github.com/sixafter/graph/internal/paths"
	"github.com/stretchr/testify/assert"
)

func TestCreatesCycle(t *testing.T) {
	t.Parallel()

	t.Run("Detects a cycle in a directedGraph graph", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.IntHash, graph.Directed(), graph.PreventCycles())
		is.NoError(g.AddVertexWithOptions(1))
		is.NoError(g.AddVertexWithOptions(2))
		is.NoError(g.AddVertexWithOptions(3))
		is.NoError(g.AddEdgeWithOptions(1, 2))
		is.NoError(g.AddEdgeWithOptions(2, 3))

		// Adding edge 3 -> 1 should create a cycle
		isCycle, err := internal.WouldCreateCycle(g, 3, 1)
		is.NoError(err)
		is.True(isCycle, "Adding edge 3 -> 1 should create a cycle")
	})

	t.Run("No cycle when adding edge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		g, _ := New(graph.IntHash, graph.Directed(), graph.PreventCycles())
		is.NoError(g.AddVertexWithOptions(1))
		is.NoError(g.AddVertexWithOptions(2))
		is.NoError(g.AddVertexWithOptions(3)) // Add missing vertex 3
		is.NoError(g.AddEdgeWithOptions(1, 2))

		// Adding edge 2 -> 3 should not create a cycle
		isCycle, err := internal.WouldCreateCycle(g, 2, 3)
		is.NoError(err)
		is.False(isCycle, "Adding edge 2 -> 3 should not create a cycle")
	})
}

func TestDirectedAcyclicGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directedGraph acyclic graph
	g, _ := New[string, string](graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())

	// Add vertices
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddVertexWithOptions("C"))
	is.NoError(g.AddVertexWithOptions("D"))
	is.NoError(g.AddVertexWithOptions("E"))

	// Add edges to form a DAG
	is.NoError(g.AddEdgeWithOptions("A", "B"))
	is.NoError(g.AddEdgeWithOptions("A", "C"))
	is.NoError(g.AddEdgeWithOptions("B", "D"))
	is.NoError(g.AddEdgeWithOptions("C", "D"))
	is.NoError(g.AddEdgeWithOptions("D", "E"))

	// Attempt to add an edge that creates a cycle
	err := g.AddEdgeWithOptions("E", "A")
	is.EqualError(err, graph.ErrEdgeCreatesCycle.Error()) // Ensure the error is returned

	// Verify adjacency map
	adjacency, _ := g.AdjacencyMap()
	expectedAdjacency := map[string]map[string]graph.Edge[string]{
		"A": {

			"B": NewEdgeWithOptions("A", "B", EdgeWeight(0)),
			"C": NewEdgeWithOptions("A", "C", EdgeWeight(0)),
		},
		"B": {
			"D": NewEdgeWithOptions("B", "D", EdgeWeight(0)),
		},
		"C": {
			"D": NewEdgeWithOptions("C", "D", EdgeWeight(0)),
		},
		"D": {
			"E": NewEdgeWithOptions("D", "E", EdgeWeight(0)),
		},
		"E": {}, // No outgoing edges
	}
	is.Equal(expectedAdjacency, adjacency)
}

// TestWeightedGraph tests adding edges with weights in a directedGraph graph.
func TestWeightedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed(), graph.Weighted())

	// Add vertices
	vertices := []string{"A", "B", "C", "D", "E", "F"}
	for _, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add weighted edges
	is.NoError(g.AddEdgeWithOptions("A", "B", EdgeWeight(5)))
	is.NoError(g.AddEdgeWithOptions("A", "C", EdgeWeight(3)))
	is.NoError(g.AddEdgeWithOptions("B", "D", EdgeWeight(2)))
	is.NoError(g.AddEdgeWithOptions("C", "D", EdgeWeight(7)))
	is.NoError(g.AddEdgeWithOptions("D", "E", EdgeWeight(1)))
	is.NoError(g.AddEdgeWithOptions("E", "F", EdgeWeight(4)))

	// Verify edge weights
	edgeAB, _ := g.Edge("A", "B")
	is.Equal(float64(5), edgeAB.Properties().Weight())

	_, err := g.Edge("C", "F")
	is.Error(err) // Edge should not exist
}

// TestRootedGraph tests a rooted graph structure.
func TestRootedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed(), graph.Rooted())

	// Add vertices
	is.NoError(g.AddVertexWithOptions("Root"))
	is.NoError(g.AddVertexWithOptions("Child1"))
	is.NoError(g.AddVertexWithOptions("Child2"))
	is.NoError(g.AddVertexWithOptions("Grandchild1"))
	is.NoError(g.AddVertexWithOptions("Grandchild2"))

	// Add edges
	is.NoError(g.AddEdgeWithOptions("Root", "Child1"))
	is.NoError(g.AddEdgeWithOptions("Root", "Child2"))
	is.NoError(g.AddEdgeWithOptions("Child1", "Grandchild1"))
	is.NoError(g.AddEdgeWithOptions("Child2", "Grandchild2"))

	// Verify root node has no predecessors
	predecessors, _ := g.PredecessorMap()
	is.Empty(predecessors["Root"])

	// Verify child nodes have Root as predecessor
	is.Contains(predecessors["Child1"], "Root")
	is.Contains(predecessors["Child2"], "Root")
}

// TestPreventCyclesInAcyclicGraph ensures cycles cannot be created when PreventCycles is enabled.
func TestPreventCyclesInAcyclicGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed(), graph.Acyclic(), graph.PreventCycles())

	// Add vertices
	for _, v := range []string{"A", "B", "C", "D", "E", "F"} {
		is.NoError(g.AddVertexWithOptions(v))
	}

	// Add edges to form an acyclic graph
	is.NoError(g.AddEdgeWithOptions("A", "B"))
	is.NoError(g.AddEdgeWithOptions("B", "C"))
	is.NoError(g.AddEdgeWithOptions("C", "D"))
	is.NoError(g.AddEdgeWithOptions("D", "E"))
	is.NoError(g.AddEdgeWithOptions("E", "F"))

	// Attempt to add an edge that creates a cycle
	err := g.AddEdgeWithOptions("F", "A")
	is.EqualError(err, graph.ErrEdgeCreatesCycle.Error())
}

// TestGraphWithAttributes tests vertices and edges with attributes.
func TestGraphWithAttributes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed())

	// Add vertices with attributes
	is.NoError(g.AddVertexWithOptions("A", VertexItem("color", "red")))
	is.NoError(g.AddVertexWithOptions("B", VertexItem("color", "blue")))
	is.NoError(g.AddVertexWithOptions("C"))

	// Add edges with attributes
	is.NoError(g.AddEdgeWithOptions("A", "B", EdgeItem("weight", 10)))
	is.NoError(g.AddEdgeWithOptions("B", "C", EdgeItem("weight", 20)))

	// Retrieve and verify vertex attributes
	vp, _ := g.Vertex("A")
	is.Equal("red", vp.Properties().Items()["color"])

	// Retrieve and verify edge attributes
	edgeAB, _ := g.Edge("A", "B")
	is.Equal(10, edgeAB.Properties().Items()["weight"])
}

func TestGraphOperations(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[int, int](graph.IntHash)

	// Add vertices
	for i := 1; i <= 20; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges to create a connected graph
	for i := 1; i < 20; i++ {
		is.NoError(g.AddEdgeWithOptions(i, i+1)) // Adds (1,2), (2,3), ..., (19,20)
	}

	// Explicitly add edge (5, 7) for testing
	is.NoError(g.AddEdgeWithOptions(5, 7))

	// Verify graph order and size
	order, _ := g.Order()
	size, _ := g.Size()
	is.Equal(20, order)
	is.Equal(20, size) // One extra edge added: (5,7)

	// Remove some edges
	is.NoError(g.RemoveEdge(5, 6))   // Remove an existing edge
	is.NoError(g.RemoveEdge(10, 11)) // Remove an existing edge
	is.NoError(g.RemoveEdge(5, 7))   // Now valid because we added this edge earlier

	// Verify size after edge removal
	sizeAfterRemoval, _ := g.Size()
	is.Equal(17, sizeAfterRemoval)

	// Remove a vertex
	// First, explicitly remove all edges connected to vertex 5
	is.NoError(g.RemoveEdge(4, 5))
	is.NoError(g.RemoveVertex(5)) // Vertex 5 is fully disconnected

	// Verify order after vertex removal
	orderAfterRemoval, _ := g.Order()
	is.Equal(19, orderAfterRemoval)

	// Verify remaining edges do not include vertex 5
	for _, vertex := range []int{4, 6, 7} {
		_, err := g.Edge(5, vertex)
		is.Error(err)
		_, err = g.Edge(vertex, 5)
		is.Error(err)
	}
}

// TestCycleDetectionInDirectedGraph ensures that cycles can exist in a directedGraph graph when allowed.
func TestCycleDetectionInDirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed())

	// Add vertices
	for _, v := range []string{"A", "B", "C", "D"} {
		i := NewVertexWithOptions(v, v)
		is.NoError(g.AddVertex(i))
	}

	// Add edges to form a cycle: A -> B -> C -> A
	is.NoError(g.AddEdgeWithOptions("A", "B"))
	is.NoError(g.AddEdgeWithOptions("B", "C"))
	is.NoError(g.AddEdgeWithOptions("C", "A"))

	// Since the graph is not acyclic and PreventCycles is not enabled, this is allowed
}

// TestGraphCloning tests cloning a graph and ensuring the clone is a deep copy.
func TestGraphCloning(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Directed())

	// Add vertices and edges
	is.NoError(g.AddVertexWithOptions("A"))
	is.NoError(g.AddVertexWithOptions("B"))
	is.NoError(g.AddEdgeWithOptions("A", "B", EdgeWeight(5)))

	// Clone the graph
	clone, err := g.Clone()
	is.NoError(err)

	// Verify the clone has the same vertices and edges
	edge, _ := clone.Edge("A", "B")
	is.Equal(float64(5), edge.Properties().Weight())

	// Modify the original graph
	is.NoError(g.RemoveEdge("A", "B"))

	// Verify the clone is unaffected
	edgeInClone, err := clone.Edge("A", "B")
	is.NoError(err)
	is.Equal(float64(5), edgeInClone.Properties().Weight())
}

// TestWeightedUndirectedGraph tests an undirected graph with weighted edges.
func TestWeightedUndirectedGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, _ := New[string, string](graph.StringHash, graph.Weighted())

	// Add vertices
	vertices := []string{"Node1", "Node2", "Node3", "Node4"}
	for k, v := range vertices {
		is.NoError(g.AddVertexWithOptions(v, VertexMetadata(k)))
	}

	// Add weighted edges
	is.NoError(g.AddEdgeWithOptions("Node1", "Node2", EdgeWeight(3)))
	is.NoError(g.AddEdgeWithOptions("Node2", "Node3", EdgeWeight(5)))
	is.NoError(g.AddEdgeWithOptions("Node3", "Node4", EdgeWeight(7)))
	is.NoError(g.AddEdgeWithOptions("Node4", "Node1", EdgeWeight(2)))

	// Verify edge weights in both directions
	edge12, _ := g.Edge("Node1", "Node2")
	edge21, _ := g.Edge("Node2", "Node1")
	is.Equal(float64(3), edge12.Properties().Weight())
	is.Equal(float64(3), edge21.Properties().Weight())
}

func TestEdgeOptions(t *testing.T) {
	t.Parallel()

	t.Run("EdgeWeight sets the weight of an edge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		props := &EdgeProperties{}
		option := EdgeWeight(42)
		option(props)

		is.Equal(float64(42), props.Weight(), "Edge weight should be set to 42")
	})

	t.Run("EdgeItems sets multiple attributes on an edge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		props := &EdgeProperties{}
		attributes := map[string]any{"color": "blue", "width": 5}
		option := EdgeItems(attributes)
		option(props)

		is.Equal(attributes, props.Items(), "Edge attributes should match the input map")

		option = EdgeItem("color", "red")
		option(props)

		is.Equal("red", props.Items()["color"], "Edge attribute 'color' should be 'red'")

	})

	t.Run("EdgeData sets custom data for an edge", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		props := &EdgeProperties{}
		data := "custom edge data"
		option := EdgeData(data)
		option(props)

		is.Equal(data, props.Metadata(), "Edge data should be set to 'custom edge data'")
	})
}

func TestVertexOptions(t *testing.T) {
	t.Parallel()

	t.Run("VertexWeight sets the weight of a vertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		props := &VertexProperties{}
		option := VertexWeight(10)
		option(props)

		is.Equal(float64(10), props.Weight(), "Vertex weight should be set to 10")
	})

	t.Run("VertexItems sets multiple attributes on a vertex", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)

		props := &VertexProperties{}
		attributes := map[string]any{"color": "green", "size": 20}
		option := VertexItems(attributes)
		option(props)

		is.Equal(attributes, props.Items(), "Vertex attributes should match the input map")

		option = VertexItem("name", "nodeA")
		option(props)

		is.Equal("nodeA", props.Items()["name"], "Vertex attribute 'name' should be 'nodeA'")

	})
}

func TestDefaultCursor(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cursor := EmptyCursor()
	is.NotNil(cursor, "EmptyCursor should return a non-nil cursor")

	state := cursor.State()
	is.NotNil(state, "EmptyCursor state should be nil")

	is.NoError(cursor.SetState(state), "Setting state on an empty cursor should not return an error")
}
