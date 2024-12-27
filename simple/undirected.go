// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"context"
	"errors"
	"fmt"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/paths"
)

// undirected represents an undirected graph. It maintains vertices and edges
// in a ledger backend and ensures undirected behavior for edge-related operations.
type undirected[K graph.Ordered, T any] struct {
	hash   graph.Hash[K, T]
	traits *graph.Traits
	store  ledger[K, T]
}

// newUndirected creates and returns a new undirected graph instance.
func newUndirected[K graph.Ordered, T any](hash graph.Hash[K, T], traits *graph.Traits, store ledger[K, T]) (*undirected[K, T], error) {
	return &undirected[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}, nil
}

func (u *undirected[K, T]) Traits() *graph.Traits {
	return u.traits
}

func (u *undirected[K, T]) AddVertex(vertex graph.Vertex[K, T]) error {
	hash := u.hash(vertex.Value())

	return u.store.AddVertex(hash, vertex.Value(), vertex.Properties())
}

func (u *undirected[K, T]) AddVertexWithOptions(value T, options ...graph.VertexOption) error {
	hash := u.hash(value)

	v := NewVertexWithOptions(hash, value, func(p graph.VertexProperties) {
		for _, option := range options {
			option(p)
		}
	})

	return u.AddVertex(v)
}

func (u *undirected[K, T]) Vertex(hash K) (graph.Vertex[K, T], error) {
	vertex, props, err := u.store.FindVertex(hash)
	if err != nil {
		return nil, err
	}

	return NewVertex(hash, vertex, props), nil
}

func (u *undirected[K, T]) SetVertexWithOptions(value T, options ...graph.VertexOption) error {
	hash := u.hash(value)
	_, props, err := u.store.FindVertex(hash)
	if err != nil {
		return err
	}

	for _, option := range options {
		dp, ok := props.(*VertexProperties) // Attempt to assert the type
		if !ok {
			return fmt.Errorf("failed to modify vertex: %T", props)
		}

		option(dp)
	}

	return u.store.ModifyVertex(hash, props)
}

func (u *undirected[K, T]) RemoveVertex(hash K) error {
	return u.store.RemoveVertex(hash)
}

func (u *undirected[K, T]) HasVertex(hash K) (bool, error) {
	_, err := u.Vertex(hash)
	if errors.Is(err, graph.ErrVertexNotFound) {
		return false, nil
	}

	return err == nil, err
}

func (u *undirected[K, T]) AddEdge(edge graph.Edge[K]) error {
	sourceHash := edge.Source()
	if _, _, err := u.store.FindVertex(sourceHash); err != nil {
		return fmt.Errorf("could not find source vertex with hash %v: %w", sourceHash, err)
	}

	targetHash := edge.Target()
	if _, _, err := u.store.FindVertex(targetHash); err != nil {
		return fmt.Errorf("could not find target vertex with hash %v: %w", targetHash, err)
	}

	if _, err := u.Edge(sourceHash, targetHash); !errors.Is(err, graph.ErrEdgeNotFound) {
		return graph.ErrEdgeAlreadyExists
	}

	// If the user opted in to preventing cycles, run a cycle check.
	if u.traits.PreventCycles {
		createsCycle, err := paths.WouldCreateCycle[K, T](u, sourceHash, targetHash)
		if err != nil {
			return fmt.Errorf("check for cycles: %w", err)
		}
		if createsCycle {
			return graph.ErrEdgeCreatesCycle
		}
	}

	err := u.store.AddEdge(sourceHash, targetHash, edge)
	if err != nil {
		return err
	}

	rEdge := NewEdge(targetHash, sourceHash, edge.Properties())

	err = u.store.AddEdge(targetHash, sourceHash, rEdge)
	if err != nil {
		return err
	}

	return nil
}

func (u *undirected[K, T]) AddEdgeWithOptions(sourceHash, targetHash K, options ...graph.EdgeOption) error {
	e := NewEdgeWithOptions(sourceHash, targetHash, func(p graph.EdgeProperties) {
		for _, option := range options {
			option(p)
		}
	})

	return u.AddEdge(e)
}

func (u *undirected[K, T]) AddEdgesFrom(g graph.Interface[K, T]) error {
	edges, err := g.Edges()
	if err != nil {
		return fmt.Errorf("failed to get edges: %w", err)
	}

	for _, edge := range edges {
		e := edge.Clone()
		if err := u.AddEdge(e); err != nil {
			return fmt.Errorf("failed to add (%v, %v): %w", edge.Source(), edge.Target(), err)
		}
	}

	return nil
}

func (u *undirected[K, T]) AddVerticesFrom(g graph.Interface[K, T]) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("failed to get adjacency map: %w", err)
	}

	for hash := range adjacencyMap {
		var vertex graph.Vertex[K, T]
		vertex, err = g.Vertex(hash)
		if err != nil {
			return fmt.Errorf("failed to get vertex %v: %w", hash, err)
		}

		if err = u.AddVertex(vertex); err != nil {
			return fmt.Errorf("failed to add vertex %v: %w", hash, err)
		}
	}

	return nil
}

