// Copyright (c) 2024-2025 Six After, Inc
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

// directedGraph represents a generic directed graph structure with nodes and edges.
// It uses a hash function to map keys to values, traits for graph characteristics,
// and a ledger mechanism to manage the graph's nodes and edges.
//
// Type Parameters:
//   - K: The type of keys used to identify nodes in the graph. It must be comparable.
//   - T: The type of values stored in the nodes of the graph.
type directedGraph[K graph.Ordered, T any] struct {
	hash   graph.Hash[K, T]
	traits *graph.Traits
	store  ledger[K, T]
}

// newDirectedGraph creates a new directedGraph graph with the given hash function, traits, and
//
// Parameters:
//   - hash: A function that computes a unique hash for vertices.
//   - traits: Interface traits that define properties such as directedness and cycle prevention.
//   - ledger: A memory ledger backend for the graph.
//
// Returns:
//   - A pointer to a new directedGraph graph instance.
//
// Example:
//
//	graph := newDirectedGraph(IntHash, &Traits{IsDirected: true}, newMemoryStore())
func newDirectedGraph[K graph.Ordered, T any](hash graph.Hash[K, T], traits *graph.Traits, store ledger[K, T]) (*directedGraph[K, T], error) {
	return &directedGraph[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}, nil
}

func (d *directedGraph[K, T]) Traits() *graph.Traits {
	return d.traits
}

func (d *directedGraph[K, T]) AddVertex(vertex graph.Vertex[K, T]) error {
	hash := d.hash(vertex.Value())
	return d.store.AddVertex(hash, vertex.Value(), vertex.Properties())
}

func (d *directedGraph[K, T]) AddVertexWithOptions(value T, options ...graph.VertexOption) error {
	hash := d.hash(value)

	v := NewVertexWithOptions(hash, value, func(p graph.VertexProperties) {
		for _, option := range options {
			option(p)
		}
	})

	return d.AddVertex(v)
}

func (d *directedGraph[K, T]) AddVerticesFrom(g graph.Interface[K, T]) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("%w: %v", graph.ErrFailedToGetAdjacencyMap, err)
	}

	for hash := range adjacencyMap {
		vertex, err := g.Vertex(hash)
		if err != nil {
			return fmt.Errorf("%w: %v", graph.ErrFailedToGetVertex, hash)
		}

		if err = d.AddVertex(vertex); err != nil {
			return fmt.Errorf("%w: %v", graph.ErrFailedToAddVertex, hash)
		}
	}

	return nil
}

func (d *directedGraph[K, T]) Vertex(hash K) (graph.Vertex[K, T], error) {
	vertex, props, err := d.store.FindVertex(hash)
	if err != nil {
		return nil, err
	}

	return NewVertex(hash, vertex, props), nil
}

func (d *directedGraph[K, T]) SetVertexWithOptions(value T, options ...graph.VertexOption) error {
	hash := d.hash(value)
	_, props, err := d.store.FindVertex(hash)
	if err != nil {
		return err
	}

	for _, option := range options {
		p, ok := props.(*VertexProperties)
		if !ok {
			return fmt.Errorf("failed to modify vertex: %T", props)
		}

		option(p)
	}

	return d.store.ModifyVertex(hash, props)
}

func (d *directedGraph[K, T]) HasVertex(hash K) (bool, error) {
	_, err := d.Vertex(hash)
	if errors.Is(err, graph.ErrVertexNotFound) {
		return false, nil
	}

	return err == nil, err
}

func (d *directedGraph[K, T]) RemoveVertex(hash K) error {
	return d.store.RemoveVertex(hash)
}

func (d *directedGraph[K, T]) AddEdge(edge graph.Edge[K]) error {
	source := edge.Source()
	_, _, err := d.store.FindVertex(source)
	if err != nil {
		return fmt.Errorf("%w: %v", graph.ErrVertexNotFound, source)
	}

	target := edge.Target()
	_, _, err = d.store.FindVertex(target)
	if err != nil {
		return fmt.Errorf("%w: %v", graph.ErrVertexNotFound, target)
	}

	if _, err = d.Edge(source, target); !errors.Is(err, graph.ErrEdgeNotFound) {
		return graph.ErrEdgeAlreadyExists
	}

	if d.traits.PreventCycles {
		var createsCycle bool
		createsCycle, err = d.wouldCreateCycle(source, target)
		if err != nil {
			return fmt.Errorf("%w: %v", graph.ErrEdgeCreatesCycle, err)
		}
		if createsCycle {
			return graph.ErrEdgeCreatesCycle
		}
	}

	return d.store.AddEdge(source, target, edge)
}

