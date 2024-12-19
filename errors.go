// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package graph

import (
	"errors"
)

var (
	// ErrAdjacencyMap is used as the base error when failing to retrieve the adjacency map.
	ErrAdjacencyMap = errors.New("failed to get adjacency map")

	// ErrVertexRetrieval is used as the base error when failing to retrieve a vertex.
	ErrVertexRetrieval = errors.New("failed to get vertex")

	// ErrAddVertex is used as the base error when failing to add a vertex.
	ErrAddVertex = errors.New("failed to add vertex")

	// ErrAddEdge is used as the base error when failing to add an edge.
	ErrAddEdge = errors.New("failed to add edge")

	// ErrCloneGraph is returned when cloning a graph fails.
	ErrCloneGraph = errors.New("failed to clone the graph")

	// ErrCyclicGraph indicates that the graph Contains cycles and the operation cannot proceed.
	ErrCyclicGraph = errors.New("operation cannot be performed on graph with cycles")

	// ErrEdgeAlreadyExists is returned when attempting to add an edge that already exists.
	ErrEdgeAlreadyExists = errors.New("edge already exists")

	// ErrEdgeCreatesCycle is returned when an edge would create a cycle in a graph
	// where cycles are not allowed.
	ErrEdgeCreatesCycle = errors.New("edge would create a cycle")

	// ErrEdgeNotFound is returned when an edge is not found in the graph.
	ErrEdgeNotFound = errors.New("edge not found")

	// ErrFailedToAddEdge is returned when adding an edge to the graph fails.
	ErrFailedToAddEdge = errors.New("failed to add edge")

	// ErrFailedToAddEdges is returned when edges cannot be added during a graph operation.
	ErrFailedToAddEdges = errors.New("failed to add edges")

	// ErrFailedToAddVertex is returned when adding a vertex to the graph fails.
	ErrFailedToAddVertex = errors.New("failed to add vertex")

	// ErrFailedToAddVertices is returned when vertices cannot be added during a graph operation.
	ErrFailedToAddVertices = errors.New("failed to add vertices")

	// ErrFailedToCloneGraph indicates a failure in cloning the graph.
	ErrFailedToCloneGraph = errors.New("failed to clone the graph")

	// ErrFailedToGetAdjacencyMap indicates a failure in retrieving the graph's adjacency map.
	ErrFailedToGetAdjacencyMap = errors.New("failed to get adjacency map")

	// ErrFailedToGetEdges is returned when the edge list of a graph cannot be retrieved.
	ErrFailedToGetEdges = errors.New("failed to get edges")

	// ErrFailedToGetGraphOrder indicates a failure in retrieving the graph's order.
	ErrFailedToGetGraphOrder = errors.New("failed to get graph order")

	// ErrFailedToGetPredecessorMap indicates a failure in retrieving the graph's predecessor map.
	ErrFailedToGetPredecessorMap = errors.New("failed to get predecessor map")

	// ErrFailedToGetVertex is returned when a vertex cannot be retrieved from the graph.
	ErrFailedToGetVertex = errors.New("failed to get vertex")

	// ErrFailedToListEdges is returned when the graph fails to list edges during an operation.
	ErrFailedToListEdges = errors.New("failed to list edges")

	// ErrFailedToListVertices is returned when the graph fails to list vertices during an operation.
	ErrFailedToListVertices = errors.New("failed to list vertices")

	// ErrFailedToRemoveEdge is returned when removing an edge from the graph fails.
	ErrFailedToRemoveEdge = errors.New("failed to remove edge")

	// ErrGetAdjacencyMap is returned when the adjacency map retrieval fails.
	ErrGetAdjacencyMap = errors.New("failed to get adjacency map")

	// ErrGetVertex is returned when a vertex cannot be retrieved.
	ErrGetVertex = errors.New("failed to get vertex")

	// ErrPredecessorMapFailed is returned when there is an error obtaining the predecessor
	// map of a graph. The predecessor map is used for operations like cycle detection.
	ErrPredecessorMapFailed = errors.New("could not get predecessor map")

	// ErrSameSourceAndTarget is returned when the source and target vertices of an operation
	// are the same. This is typically used in scenarios where cycles are being detected or avoided.
	ErrSameSourceAndTarget = errors.New("source and target vertices are the same")

	// ErrSCCDetectionNotDirected is returned when an attempt is made to detect strongly
	// connected components (SCCs) in a graph that is not a directed graph. SCC detection is only
	// valid for DirectedGraph graphs.
	ErrSCCDetectionNotDirected = errors.New("strongly connected components (SCCs) can only be detected in directed graph graphs")

	// ErrTargetNotReachable is returned when the target vertex is not reachable
	// from the source vertex in a graph operation such as ShortestPath.
	ErrTargetNotReachable = errors.New("target vertex not reachable from source")

	// ErrUndirectedGraph indicates that the operation cannot be performed on an Undirected graph.
	ErrUndirectedGraph = errors.New("operation cannot be performed on Undirected graph")

	// ErrVertexAlreadyExists is returned when attempting to add a vertex that already exists.
	ErrVertexAlreadyExists = errors.New("vertex already exists")

	// ErrVertexHasEdges is returned when trying to remove a vertex that still has edges.
	ErrVertexHasEdges = errors.New("vertex has edges")

	// ErrVertexNotFound is returned when a vertex is not found in the graph.
	ErrVertexNotFound = errors.New("vertex not found")

	// ErrGraphTypeMismatch is returned when attempting to perform set operations
	// on graphs with differing types or traits.
	ErrGraphTypeMismatch = errors.New("graph type mismatch")

	ErrNilInputGraph = errors.New("input graph cannot be nil")
)