func (u *undirected[K, T]) Edge(sourceHash, targetHash K) (graph.Edge[T], error) {
	// In an undirected graph, since multi-graphs aren't supported, the edge AB
	// is the same as BA. Therefore, if source[target] cannot be found, this
	// function also looks for target[source].
	edge, err := u.store.FindEdge(sourceHash, targetHash)
	if errors.Is(err, graph.ErrEdgeNotFound) {
		edge, err = u.store.FindEdge(targetHash, sourceHash)
	}

	if err != nil {
		return nil, err
	}

	sourceVertex, _, err := u.store.FindVertex(sourceHash)
	if err != nil {
		return nil, err
	}

	targetVertex, _, err := u.store.FindVertex(targetHash)
	if err != nil {
		return nil, err
	}

	if edge == nil {
		return nil, graph.ErrEdgeNotFound
	}

	return NewEdge(sourceVertex, targetVertex, edge.Properties()), nil
}

type tuple[K comparable] struct {
	source, target K
}

func (u *undirected[K, T]) Edges() ([]graph.Edge[K], error) {
	storedEdges, err := u.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to get edges: %w", err)
	}

	// An undirected graph creates each edge twice internally: The edge (A,B) is
	// stored both as (A,B) and (B,A). The Edges method is supposed to return
	// one of these two edges, because from an outside perspective, it only is
	// a single edge.
	//
	// To achieve this, Edges keeps track of already-added edges. For each edge,
	// it also checks if the reversed edge has already been added - e.g., for
	// an edge (A,B), Edges checks if the edge has been added as (B,A).
	//
	// These reversed edges are built as a custom tuple type, which is then used
	// as a map key for access in O(1) time. It looks scarier than it is.
	edges := make([]graph.Edge[K], 0, len(storedEdges)/2)

	added := make(map[tuple[K]]struct{})

	for _, storedEdge := range storedEdges {
		reversedEdge := tuple[K]{
			source: storedEdge.Target(),
			target: storedEdge.Source(),
		}
		if _, ok := added[reversedEdge]; ok {
			continue
		}

		edges = append(edges, storedEdge)

		addedEdge := tuple[K]{
			source: storedEdge.Source(),
			target: storedEdge.Target(),
		}

		added[addedEdge] = struct{}{}
	}

	return edges, nil
}

func (u *undirected[K, T]) SetEdgeWithOptions(source, target K, options ...graph.EdgeOption) error {
	edge, err := u.store.FindEdge(source, target)
	if err != nil {
		return err
	}

	for _, option := range options {
		p, ok := edge.Properties().(*EdgeProperties) // Attempt to assert the type
		if !ok {
			return fmt.Errorf("failed to modify edge: %T", edge.Properties())
		}

		option(p)
	}

	if err = u.store.ModifyEdge(source, target, edge); err != nil {
		return err
	}

	reversedEdge := NewEdge(target, source, edge.Properties())

	return u.store.ModifyEdge(target, source, reversedEdge)
}

func (u *undirected[K, T]) RemoveEdge(source, target K) error {
	if _, err := u.Edge(source, target); err != nil {
		return err
	}

	if err := u.store.RemoveEdge(source, target); err != nil {
		return fmt.Errorf("failed to remove edge from %v to %v: %w", source, target, err)
	}

	if err := u.store.RemoveEdge(target, source); err != nil {
		return fmt.Errorf("failed to remove edge from %v to %v: %w", target, source, err)
	}

	return nil
}

func (u *undirected[K, T]) AdjacencyMap() (map[K]map[K]graph.Edge[K], error) {
	vertices, err := u.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := u.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[K]map[K]graph.Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]graph.Edge[K])
	}

	for _, edge := range edges {
		m[edge.Source()][edge.Target()] = edge
		m[edge.Target()][edge.Source()] = edge // Ensure edges are bidirectional
	}

	return m, nil
}

func (u *undirected[K, T]) PredecessorMap() (map[K]map[K]graph.Edge[K], error) {
	return u.AdjacencyMap()
}

func (u *undirected[K, T]) Clone() (graph.Interface[K, T], error) {
	traits := u.traits.Clone()
	store, err := u.store.New()
	if err != nil {
		return nil, err
	}

	clone := &undirected[K, T]{
		hash:   u.hash,
		traits: traits,
		store:  store,
	}

	if err := clone.AddVerticesFrom(u); err != nil {
		return nil, fmt.Errorf("failed to add vertices: %w", err)
	}

	if err := clone.AddEdgesFrom(u); err != nil {
		return nil, fmt.Errorf("failed to add edges: %w", err)
	}

	return clone, nil
}

func (u *undirected[K, T]) Order() (int, error) {
	return u.store.CountVertices()
}

func (u *undirected[K, T]) Size() (int, error) {
	edgeCount, err := u.store.CountEdges()

	// Divide by 2 since every add edge operation on undirected graph is counted
	// twice.
	return edgeCount / 2, err
}

