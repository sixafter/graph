// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package simple

import (
	"strconv"

	"github.com/sixafter/graph"
)

// New creates a new graph with vertices of type T, identified by hash values of
// type K. These hash values will be obtained using the provided hash function.
//
// The graph will use the default in-memory ledger for persisting vertices and
// edges. To use a different [Store], use [newWithStore].
func New[K graph.Ordered, T any](hash graph.Hash[K, T], options ...func(*graph.Traits)) (graph.Interface[K, T], error) {
	s, err := newMemoryStore[K, T]()
	if err != nil {
		return nil, err
	}

	return newWithStore(hash, s, options...)
}

// NewLike creates a graph that is "like" the given graph: It has the same type,
// the same hashing function, and the same traits. The new graph is independent
// of the original graph and uses the default in-memory ledger.
//
//	g := simple.New(graph.IntHash, graph.Directed())
//	h := simple.NewLike(g)
//
// In the example above, h is a new directed graph of integers derived from g.
func NewLike[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	traits := func(t *graph.Traits) {
		*t = *g.Traits().Clone()
	}

	s, err := newMemoryStore[K, T]()
	if err != nil {
		return nil, err
	}

	return newWithStore(g.Hash(), s, traits)
}

// newWithStore creates a new graph using the specified hash function, ledger
// backend, and optional traits.
func newWithStore[K graph.Ordered, T any](hash graph.Hash[K, T], store ledger[K, T], options ...func(*graph.Traits)) (graph.Interface[K, T], error) {
	var p graph.Traits

	for _, option := range options {
		option(&p)
	}

	if p.IsDirected {
		return newDirectedGraph(hash, &p, store)
	}

	return newUndirected(hash, &p, store)
}

// Cursor represents the state of streaming operations in the graph.
type Cursor struct {
	position int64
}

// State serializes the current cursor position into a byte slice by converting it to a string first.
func (c *Cursor) State() []byte {
	return []byte(strconv.FormatInt(c.position, 10))
}

// SetState deserializes the cursor position from a byte slice by converting it back from a string.
func (c *Cursor) SetState(state []byte) error {
	pos, err := strconv.ParseInt(string(state), 10, 64)
	if err != nil {
		return err
	}
	c.position = pos
	return nil
}

// EmptyCursor creates a new cursor with newly initialized state.
func EmptyCursor() graph.Cursor {
	return &Cursor{
		position: 0,
	}
}

// Edge represents a connection between two vertices in a graph.
// It includes the source and target vertices, along with any associated properties.
//
// Type Parameters:
//   - T: The type of the vertices in the graph. It can be any type that
//     represents a Vertex (e.g., int, string, or a custom struct).
type Edge[T any] struct {
	graph.Edge[T]

	// source is the starting Vertex of the edge.
	source T

	// Target is the ending Vertex of the edge.
	target T

	// Properties Contains additional information or metadata about the edge,
	// such as weight, capacity, or any custom v.
	properties graph.EdgeProperties
}

type EdgeOption func(properties *EdgeProperties)

// NewEdge creates a new instance of an edge with the specified source and target.
func NewEdge[T any](source, target T, properties graph.EdgeProperties) graph.Edge[T] {
	return &Edge[T]{
		source:     source,
		target:     target,
		properties: properties,
	}
}

// NewEdgeWithOptions creates a new instance of an edge with the specified source and target.
func NewEdgeWithOptions[T any](source, target T, options ...graph.EdgeOption) graph.Edge[T] {
	properties := EdgeProperties{
		items: make(map[string]any),
	}

	for _, option := range options {
		option(&properties)
	}

	return &Edge[T]{
		source:     source,
		target:     target,
		properties: &properties,
	}
}

func (e *Edge[T]) Source() T {
	return e.source
}

func (e *Edge[T]) Target() T {
	return e.target
}

func (e *Edge[T]) Properties() graph.EdgeProperties {
	return e.properties
}

func (e *Edge[T]) Clone() graph.Edge[T] {
	return &Edge[T]{
		source:     e.source,
		target:     e.target,
		properties: e.properties.Clone(),
	}
}

// EdgeProperties defines the v, weight, and additional metadata associated
// with an edge in the graph.
type EdgeProperties struct {
	graph.EdgeProperties
	items    map[string]any
	metadata any
	weight   float64
}

// EdgeWeight returns a functional option to set the weight of an edge.
func EdgeWeight(weight float64) graph.EdgeOption {
	return func(e graph.EdgeProperties) {
		p := e.(*EdgeProperties)
		p.weight = weight
	}
}

// EdgeItem returns a functional option to add a key-value pair as an
// attribute of an edge.
func EdgeItem(key string, value any) graph.EdgeOption {
	return func(e graph.EdgeProperties) {
		p := e.(*EdgeProperties)
		p.items[key] = value
	}
}

