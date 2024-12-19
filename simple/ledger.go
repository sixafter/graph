// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"github.com/sixafter/graph"
)

// ledger defines the interface for managing both vertices and edges in a graph.
// It combines methods for adding, modifying, deleting, and retrieving vertices and edges,
// as well as counting and listing them.
//
// Type Parameters:
//   - K: The type used to uniquely identify vertices (keys), must be comparable.
//   - T: The type of metadata stored in each Vertex.
type ledger[K comparable, T any] interface {
	// New creates and returns a new instance of the ledger implementation.
	// This method allows for creating fresh, independent ledger instances
	// that adhere to the ledger interface.
	New() (ledger[K, T], error)

	// Vertex Management

	// AddVertex adds a Vertex with the specified key, value, and properties to the graph.
	// If the Vertex already exists, ErrVertexAlreadyExists must be returned.
	//
	// Parameters:
	//   - key: The unique identifier for the Vertex.
	//   - value: The value associated with the Vertex.
	//   - properties: The properties of the Vertex, such as v and weight.
	//
	// Returns:
	//   - An error if the Vertex already exists or the operation fails.
	AddVertex(key K, value T, properties graph.VertexProperties) error

	// FindVertex retrieves a Vertex and its associated properties by its key.
	// If the Vertex does not exist, ErrVertexNotFound must be returned.
	//
	// Parameters:
	//   - key: The unique identifier of the Vertex to retrieve.
	//
	// Returns:
	//   - The Vertex value, its properties, and an error if the Vertex is not found.
	FindVertex(key K) (T, graph.VertexProperties, error)

	// ModifyVertex updates the properties of an existing Vertex.
	// If the Vertex does not exist, ErrVertexNotFound must be returned.
	//
	// Parameters:
	//   - key: The unique identifier of the Vertex to modify.
	//   - properties: The new properties to associate with the Vertex.
	//
	// Returns:
	//   - An error if the Vertex does not exist or the operation fails.
	ModifyVertex(key K, properties graph.VertexProperties) error

	// RemoveVertex removes a Vertex by its key.
	// If the Vertex does not exist, ErrVertexNotFound must be returned.
	// If the Vertex has edges connected to it, ErrVertexHasEdges must be returned.
	//
	// Parameters:
	//   - key: The unique identifier of the Vertex to delete.
	//
	// Returns:
	//   - An error if the Vertex does not exist, has connected edges, or the operation fails.
	RemoveVertex(key K) error

	// ListVertices retrieves all Vertex keys in the graph.
	//
	// Returns:
	//   - A slice of Vertex keys and an error if the operation fails.
	ListVertices() ([]K, error)

	// CountVertices returns the total number of vertices in the graph.
	//
	// Returns:
	//   - The number of vertices and an error if the operation fails.
	CountVertices() (int, error)

	// Edge Management

	// AddEdge adds an edge between the vertices with the specified source and target keys.
	// If either Vertex does not exist, ErrVertexNotFound must be returned for the respective Vertex.
	// If the edge already exists, ErrEdgeAlreadyExists must be returned.
	//
	// Parameters:
	//   - source: The key of the source Vertex.
	//   - target: The key of the target Vertex.
	//   - edge: The edge metadata, including properties such as v and weight.
	//
	// Returns:
	//   - An error if the edge already exists, a Vertex does not exist, or the operation fails.
	AddEdge(source, target K, edge graph.Edge[K]) error

	// FindEdge retrieves the edge between the specified source and target vertices.
	// If the edge does not exist, ErrEdgeNotFound must be returned.
	//
	// Parameters:
	//   - source: The key of the source Vertex.
	//   - target: The key of the target Vertex.
	//
	// Returns:
	//   - The edge metadata and an error if the edge does not exist.
	FindEdge(source, target K) (graph.Edge[K], error)

	// ModifyEdge updates the properties of an existing edge between the specified source and target vertices.
	//
	// Parameters:
	//   - source: The unique identifier of the source Vertex.
	//   - target: The unique identifier of the target Vertex.
	//   - edge: The updated edge metadata, including new properties.
	//
	// Returns:
	//   - An error if the edge does not exist (ErrEdgeNotFound) or the operation fails.
	ModifyEdge(source, target K, edge graph.Edge[K]) error

	// RemoveEdge removes the edge between the specified vertices.
	// If the edge does not exist, ErrEdgeNotFound must be returned.
	//
	// Parameters:
	//   - source: The key of the source Vertex.
	//   - target: The key of the target Vertex.
	//
	// Returns:
	//   - An error if the operation fails.
	RemoveEdge(source, target K) error

	// ListEdges retrieves all edges in the graph as a slice.
	//
	// Returns:
	//   - A slice of edges and an error if the operation fails.
	ListEdges() ([]graph.Edge[K], error)

	// CountEdges returns the total number of edges in the graph.
	//
	// Returns:
	//   - The number of edges and an error if the operation fails.
	CountEdges() (int, error)
}
