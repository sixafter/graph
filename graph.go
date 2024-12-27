// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

// Package graph provides interfaces and types to model and manipulate graphs,
// including vertices, edges, and their associated properties.
package graph

import (
	"context"

	"golang.org/x/exp/constraints"
)

type Ordered interface {
	comparable // Ensures equality (`==`, `!=`) is supported.
	constraints.Ordered
}

// Interface represents a generic graph metadata structure consisting of vertices of
// type T identified by a hash of type K. It provides methods for managing
// vertices, edges, and graph properties.
type Interface[K Ordered, T any] interface {
	// AddVertex creates a new Vertex in the graph. If the Vertex already exists,
	// it returns ErrVertexAlreadyExists. Functional options can be used to set
	// Vertex properties such as weight or v.
	//
	// Example:
	//  hash := graph.StringHash("A")
	//  vertex := graph.NewVertexWithOptions(hash, "A", graph.VertexWeight(4), graph.VertexItem("label", "Node A"))
	//	err := graph.AddVertex(vertex)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	AddVertex(value Vertex[K, T]) error

	// AddVertexWithOptions creates a new Vertex in the graph. If the Vertex already exists,
	// it returns ErrVertexAlreadyExists. Functional options can be used to set
	// Vertex properties such as weight or v.
	//
	// Example:
	//	err := graph.AddVertexWithOptions("A", graph.VertexWeight(4), graph.VertexItem("label", "Node A"))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	AddVertexWithOptions(value T, options ...VertexOption) error

	// AddVerticesFrom imports all vertices along with their properties from the
	// given graph into the current graph. Stops and returns an error if any Vertex
	// already exists.
	//
	// Example:
	//	err := graph.AddVerticesFrom(otherGraph)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	AddVerticesFrom(g Interface[K, T]) error

	// Vertex retrieves the Vertex with the given hash value. Returns
	// ErrVertexNotFound if the Vertex does not exist.
	//
	// Example:
	//	vertex, err := graph.Vertex("A")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(vertex)
	Vertex(hash K) (Vertex[K, T], error)

	// SetVertexWithOptions updates the properties of an existing vertex using functional
	// options. Returns ErrVertexNotFound if the vertex does not exist.
	//
	// Example:
	//	err := graph.SetVertexWithOptions("A", "B", graph.VertexWeight(20))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	SetVertexWithOptions(value T, options ...VertexOption) error

	// RemoveVertex removes the Vertex identified by the given hash from the graph.
	// Returns ErrVertexHasEdges if the Vertex still has edges, or ErrVertexNotFound
	// if the Vertex does not exist.
	//
	// Example:
	//	err := graph.RemoveVertex("A")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	RemoveVertex(hash K) error

	// Vertices returns a slice containing all Vertex values in the graph.
	//
	// Example:
	//	vertices, err := graph.Vertices()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	for _, vertex := range vertices {
	//		fmt.Println(vertex)
	//	}
	Vertices() ([]Vertex[K, T], error)

	// StreamVerticesWithContext streams vertices from the graph in paginated batches. The stream can be canceled
	// or timed out using a context, and the state of the cursor can be used to resume from the last point.
	//
	// Parameters:
	//   - ctx: The context to manage cancellation or timeout of the stream.
	//   - cursor: A Cursor object to track and manage the streaming state.
	//   - limit: The maximum number of vertices to include in each batch.
	//   - ch: The channel to which vertices will be sent.
	//
	// Returns:
	//   - An updated Cursor reflecting the new position in the stream.
	//   - An error if the operation is canceled or encounters an issue.
	//
	// Example:
	//   cursor := simple.EmptyCursor()
	//   ch := make(chan []Vertex[int, string])
	//   go func() {
	//       cursor, err := g.StreamVerticesWithContext(ctx, cursor, 10, ch)
	//       if err != nil {
	//           log.Printf("Stream failed: %v", err)
	//       }
	//   }()
	//   for batch := range ch {
	//       for _, vertex := range batch {
	//           fmt.Printf("Vertex: %v\n", vertex)
	//       }
	//   }
	StreamVerticesWithContext(ctx context.Context, cursor Cursor, limit int, ch chan<- []Vertex[K, T]) (Cursor, error)

	// HasVertex checks if a Vertex with the specified identifier exists in the graph.
	//
	// Parameters:
	//   - hash: The identifier of the Vertex to check. Its type (K) must be comparable.
	//
	// Returns:
	//   - bool: A boolean value indicating whether the Vertex exists in the graph (true if it exists, false otherwise).
	//   - error: An error if the operation fails (e.g., due to an underlying storage or implementation issue).
	//
	// Example:
	//   exists, err := graph.HasVertex("vertex1")
	//   if err != nil {
	//       log.Fatalf("Error checking Vertex: %v", err)
	//   }
	//   if exists {
	//       fmt.Println("Vertex exists in the graph.")
	//   } else {
	//       fmt.Println("Vertex does not exist in the graph.")
	//   }
	HasVertex(hash K) (bool, error)

	// AddEdge creates an edge between the source and target vertices. Returns
	// ErrVertexNotFound if either Vertex is missing, ErrEdgeAlreadyExists if the
	// edge already exists, or ErrEdgeCreatesCycle if adding the edge creates a
	// cycle in a cycle-preventing graph. Functional options can be used to set
	// edge properties such as weight or v.
	//
	// Example:
	//	edge := simple.NewEdgeWithOptions("A", "B", graph.EdgeWeight(10), graph.EdgeItem("color", "blue"))
	//	err := graph.AddEdge(edge)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	AddEdge(edge Edge[K]) error

	// AddEdgeWithOptions creates an edge between the source and target vertices. Returns
	// ErrVertexNotFound if either Vertex is missing, ErrEdgeAlreadyExists if the
	// edge already exists, or ErrEdgeCreatesCycle if adding the edge creates a
	// cycle in a cycle-preventing graph. Functional options can be used to set
	// edge properties such as weight or v.
	//
	// Example:
	//	err := graph.AddEdgeWithOptions("A", "B", graph.EdgeWeight(10), graph.EdgeItem("color", "blue"))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	AddEdgeWithOptions(source, target K, options ...EdgeOption) error

	// AddEdgesFrom imports all edges from another graph into the current graph.
	// The vertices that the edges connect must already exist in the current graph.
	//
	// Example:
	//	err := graph.AddEdgesFrom(otherGraph)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	AddEdgesFrom(g Interface[K, T]) error

	// Edge retrieves the edge between two vertices. For Undirected graphs, the
	// order of the source and target vertices does not matter. Returns ErrEdgeNotFound
	// if the edge does not exist.
	//
	// Example:
	//	edge, err := graph.Edge("A", "B")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(edge)
	Edge(source, target K) (Edge[T], error)

	// Edges returns a slice of all edges in the graph. Each edge Contains the
	// source and target Vertex hashes, along with edge properties.
	//
	// Example:
	//	edges, err := graph.Edges()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	for _, edge := range edges {
	//		fmt.Printf("Edge: %+v\n", edge)
	//	}
	Edges() ([]Edge[K], error)

	// StreamEdgesWithContext streams edges from the graph in paginated batches. The stream can be canceled
	// or timed out using a context, and the state of the cursor can be used to resume from the last point.
	//
	// Parameters:
	//   - ctx: The context to manage cancellation or timeout of the stream.
	//   - cursor: A Cursor object to track and manage the streaming state.
	//   - limit: The maximum number of edges to include in each batch.
	//   - ch: The channel to which edges will be sent.
	//
	// Returns:
	//   - An updated Cursor reflecting the new position in the stream.
	//   - An error if the operation is canceled or encounters an issue.
	//
	// Example:
	//   cursor := simple.EmptyCursor()
	//   ch := make(chan Edge[int])
	//   go func() {
	//       cursor, err := g.StreamEdgesWithContext(ctx, cursor, 10, ch)
	//       if err != nil {
	//           log.Printf("Stream failed: %v", err)
	//       }
	//   }()
	//   for edge := range ch {
	//       fmt.Printf("Edge: %v\n", edge)
	//   }
	StreamEdgesWithContext(ctx context.Context, cursor Cursor, limit int, ch chan<- Edge[K]) (Cursor, error)

	// SetEdgeWithOptions updates the properties of an existing edge using functional
	// options. Returns ErrEdgeNotFound if the edge does not exist.
	//
	// Example:
	//	err := graph.SetEdgeWithOptions("A", "B", graph.EdgeWeight(20))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	SetEdgeWithOptions(source, target K, options ...EdgeOption) error

	// RemoveEdge removes the edge between the given source and target vertices.
	// Returns ErrEdgeNotFound if the edge does not exist.
	//
	// Example:
	//	err := graph.RemoveEdge("A", "B")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	RemoveEdge(source, target K) error

	// HasEdge checks if an edge exists between the source and target vertices.
	// Returns true if the edge exists, false otherwise.
	//
	// Example:
	//	exists, err := graph.HasEdge("A", "B")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Edge exists: %v\n", exists)
	HasEdge(source, target K) (bool, error)

	// AdjacencyMap generates and returns an adjacency map representing all
	// outgoing edges for each Vertex in the graph.
	//
	// Example:
	//	adjMap, err := graph.AdjacencyMap()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Adjacency Map: %+v\n", adjMap)
	AdjacencyMap() (map[K]map[K]Edge[K], error)

	// PredecessorMap generates and returns a map representing all incoming
	// edges for each Vertex in the graph.
	//
	// Example:
	//	predMap, err := graph.PredecessorMap()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Predecessor Map: %+v\n", predMap)
	PredecessorMap() (map[K]map[K]Edge[K], error)

	// Clone creates a deep copy of the graph and returns the new instance.
	//
	// Example:
	//	clonedGraph, err := graph.Clone()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	Clone() (Interface[K, T], error)

	// Order returns the number of vertices in the graph.
	//
	// Example:
	//	order, err := graph.Order()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Interface order: %d\n", order)
	Order() (int, error)

	// Size returns the number of edges in the graph.
	//
	// Example:
	//	size, err := graph.Size()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Interface size: %d\n", size)
	Size() (int, error)

	// Hash retrieves the hash function used by the graph. The hash function maps a Vertex
	// of type T to a hash of type K, which is used internally to identify vertices.
	//
	// This allows clients to understand or reuse the hashing logic for operations
	// that involve Vertex identification, such as comparing vertices or exporting
	// the graph structure.
	//
	// Returns:
	//   - A Hash[K, T] function that maps vertices to their unique hash values.
	//
	// Example:
	//   hashFunc := graph.Hash()
	//   vertex := "A"
	//   vertexHash := hashFunc(Vertex)
	//   fmt.Printf("Hash for Vertex %v: %v\n", vertex, vertexHash)
	Hash() Hash[K, T]

	// Traits returns the graph's traits, such as whether it is directed, weighted,
	// or acyclic. These traits must be set when creating a graph using the New function.
	//
	// Example:
	//	traits := graph.Traits()
	//	fmt.Printf("Directed: %v, Weighted: %v, Acyclic: %v\n", traits.Directed, traits.Weighted, traits.Acyclic)
	Traits() *Traits

	// Neighbors retrieves the vertices adjacent to the specified vertex in the graph.
	//
	// For directed graphs, it returns the outgoing neighbors of the vertex.
	// For undirected graphs, it returns all vertices connected to the specified vertex.
	//
	// Parameters:
	//   - hash: The identifier of the vertex whose neighbors are to be retrieved.
	//
	// Returns:
	//   - A slice of Vertex[K, T] values that are adjacent to the specified vertex.
	//   - An error if the vertex does not exist in the graph.
	//
	// Notes:
	//   - If the vertex exists but has no neighbors (e.g., an isolated vertex), the returned slice will be empty.
	//   - For undirected graphs, both vertices connected by an edge are considered neighbors of each other.
	//
	// Example:
	//	neighbors, err := graph.Neighbors("A")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	for _, neighbor := range neighbors {
	//		fmt.Printf("Neighbor ID: %v, Value: %v\n", neighbor.ID, neighbor.Value)
	//	}
	Neighbors(hash K) ([]Vertex[K, T], error)

	// Degree returns the total degree of the specified Vertex.
	//
	// For Undirected graphs, this is the number of adjacent vertices.
	// For directed graphs, this is the sum of in-degree and out-degree.
	//
	// Parameters:
	//   - hash: The identifier of the Vertex.
	//
	// Returns:
	//   - An integer representing the degree of the Vertex.
	//   - An error if the Vertex does not exist or the operation fails.
	//
	// Example:
	//	degree, err := graph.Degree("A")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Degree of A: %d\n", degree)
	Degree(hash K) (int, error)

	// InDegree returns the number of incoming edges to the specified Vertex.
	//
	// Applicable only for directed graphs. For Undirected graphs, this is
	// the number of adjacent vertices and thus identical in value to Degree.
	//
	// Parameters:
	//   - hash: The identifier of the Vertex.
	//
	// Returns:
	//   - An integer representing the in-degree of the Vertex.
	//   - An error if the Vertex does not exist, the graph is Undirected, or the operation fails.
	//
	// Example:
	//	inDegree, err := graph.InDegree("A")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("In-Degree of A: %d\n", inDegree)
	InDegree(hash K) (int, error)

	// OutDegree returns the number of outgoing edges from the specified Vertex.
	//
	// Applicable only for directed graphs. For Undirected graphs, this is
	// the number of adjacent vertices and thus identical in value to Degree.
	//
	// Parameters:
	//   - hash: The identifier of the Vertex.
	//
	// Returns:
	//   - An integer representing the out-degree of the Vertex.
	//   - An error if the Vertex does not exist, the graph is Undirected, or the operation fails.
	//
	// Example:
	//	outDegree, err := graph.OutDegree("A")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Out-Degree of A: %d\n", outDegree)
	OutDegree(hash K) (int, error)
}