func (u *undirected[K, T]) HasEdge(vertex1, vertex2 K) (bool, error) {
	adjacencyMap, err := u.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("%w: %v", graph.ErrGetAdjacencyMap, err)
	}

	// Check vertex1 -> vertex2
	edges, exists := adjacencyMap[vertex1]
	if exists {
		for _, edge := range edges {
			if edge.Target() == vertex2 {
				return true, nil
			}
		}
	}

	// Check vertex2 -> vertex1
	edges, exists = adjacencyMap[vertex2]
	if exists {
		for _, edge := range edges {
			if edge.Target() == vertex1 {
				return true, nil
			}
		}
	}

	return false, nil
}

func (u *undirected[K, T]) Hash() graph.Hash[K, T] {
	return u.hash
}

func (u *undirected[K, T]) Vertices() ([]graph.Vertex[K, T], error) {
	vHashes, err := u.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToListVertices, err)
	}

	vertices := make([]graph.Vertex[K, T], len(vHashes))
	for i, hash := range vHashes {
		var vertex graph.Vertex[K, T]
		vertex, err = u.Vertex(hash)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetVertex, hash)
		}
		vertices[i] = vertex
	}

	return vertices, nil
}

func (u *undirected[K, T]) Neighbors(hash K) ([]graph.Vertex[K, T], error) {
	adjacencyMap, err := u.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	neighborsMap, exists := adjacencyMap[hash]
	if !exists {
		return nil, fmt.Errorf("%w: vertex not found", graph.ErrVertexNotFound)
	}

	neighbors := make([]graph.Vertex[K, T], 0, len(neighborsMap))
	for neighborHash := range neighborsMap {
		vertex, err := u.Vertex(neighborHash)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to get neighbor vertex", graph.ErrFailedToGetVertex)
		}
		neighbors = append(neighbors, vertex)
	}

	return neighbors, nil
}

func (u *undirected[K, T]) Degree(hash K) (int, error) {
	adjacencyMap, err := u.AdjacencyMap()
	if err != nil {
		return 0, fmt.Errorf("%w: failed to retrieve adjacency map", graph.ErrFailedToGetEdges)
	}

	edges, exists := adjacencyMap[hash]
	if !exists {
		return 0, fmt.Errorf("%w: vertex not found", graph.ErrVertexNotFound)
	}

	return len(edges), nil
}

func (u *undirected[K, T]) InDegree(hash K) (int, error) {
	return u.Degree(hash)
}

func (u *undirected[K, T]) OutDegree(hash K) (int, error) {
	return u.Degree(hash)
}

// StreamEdgesWithContext streams edges from the undirected graph in paginated batches. The stream can be canceled
// or timed out using a context, and the state of the cursor can be used to resume from the last point.
func (u *undirected[K, T]) StreamEdgesWithContext(ctx context.Context, cursor graph.Cursor, limit int, ch chan<- graph.Edge[K]) (graph.Cursor, error) {
	defer close(ch) // Ensure the channel is closed when the function returns

	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}

	// Deserialize the cursor state to determine the starting position
	var currentPosition int64
	if err := cursor.SetState(cursor.State()); err != nil {
		return nil, err
	}
	state, ok := cursor.(*Cursor)
	if !ok {
		return nil, errors.New("invalid cursor type")
	}
	currentPosition = state.position

	// Fetch edges from the underlying ledger
	edges, err := u.Edges()
	if err != nil {
		return nil, err
	}

	// Stream edges in batches
	for currentPosition < int64(len(edges)) {
		end := currentPosition + int64(limit)
		if end > int64(len(edges)) {
			end = int64(len(edges))
		}

		for _, edge := range edges[currentPosition:end] {
			select {
			case <-ctx.Done(): // Handle cancellation
				return cursor, ctx.Err()
			case ch <- edge: // Send each edge to the channel
			}
		}

		currentPosition = end
		state.position = currentPosition
	}

	return cursor, nil
}

// StreamVerticesWithContext streams vertices from the undirected graph in paginated batches. The stream can be canceled
// or timed out using a context, and the state of the cursor can be used to resume from the last point.
func (u *undirected[K, T]) StreamVerticesWithContext(ctx context.Context, cursor graph.Cursor, limit int, ch chan<- []graph.Vertex[K, T]) (graph.Cursor, error) {
	defer close(ch) // Ensure the channel is closed when the function returns

	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}

	// Deserialize the cursor state to determine the starting position
	var currentPosition int64
	if err := cursor.SetState(cursor.State()); err != nil {
		return nil, err
	}
	state, ok := cursor.(*Cursor)
	if !ok {
		return nil, errors.New("invalid cursor type")
	}
	currentPosition = state.position

	// Fetch vertices from the underlying ledger
	vertices, err := u.Vertices()
	if err != nil {
		return nil, err
	}

	// Stream vertices in batches
	for currentPosition < int64(len(vertices)) {
		end := currentPosition + int64(limit)
		if end > int64(len(vertices)) {
			end = int64(len(vertices))
		}

		batch := vertices[currentPosition:end]
		select {
		case <-ctx.Done(): // Handle cancellation
			return cursor, ctx.Err()
		case ch <- batch: // Send batch to channel
		}

		currentPosition = end
		state.position = currentPosition
	}

	return cursor, nil
}