// EdgeItems returns a functional option to set multiple v for an
// edge as a map.
func EdgeItems(items map[string]any) graph.EdgeOption {
	return func(e graph.EdgeProperties) {
		p := e.(*EdgeProperties)
		p.items = items
	}
}

// EdgeData returns a functional option to set additional metadata for an edge.
func EdgeData(data any) graph.EdgeOption {
	return func(e graph.EdgeProperties) {
		p := e.(*EdgeProperties)
		p.metadata = data
	}
}

func (e *EdgeProperties) Items() map[string]any {
	return e.items
}

func (e *EdgeProperties) Metadata() any {
	return e.metadata
}

func (e *EdgeProperties) Weight() float64 {
	return e.weight
}

// Clone creates a deep copy of the edge properties and returns the new instance.
func (e *EdgeProperties) Clone() graph.EdgeProperties {
	clone := EdgeProperties{
		items:    make(map[string]any),
		weight:   e.weight,
		metadata: e.metadata,
	}

	for k, v := range e.items {
		clone.items[k] = v
	}

	return &clone
}

type VertexOption func(*VertexProperties)

type Vertex[K comparable, T any] struct {
	id         K
	value      T
	properties graph.VertexProperties
}

// NewVertex creates a new instance of a Vertex with the specified identifier and value.
func NewVertex[K comparable, T any](id K, value T, properties graph.VertexProperties) graph.Vertex[K, T] {
	return &Vertex[K, T]{
		id:         id,
		value:      value,
		properties: properties,
	}
}

// NewVertexWithOptions creates a new instance of a Vertex with the specified identifier, value, and properties.
func NewVertexWithOptions[K comparable, T any](id K, value T, options ...graph.VertexOption) graph.Vertex[K, T] {
	properties := &VertexProperties{
		v: make(map[string]any),
	}

	for _, option := range options {
		option(properties)
	}

	return &Vertex[K, T]{
		id:         id,
		value:      value,
		properties: properties,
	}
}

// ID returns the unique identifier of the Vertex.
func (v *Vertex[K, T]) ID() K {
	return v.id
}

// Value returns the value associated with the Vertex.
func (v *Vertex[K, T]) Value() T {
	return v.value
}

// Properties returns the properties associated with the Vertex, such as its v and weight.
func (v *Vertex[K, T]) Properties() graph.VertexProperties {
	return v.properties
}

// VertexProperties defines the v and weight of a Vertex in the graph.
//
// Fields:
//   - Attributes: A map of key-value pairs representing custom properties
//     associated with the Vertex.
//   - Weight: An integer representing the weight of the Vertex, which can
//     be used in graph algorithms.
type VertexProperties struct {
	v        map[string]any
	metadata any
	weight   float64
}

// VertexWeight returns a functional option to set the weight of a Vertex.
func VertexWeight(weight float64) graph.VertexOption {
	return func(e graph.VertexProperties) {
		p := e.(*VertexProperties)
		p.weight = weight
	}
}

// VertexItem returns a functional option to add a key-value pair as an
// attribute of a Vertex.
func VertexItem(key string, value any) graph.VertexOption {
	return func(e graph.VertexProperties) {
		p := e.(*VertexProperties)
		p.v[key] = value
	}
}

// VertexItems returns a functional option to set multiple v for
// a Vertex as a map.
func VertexItems(items map[string]any) graph.VertexOption {
	return func(e graph.VertexProperties) {
		p := e.(*VertexProperties)
		p.v = items
	}
}

func VertexMetadata(data any) graph.VertexOption {
	return func(e graph.VertexProperties) {
		p := e.(*VertexProperties)
		p.metadata = data
	}
}

func (p *VertexProperties) Items() map[string]any {
	return p.v
}

func (p *VertexProperties) Weight() float64 {
	return p.weight
}

func (p *VertexProperties) Metadata() any {
	return p.metadata
}

// Clone creates a deep copy of the Vertex properties and returns the new instance.
func (p *VertexProperties) Clone() graph.VertexProperties {
	clone := VertexProperties{
		v:      make(map[string]any),
		weight: p.weight,
	}

	for k, v := range p.v {
		clone.v[k] = v
	}

	return &clone
}

// Clone creates a deep copy of the Vertex and returns the new instance.
func (v *Vertex[K, T]) Clone() graph.Vertex[K, T] {
	clone := Vertex[K, T]{
		id:    v.id,
		value: v.value,
	}

	if v.properties == nil {
		return &clone
	}

	p, ok := v.properties.(*VertexProperties)
	if !ok {
		return &clone
	}

	clone.properties = p.Clone()

	return &clone
}