// Cursor is an interface representing a position in the graph stream. It allows for serialization
// and deserialization of the cursor state, enabling resumption of streams at a specific point.
//
// Methods:
//   - State: Returns the current state of the cursor as a serialized byte slice.
//   - SetState: Restores the cursor state from a serialized byte slice.
//
// Example:
//
//	cursor := &simple.EmptyCursor()
//	state := cursor.State()
//	fmt.Printf("Serialized Cursor State: %s\n", state)
//
//	err := cursor.SetState(state)
//	if err != nil {
//	    log.Fatalf("Failed to restore cursor state: %v", err)
//	}
//	fmt.Printf("Cursor State Restored")
type Cursor interface {
	// State returns the current state of the cursor as a serialized byte slice.
	State() []byte

	// SetState restores the cursor state from a serialized byte slice.
	SetState(state []byte) error
}

// Cloneable represents a type that can produce a deep copy of itself.
//
// Example:
//
//	original := obj
//	cloned := original.Clone()
//	if !reflect.DeepEqual(original, cloned) {
//		log.Fatal("Clone is not identical to original")
//	}
type Cloneable[T any] interface {
	// Clone creates a deep copy of the object.
	Clone() T
}

// Edge represents a connection between two vertices in a graph.
//
// Example:
//
//	edge := graph.NewEdge("A", "B", graph.EdgeWeight(10))
//	fmt.Printf("Edge from %v to %v\n", edge.Source(), edge.Target())
type Edge[T any] interface {
	Cloneable[Edge[T]]

	// Source is the starting vertex of the edge.
	Source() T

	// Target is the ending vertex of the edge.
	Target() T

	// Properties contains additional information or metadata about the edge,
	// such as weight, capacity, or any custom values.
	Properties() EdgeProperties
}

