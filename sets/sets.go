// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"github.com/sixafter/graph"
)

// Interface defines an interface for performing set-theoretic operations on graphs.
// These operations include union, intersection, difference, and comparisons
// like subset, superset, and equality. The interface is generic and supports
// graphs with any key and vertex types.
//
// Type Parameters:
//   - K: The type of keys used to uniquely identify vertices in the graph (e.g., int, string). Must be comparable.
//   - T: The type of data associated with each vertex in the graph.
type Interface[K graph.Ordered, T any] interface {
	// Union computes the union of two graphs, `g` and `h`.
	// The resulting graph Contains all vertices and edges from both `g` and `h`.
	//
	// Parameters:
	//   - g: The first input graph.
	//   - h: The second input graph.
	//
	// Returns:
	//   - A new graph representing the union of `g` and `h`.
	//   - An error if the operation fails (e.g., incompatible graph types).
	Union(g, h graph.Interface[K, T]) (graph.Interface[K, T], error)

	// Intersection computes the intersection of two graphs, `g` and `h`.
	// The resulting graph Contains only the vertices and edges that are common to both `g` and `h`.
	//
	// Parameters:
	//   - g: The first input graph.
	//   - h: The second input graph.
	//
	// Returns:
	//   - A new graph representing the intersection of `g` and `h`.
	//   - An error if the operation fails (e.g., incompatible graph types).
	Intersection(g, h graph.Interface[K, T]) (graph.Interface[K, T], error)

	// Difference computes the difference of two graphs, `g` and `h`.
	// The resulting graph Contains vertices and edges from `g` that are not in `h`.
	//
	// Parameters:
	//   - g: The first input graph (the minuend).
	//   - h: The second input graph (the subtrahend).
	//
	// Returns:
	//   - A new graph representing the difference of `g` and `h`.
	//   - An error if the operation fails (e.g., incompatible graph types).
	Difference(g, h graph.Interface[K, T]) (graph.Interface[K, T], error)

	// IsSubset checks whether graph `g` is a subset of graph `h`.
	// A subset means all vertices and edges of `g` are contained in `h`.
	//
	// Parameters:
	//   - g: The potential subset graph.
	//   - h: The potential superset graph.
	//
	// Returns:
	//   - A boolean indicating whether `g` is a subset of `h`.
	//   - An error if the operation fails (e.g., incompatible graph types).
	IsSubset(g, h graph.Interface[K, T]) (bool, error)

	// IsSuperset checks whether graph `g` is a superset of graph `h`.
	// A superset means all vertices and edges of `h` are contained in `g`.
	//
	// Parameters:
	//   - g: The potential superset graph.
	//   - h: The potential subset graph.
	//
	// Returns:
	//   - A boolean indicating whether `g` is a superset of `h`.
	//   - An error if the operation fails (e.g., incompatible graph types).
	IsSuperset(g, h graph.Interface[K, T]) (bool, error)

	// Equals checks whether two graphs, `g` and `h`, are equal.
	// Equality means both graphs have the same vertices and edges.
	//
	// Parameters:
	//   - g: The first input graph.
	//   - h: The second input graph.
	//
	// Returns:
	//   - A boolean indicating whether the two graphs are equal.
	//   - An error if the operation fails (e.g., incompatible graph types).
	Equals(g, h graph.Interface[K, T]) (bool, error)

	// Complement computes the complement of a given graph `g`.
	// The resulting graph Contains the same vertices as `g` but with the edges inverted.
	// Edges that exist in `g` are removed, and edges that do not exist in `g` are added.
	//
	// Parameters:
	//   - g: The input graph.
	//
	// Returns:
	//   - A new graph representing the complement of `g`.
	//   - An error if the operation fails (e.g., graph is invalid or inaccessible).
	Complement(g graph.Interface[K, T]) (graph.Interface[K, T], error)

	// SymmetricDifference computes the symmetric difference of two graphs, `g` and `h`.
	// The resulting graph Contains vertices and edges that are in either `g` or `h` but not in both.
	//
	// Parameters:
	//   - g: The first input graph.
	//   - h: The second input graph.
	//
	// Returns:
	//   - A new graph representing the symmetric difference of `g` and `h`.
	//   - An error if the operation fails (e.g., incompatible graph types).
	SymmetricDifference(g, h graph.Interface[K, T]) (graph.Interface[K, T], error)

	// IsDisjoint checks whether two graphs, `g` and `h`, have no vertices or edges in common.
	//
	// Parameters:
	//   - g: The first input graph.
	//   - h: The second input graph.
	//
	// Returns:
	//   - A boolean indicating whether `g` and `h` are disjoint (i.e., have no overlap).
	//   - An error if the operation fails (e.g., incompatible graph types).
	IsDisjoint(g, h graph.Interface[K, T]) (bool, error)
}
