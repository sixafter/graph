// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package mst

import (
	"errors"
	"fmt"
	"sort"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
)

var (
	// ErrDirectedGraph is returned when a spanning tree is requested for a directed graph.
	ErrDirectedGraph = errors.New("spanning trees can only be determined for undirected graphs")
)

// MinimumSpanningTree returns a minimum spanning tree within the given graph.
//
// The MST Contains all vertices from the given graph as well as the required
// edges for building the MST. The original graph remains unchanged.
func MinimumSpanningTree[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	return spanningTree(g, false)
}

// MaximumSpanningTree returns a minimum spanning tree within the given graph.
//
// The MST Contains all vertices from the given graph as well as the required
// edges for building the MST. The original graph remains unchanged.
func MaximumSpanningTree[K graph.Ordered, T any](g graph.Interface[K, T]) (graph.Interface[K, T], error) {
	return spanningTree(g, true)
}

// spanningTree computes the minimum or maximum spanning tree of a given graph.
// The spanning tree is constructed using Kruskal's algorithm, which relies on a
// Union-Find data structure for efficiently managing connected components.
//
// Type Parameters:
//   - K: The type used to uniquely identify vertices in the graph (e.g., int, string). Must be comparable.
//   - T: The type of data associated with each vertex in the graph.
//
// Parameters:
//   - g: The input graph, which must be undirected. The graph should implement the `Interface` interface.
//   - maximum: A boolean flag indicating whether to compute a maximum spanning tree (if true) or
//     a minimum spanning tree (if false).
//
// Returns:
//   - A new graph representing the spanning tree. The graph retains the same structure and properties
//     as the input graph but Contains only the edges of the spanning tree.
//   - An error if the input graph is directed or if any operation (e.g., accessing vertices, adding edges) fails.
//
// Errors:
//   - Returns an error if the input graph is directed, as spanning trees are only defined for undirected graphs.
//   - Returns an error if the adjacency map cannot be retrieved or if vertices or edges cannot be added
//     to the resulting spanning tree.
//
// Algorithm Details:
//  1. Checks if the input graph is undirected. If not, returns an error.
//  2. Extracts the adjacency map and initializes a Union-Find data structure to manage connected components.
//  3. Adds all vertices from the input graph to the spanning tree graph, copying their properties.
//  4. Sorts edges by weight in ascending order for a minimum spanning tree or descending order for a maximum spanning tree.
//  5. Iterates through the sorted edges, adding each edge to the spanning tree if it connects two previously unconnected components.
//  6. Returns the resulting spanning tree.
//
// Example:
//
//	g := SomeUndirectedGraph() // Interface implementation
//	mst, err := spanningTree(g, false) // Compute minimum spanning tree
//	if err != nil {
//	    log.Fatalf("Error computing spanning tree: %v", err)
//	}
//	fmt.Println("Minimum Spanning Tree:", mst)
func spanningTree[K graph.Ordered, T any](g graph.Interface[K, T], maximum bool) (graph.Interface[K, T], error) {
	if g.Traits().IsDirected {
		return nil, ErrDirectedGraph
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	edges := make([]graph.Edge[K], 0)
	subtrees := newUnionFind[K]()

	mst, err := simple.NewLike(g)
	if err != nil {
		return nil, fmt.Errorf("failed to create new graph: %w", err)
	}

	for v, adjacency := range adjacencyMap {
		var vertex graph.Vertex[K, T]
		vertex, err = g.Vertex(v)
		if err != nil {
			return nil, fmt.Errorf("failed to get vertex %v: %w", v, err)
		}

		err = mst.AddVertex(vertex)
		if err != nil {
			return nil, fmt.Errorf("failed to Add vertex %v: %w", v, err)
		}

		subtrees.Add(v)

		for _, edge := range adjacency {
			edges = append(edges, edge)
		}
	}

	if maximum {
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].Properties().Weight() > edges[j].Properties().Weight()
		})
	} else {
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].Properties().Weight() < edges[j].Properties().Weight()
		})
	}

	for _, edge := range edges {
		sourceRoot := subtrees.Find(edge.Source())
		targetRoot := subtrees.Find(edge.Target())

		if sourceRoot != targetRoot {
			subtrees.Union(sourceRoot, targetRoot)

			if err = mst.AddEdge(edge.Clone()); err != nil {
				return nil, fmt.Errorf("failed to Add edge (%v, %v): %w", edge.Source(), edge.Target(), err)
			}
		}
	}

	return mst, nil
}