func (d *directedGraph[K, T]) AddEdgeWithOptions(source, target K, options ...graph.EdgeOption) error {
	e := NewEdgeWithOptions(source, target, func(p graph.EdgeProperties) {
		for _, option := range options {
			option(p)
		}
	})

	return d.AddEdge(e)
}

func (d *directedGraph[K, T]) AddEdgesFrom(g graph.Interface[K, T]) error {
	edges, err := g.Edges()
	if err != nil {
		return fmt.Errorf("%w: %v", graph.ErrFailedToGetEdges, err)
	}

	for _, edge := range edges {
		e := edge.Clone()
		if err = d.AddEdge(e); err != nil {
			return fmt.Errorf("%w (%v, %v): %v", graph.ErrFailedToAddEdge, edge.Source(), edge.Target(), err)
		}
	}

	return nil
}

func (d *directedGraph[K, T]) Edge(source, target K) (graph.Edge[T], error) {
	edge, err := d.store.FindEdge(source, target)
	if err != nil {
		return nil, err
	}

	sourceVertex, _, err := d.store.FindVertex(source)
	if err != nil {
		return nil, err
	}

	targetVertex, _, err := d.store.FindVertex(target)
	if err != nil {
		return nil, err
	}

	if edge == nil {
		return nil, graph.ErrEdgeNotFound
	}

	return NewEdge(sourceVertex, targetVertex, edge.Properties()), nil
}

func (d *directedGraph[K, T]) Edges() ([]graph.Edge[K], error) {
	return d.store.ListEdges()
}

func (d *directedGraph[K, T]) SetEdgeWithOptions(source, target K, options ...graph.EdgeOption) error {
	existingEdge, err := d.store.FindEdge(source, target)
	if err != nil {
		return err
	}

	for _, option := range options {
		ep := existingEdge.Properties()
		dp, ok := ep.(*EdgeProperties) // Attempt to assert the type
		if !ok {
			return fmt.Errorf("failed to modify edge: %T", ep)
		}

		option(dp)
	}

	return d.store.ModifyEdge(source, target, existingEdge)
}

func (d *directedGraph[K, T]) RemoveEdge(source, target K) error {
	if _, err := d.Edge(source, target); err != nil {
		return err
	}

	err := d.store.RemoveEdge(source, target)
	if err != nil {
		return fmt.Errorf("%w: %v -> %v", graph.ErrFailedToRemoveEdge, source, target)
	}

	return nil
}

func (d *directedGraph[K, T]) AdjacencyMap() (map[K]map[K]graph.Edge[K], error) {
	vertices, err := d.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToListVertices, err)
	}

	edges, err := d.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToListEdges, err)
	}

	m := make(map[K]map[K]graph.Edge[K], len(vertices))
	for _, vertex := range vertices {
		m[vertex] = make(map[K]graph.Edge[K])
	}

	for _, edge := range edges {
		m[edge.Source()][edge.Target()] = edge
	}

	return m, nil
}

func (d *directedGraph[K, T]) PredecessorMap() (map[K]map[K]graph.Edge[K], error) {
	vertices, err := d.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToListVertices, err)
	}

	edges, err := d.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToListEdges, err)
	}

	m := make(map[K]map[K]graph.Edge[K], len(vertices))
	for _, vertex := range vertices {
		m[vertex] = make(map[K]graph.Edge[K])
	}

	for _, edge := range edges {
		if _, ok := m[edge.Target()]; !ok {
			m[edge.Target()] = make(map[K]graph.Edge[K])
		}
		m[edge.Target()][edge.Source()] = edge
	}

	return m, nil
}

func (d *directedGraph[K, T]) Clone() (graph.Interface[K, T], error) {
	s, err := d.store.New()
	if err != nil {
		return nil, err
	}

	clone := &directedGraph[K, T]{
		hash:   d.hash,
		traits: d.traits.Clone(),
		store:  s,
	}

	if err := clone.AddVerticesFrom(d); err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToAddVertices, err)
	}

	if err := clone.AddEdgesFrom(d); err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToAddEdges, err)
	}

	return clone, nil
}

func (d *directedGraph[K, T]) Order() (int, error) {
	return d.store.CountVertices()
}