// EdgeProperties defines properties associated with a graph edge.
//
// Example:
//
//	edgeProps := edge.Properties()
//	weight := edgeProps.Weight()
//	customValue := edgeProps.Items()["custom"]
//	fmt.Printf("Weight: %v, Custom Value: %v\n", weight, customValue)
type EdgeProperties interface {
	Cloneable[EdgeProperties]

	// Items retrieves key-value pairs associated with the edge.
	Items() map[string]any

	// Metadata returns custom user-defined information for the edge.
	Metadata() any

	// Weight specifies the weight of the edge. Default is 0 if not explicitly set.
	Weight() float64
}

// EdgeOption defines a functional option for configuring edge properties.
//
// Example:
//
//	edge := graph.NewEdge("A", "B", graph.EdgeWeight(15), graph.EdgeItem("color", "red"))
type EdgeOption func(EdgeProperties)

// Vertex represents a graph vertex with an identifier, metadata, and properties.
//
// Example:
//
//	vertex := graph.NewVertex("A", graph.VertexWeight(10))
//	fmt.Printf("Vertex ID: %v, Weight: %v\n", vertex.ID(), vertex.Properties().Weight())
type Vertex[K comparable, T any] interface {
	Cloneable[Vertex[K, T]]

	// ID returns the unique identifier of the vertex.
	ID() K

	// Value retrieves the value associated with the vertex.
	Value() T

	// Properties returns the properties associated with the vertex, such as
	// its attributes and weight.
	Properties() VertexProperties
}

// VertexProperties defines properties associated with a graph vertex.
//
// Example:
//
//	vertexProps := vertex.Properties()
//	weight := vertexProps.Weight()
//	metadata := vertexProps.Metadata()
//	fmt.Printf("Weight: %v, Metadata: %v\n", weight, metadata)
type VertexProperties interface {
	// Items retrieves key-value pairs associated with the vertex.
	Items() map[string]any

	// Weight specifies the weight of the vertex. Default is 0 if not explicitly set.
	Weight() float64

	// Metadata returns custom user-defined information for the vertex.
	Metadata() any
}

// VertexOption defines a functional option for configuring vertex properties.
//
// Example:
//
//	vertex := graph.NewVertex("A", graph.VertexWeight(20), graph.VertexItem("type", "root"))
type VertexOption func(VertexProperties)

type Comparable[T any] interface {
	CompareTo(other T) int
}