// UnionFind implements the Union-Find (or Disjoint-Set) data structure,
// which is used to efficiently perform Union and Find operations on disjoint sets.
//
// The Union-Find data structure is commonly used in algorithms for tasks such as
// finding connected components, Kruskal's algorithm for minimum spanning trees,
// and more.
//
// Type Parameters:
//   - K: The type of the elements in the sets. It must be comparable.
type UnionFind[K graph.Ordered] struct {
	// parents maps each element (of type K) to its parent in the Union-Find structure.
	// If an element is its own parent, it is the root of its set.
	parents map[K]K
}

// newUnionFind creates a new Union-Find (disjoint set) data structure
// initialized with the given vertices. Each vertex starts as its own parent.
//
// Parameters:
//   - vertices: The vertices to initialize in the Union-Find structure.
//
// Returns:
//   - A pointer to a new `UnionFindunionFind` instance.
//
// Complexity: O(Items), where Items is the number of vertices.
//
// Example:
//
//	uf := newUnionFind(1, 2, 3, 4)
//	fmt.Println(uf.Find(1)) // Output: 1
func newUnionFind[K graph.Ordered](vertices ...K) *UnionFind[K] {
	u := &UnionFind[K]{
		parents: make(map[K]K, len(vertices)),
	}

	for _, vertex := range vertices {
		u.parents[vertex] = vertex
	}

	return u
}

// Add inserts a new vertex into the Union-Find data structure, initializing
// it as its own parent.
//
// Parameters:
//   - vertex: The vertex to Add to the Union-Find structure.
//
// Complexity: O(1).
//
// Example:
//
//	uf.Add(5)
func (u *UnionFind[K]) Add(vertex K) {
	u.parents[vertex] = vertex
}

// Union merges the sets containing `vertex1` and `vertex2` by connecting their roots.
// If the vertices already belong to the same set, the operation is a no-op.
//
// Parameters:
//   - vertex1: A vertex in the first set to merge.
//   - vertex2: A vertex in the second set to merge.
//
// Complexity: O(α(Items)), where α is the inverse Ackermann function, effectively constant.
//
// Example:
//
//	uf.Union(1, 2)
func (u *UnionFind[K]) Union(vertex1, vertex2 K) {
	root1 := u.Find(vertex1)
	root2 := u.Find(vertex2)

	if root1 == root2 {
		return
	}

	u.parents[root2] = root1
}

// Find locates the root (representative) of the set containing the given vertex.
// It also applies path compression to optimize the structure of the Union-Find
// data structure for future Find calls.
//
// The function operates in two phases:
//  1. Traverse the tree to find the root of the set containing the vertex. The
//     root is identified as the element that is its own parent in the `parents` map.
//  2. Perform path compression, flattening the tree by making every node on the
//     path from the vertex to the root point directly to the root. This reduces
//     the depth of the tree and ensures that future Find operations for these
//     nodes are faster.
//
// Time Complexity:
//   - Without path compression, the worst-case time complexity of Find is O(n)
//     for a deeply nested tree.
//   - With path compression, the amortized time complexity is O(α(n)), where
//     α(n) is the inverse Ackermann function, which grows extremely slowly and
//     is effectively constant for all practical inputs.
//
// Parameters:
//   - vertex: The element whose set representative (root) is to be found.
//
// Returns:
//   - The root of the set containing the given vertex.
//
// Example Usage:
//
//	uf := NewUnionFind[int]()
//	uf.Union(1, 2)
//	uf.Union(2, 3)
//	fmt.Println(uf.Find(1)) // Output: 3 (assuming 3 becomes the root)
//	fmt.Println(uf.Find(2)) // Output: 3 (optimized due to path compression)
//
// Path Compression Explanation:
//
//	During the traversal, each node encountered in the path from the given vertex
//	to the root is updated to point directly to the root. This reduces the depth
//	of the tree, allowing subsequent Find operations to complete in nearly
//	constant time.
func (u *UnionFind[K]) Find(vertex K) K {
	// Phase 1: Traverse upward to find the root of the set.
	root := vertex
	for u.parents[root] != root {
		root = u.parents[root]
	}

	// Phase 2: Perform path compression to flatten the tree structure.
	current := vertex
	for u.parents[current] != root {
		parent := u.parents[current]
		u.parents[current] = root // Update the parent to point directly to the root.
		current = parent
	}

	return root
}
