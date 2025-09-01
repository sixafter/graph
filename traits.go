// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package graph

// Traits represents a set of graph traits and types, such as whether the graph
// is directed or acyclic. These traits determine the behavior and properties of
// the graph. They can be set during graph creation using functional options.
//
// Example:
//
//	// Create a directed graph.
//	g := simple.New(graph.IntHash, graph.Directed())
type Traits struct {
	// IsAcyclic indicates whether the graph is acyclic. An acyclic graph does not
	// contain cycles (closed paths).
	IsAcyclic bool

	// IsDirected indicates whether the graph is DirectedGraph. In a DirectedGraph graph,
	// edges have a specific direction from a source to a target.
	IsDirected bool

	// IsMultiGraph indicates whether the graph is a multigraph. A multigraph is a
	// graph that allows multiple edges between the same pair of vertices.
	IsMultiGraph bool

	// IsRooted indicates whether the graph is rooted. A rooted graph is a graph
	// with a designated root node, common in tree structures.
	IsRooted bool

	// IsWeighted indicates whether the graph is weighted. Weighted graphs have
	// edges with associated weights.
	IsWeighted bool

	// PreventCycles indicates whether the graph proactively prevents the creation
	// of cycles. This adds checks during edge creation to ensure cycles are not
	// introduced.
	PreventCycles bool
}

type TraitsOption func(*Traits)

// Acyclic creates an acyclic graph. An acyclic graph does not contain cycles.
// Note: This does not prevent cycles from being created explicitly. Use PreventCycles
// to ensure cycles are avoided during edge creation.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.Acyclic())
func Acyclic() TraitsOption {
	return func(t *Traits) {
		t.IsAcyclic = true
	}
}

// Directed creates a DirectedGraph graph. In a DirectedGraph graph, edges have a specific
// direction from a source to a target. This functional option sets the IsDirected
// field in Traits.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.Directed())
func Directed() TraitsOption {
	return func(t *Traits) {
		t.IsDirected = true
	}
}

// MultiGraph creates a multigraph. A multigraph is a graph that allows multiple
// edges between the same pair of vertices. This functional option sets the
// IsMultiGraph field in Traits.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.MultiGraph())
func MultiGraph() TraitsOption {
	return func(t *Traits) {
		t.IsMultiGraph = true
	}
}

// PreventCycles creates an acyclic graph that explicitly prevents the creation
// of cycles. This functional option adds checks during edge creation to
// proactively avoid cycles. Using PreventCycles can affect performance during
// operations like AddEdge.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.PreventCycles())
func PreventCycles() TraitsOption {
	return func(t *Traits) {
		Acyclic()(t)
		t.PreventCycles = true
	}
}

// Rooted creates a rooted graph. Rooted graphs have a designated root DefaultVertex,
// commonly used in tree metadata structures.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.Rooted())
func Rooted() TraitsOption {
	return func(t *Traits) {
		t.IsRooted = true
	}
}

// Tree is a shorthand for creating a rooted, acyclic graph. It combines the
// Rooted and Acyclic functional options. Trees are commonly used in metadata
// structures like binary trees or general trees.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.Tree())
func Tree() TraitsOption {
	return func(t *Traits) {
		Acyclic()(t)
		Rooted()(t)
	}
}

// Weighted creates a weighted graph. Weighted graphs have edges with associated
// weights, which can be set using the Edge or AddEdge methods.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.Weighted())
func Weighted() TraitsOption {
	return func(t *Traits) {
		t.IsWeighted = true
	}
}

// Equals returns true if the provided Traits are equal to the current Traits.
func (t *Traits) Equals(other *Traits) bool {
	if other == nil {
		return false
	}
	return t.IsAcyclic == other.IsAcyclic &&
		t.IsDirected == other.IsDirected &&
		t.IsMultiGraph == other.IsMultiGraph &&
		t.IsRooted == other.IsRooted &&
		t.IsWeighted == other.IsWeighted &&
		t.PreventCycles == other.PreventCycles
}

// Clone creates a deep copy of the provided Traits.
func (t *Traits) Clone() *Traits {
	return &Traits{
		IsAcyclic:     t.IsAcyclic,
		IsDirected:    t.IsDirected,
		IsMultiGraph:  t.IsMultiGraph,
		IsRooted:      t.IsRooted,
		IsWeighted:    t.IsWeighted,
		PreventCycles: t.PreventCycles,
	}
}
