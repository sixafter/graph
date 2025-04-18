// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package graph

// Hash is a function type that maps a vertex of type T to a hash of type K.
type Hash[K Ordered, T any] func(T) K

// StringHash returns the given string as its hash.
//
// This is a simple identity function that is useful when the string itself
// uniquely identifies the vertex in the graph.
//
// Example:
//
//	g := simple.New(graph.StringHash, graph.Directed())
//	g.AddVertex("A")
//	g.AddVertex("B")
//	g.AddEdge("A", "B")
func StringHash(v string) string {
	return v
}

// IntHash returns the given integer as its hash.
//
// This is a simple identity function that is useful when the integer itself
// uniquely identifies the vertex in the graph.
//
// Example:
//
//	g := simple.New(graph.IntHash, graph.Directed())
//	g.AddVertex(1)
//	g.AddVertex(2)
//	g.AddEdge(1, 2)
func IntHash(v int) int {
	return v
}

// Float64Hash returns the given float as its hash.
//
// This is a simple identity function that is useful when the float itself
// uniquely identifies the vertex in the graph.
//
// Example:
//
//	g := simple.New(graph.Float64Hash, graph.Directed())
//	g.AddVertex(1.0)
//	g.AddVertex(2.0)
//	g.AddEdge(1.0, 2.0)
func Float64Hash(v float64) float64 {
	return v
}

// Float32Hash returns the given float as its hash.
//
// This is a simple identity function that is useful when the float itself
// uniquely identifies the vertex in the graph.
//
// Example:
//
//	g := simple.New(graph.Float32Hash, graph.Directed())
//	g.AddVertex(1.0)
//	g.AddVertex(2.0)
//	g.AddEdge(1.0, 2.0)
func Float32Hash(v float32) float32 {
	return v
}