func (d *directedGraph[K, T]) Size() (int, error) {
	return d.store.CountEdges()
}

func (d *directedGraph[K, T]) HasEdge(source, target K) (bool, error) {
	adjacencyMap, err := d.AdjacencyMap()
	if err != nil {
		return false, fmt.Errorf("%w: %v", graph.ErrGetAdjacencyMap, err)
	}

	edges, exists := adjacencyMap[source]
	if !exists {
		return false, nil
	}

	for _, edge := range edges {
		if edge.Target() == target {
			return true, nil
		}
	}

	return false, nil
}

func (d *directedGraph[K, T]) Hash() graph.Hash[K, T] {
	return d.hash
}

func (d *directedGraph[K, T]) Vertices() ([]graph.Vertex[K, T], error) {
	hashes, err := d.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrFailedToListVertices, err)
	}

	vertices := make([]graph.Vertex[K, T], len(hashes))
	for i, hash := range hashes {
		var vertex graph.Vertex[K, T]
		vertex, err = d.Vertex(hash)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", graph.ErrFailedToGetVertex, hash)
		}
		vertices[i] = vertex
	}

	return vertices, nil
}

func (d *directedGraph[K, T]) Neighbors(hash K) ([]graph.Vertex[K, T], error) {
	adjacencyMap, err := d.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	neighborsMap, exists := adjacencyMap[hash]
	if !exists {
		return nil, fmt.Errorf("%w: vertex not found", graph.ErrVertexNotFound)
	}

	neighbors := make([]graph.Vertex[K, T], 0, len(neighborsMap))
	for neighborHash := range neighborsMap {
		vertex, err := d.Vertex(neighborHash)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to get neighbor vertex", graph.ErrFailedToGetVertex)
		}
		neighbors = append(neighbors, vertex)
	}

	return neighbors, nil
}

func (d *directedGraph[K, T]) Degree(key K) (int, error) {
	inDegree, err := d.InDegree(key)
	if err != nil {
		return 0, err
	}

	outDegree, err := d.OutDegree(key)
	if err != nil {
		return 0, err
	}

	return inDegree + outDegree, nil
}

func (d *directedGraph[K, T]) InDegree(hash K) (int, error) {
	predecessorMap, err := d.PredecessorMap()
	if err != nil {
		return 0, fmt.Errorf("%w: failed to retrieve predecessor map", graph.ErrFailedToGetEdges)
	}

	edges, exists := predecessorMap[hash]
	if !exists {
		return 0, fmt.Errorf("%w: vertex not found", graph.ErrVertexNotFound)
	}

	return len(edges), nil
}

func (d *directedGraph[K, T]) OutDegree(hash K) (int, error) {
	adjacencyMap, err := d.AdjacencyMap()
	if err != nil {
		return 0, fmt.Errorf("%w: failed to retrieve adjacency map", graph.ErrFailedToGetEdges)
	}

	edges, exists := adjacencyMap[hash]
	if !exists {
		return 0, fmt.Errorf("%w: vertex not found", graph.ErrVertexNotFound)
	}

	return len(edges), nil
}

func (d *directedGraph[K, T]) wouldCreateCycle(source, target K) (bool, error) {
	// If the underlying ledger implements WouldCreateCycle, use that fast path.
	if cc, ok := d.store.(interface {
		WouldCreateCycle(source, target K) (bool, error)
	}); ok {
		return cc.WouldCreateCycle(source, target)
	}

	// Slow path.
	return paths.WouldCreateCycle(graph.Interface[K, T](d), source, target)
}

// StreamVerticesWithContext streams vertices from the graph in paginated batches. The stream can be canceled
// or timed out using a context, and the state of the cursor can be used to resume from the last point.
func (d *directedGraph[K, T]) StreamVerticesWithContext(ctx context.Context, cursor graph.Cursor, limit int, ch chan<- []graph.Vertex[K, T]) (graph.Cursor, error) {
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
	vertices, err := d.Vertices()
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

// StreamEdgesWithContext streams edges from the graph in paginated batches. The stream can be canceled
// or timed out using a context, and the state of the cursor can be used to resume from the last point.
func (d *directedGraph[K, T]) StreamEdgesWithContext(ctx context.Context, cursor graph.Cursor, limit int, ch chan<- graph.Edge[K]) (graph.Cursor, error) {
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
	edges, err := d.Edges()
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
