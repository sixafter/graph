// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package graph

import (
	"crypto/sha256"
	"hash/fnv"
)

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

// FNV1aHash64 returns the 64-bit FNV-1a hash of the given string.
//
// FNV-1a is a fast, non-cryptographic hash function with good distribution properties,
// suitable for use in hash tables, graph vertex keys, and similar use cases.
// The returned hash value is always the same for the same input string.
//
// This function uses the Go standard library's hash/fnv package, and panics if an
// unexpected error occurs during hashing (which should never happen in practice).
//
// Example:
//
//	key := graph.FNV1aHash64("vertex-id")
func FNV1aHash64(s string) uint64 {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic("unexpected error from FNV hash Write: " + err.Error())
	}
	return h.Sum64()
}

// FNV1aHash128 returns the 128-bit FNV-1a hash of the given string as a [16]byte array.
//
// FNV-1a 128-bit is a fast, non-cryptographic hash with excellent distribution properties.
// The returned [16]byte value is suitable for use as a key in hash tables, graphs, or for deduplication.
//
// This function uses Go's standard library hash/fnv package and panics if an unexpected error occurs.
//
// Example:
//
//	key := graph.FNV1aHash128("vertex-id")
func FNV1aHash128(s string) [16]byte {
	h := fnv.New128a()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic("unexpected error from FNV hash Write: " + err.Error())
	}
	var sum [16]byte
	copy(sum[:], h.Sum(nil))
	return sum
}

// Sha256Hash returns the 256-bit SHA-2 hash of the given string as a [32]byte array.
//
// SHA-256 is a cryptographically secure hash function with excellent distribution properties.
// The returned [32]byte value is suitable for use as a cryptographic fingerprint, secure key, or unique identifier.
//
// This function uses Go's standard library crypto/sha256 package.
//
// Example:
//
//	key := graph.Sha256Hash("vertex-id")
func Sha256Hash(s string) [32]byte {
	return sha256.Sum256([]byte(s))
}
