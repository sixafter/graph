// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"
	"math"
	"testing"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/simple"
	"github.com/stretchr/testify/assert"
)

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestBrandesBasic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed chain: 1 -> 2 -> 3
	g, err := simple.New(graph.IntHash, graph.Directed())
	is.NoError(err)

	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))

	// Calculate betweenness centrality using BrandesFrom
	bc, err := BrandesFrom(g)
	is.NoError(err)

	// For the chain 1->2->3:
	// Vertex 2 has betweenness = 1.0 (all paths 1->3 go through 2)
	// Vertices 1 & 3 have betweenness = 0.0
	is.True(floatEquals(0.0, bc[g.Hash()(1)]), "Vertex 1 should have BC 0")
	is.True(floatEquals(1.0, bc[g.Hash()(2)]), "Vertex 2 should have BC 1")
	is.True(floatEquals(0.0, bc[g.Hash()(3)]), "Vertex 3 should have BC 0")
}

// TestBrandesUndirectedStar verifies that in an undirected star graph, the center
// has the highest betweenness centrality, and the leaves have zero.
func TestBrandesUndirectedStar(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create an undirected star graph with one center (vertex 1) and four leaves (2,3,4,5).
	g, err := simple.New(graph.IntHash) // defaults to undirected
	is.NoError(err)

	for i := 1; i <= 5; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges (center is 1):
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(1, 4))
	is.NoError(g.AddEdgeWithOptions(1, 5))

	bc, err := BrandesFrom(g)
	is.NoError(err)

	// The center (1) is on the path between every pair of leaves.
	// There are 4 leaves, so #pairs among leaves is C(4,2) = 6.
	// In undirected BrandesFrom, we divide by 2 at the end, so the center’s BC = 6.
	centerKey := g.Hash()(1)
	is.True(floatEquals(6.0, bc[centerKey]), "Center vertex should have BC=6.0")

	// All leaves should have 0
	for i := 2; i <= 5; i++ {
		leafKey := g.Hash()(i)
		is.True(floatEquals(0.0, bc[leafKey]), fmt.Sprintf("Leaf %d should have BC=0", i))
	}
}

// TestBrandesWeightedTriangle checks an undirected weighted triangle (K3).
// For a fully connected triangle, no single vertex lies strictly "between"
// the other two vertices. All betweenness values should be zero.
func TestBrandesWeightedTriangle(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	g, err := simple.New(graph.IntHash, graph.Weighted()) // undirected, weighted
	is.NoError(err)

	// Add vertices: 1, 2, 3
	is.NoError(g.AddVertexWithOptions(1))
	is.NoError(g.AddVertexWithOptions(2))
	is.NoError(g.AddVertexWithOptions(3))

	// For simplicity, all edges have some positive weight
	is.NoError(g.AddEdgeWithOptions(1, 2, simple.EdgeWeight(1.0)))
	is.NoError(g.AddEdgeWithOptions(2, 3, simple.EdgeWeight(2.0)))
	is.NoError(g.AddEdgeWithOptions(1, 3, simple.EdgeWeight(2.0)))

	bc, err := BrandesFrom(g)
	is.NoError(err)

	// In a complete triangle, each vertex connects directly to the other two.
	// There's no "middle" vertex for any pair, so betweenness centrality is 0.
	is.True(floatEquals(0.0, bc[g.Hash()(1)]), "BC(1) should be 0")
	is.True(floatEquals(0.0, bc[g.Hash()(2)]), "BC(2) should be 0")
	is.True(floatEquals(0.0, bc[g.Hash()(3)]), "BC(3) should be 0")
}

// TestBrandesDisjointGraph ensures that vertices in disconnected components
// have zero betweenness centrality because there are no cross-component paths.
func TestBrandesDisjointGraph(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a simple undirected graph with two disjoint components:
	// Component A: 1 -- 2
	// Component B: 3 -- 4
	g, err := simple.New(graph.IntHash)
	is.NoError(err)

	// Add vertices 1,2,3,4
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges within each component
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	bc, err := BrandesFrom(g)
	is.NoError(err)

	// Since components are disconnected from each other, no vertex can lie on a
	// path between any two vertices in the other component. Hence, all BC=0.
	for i := 1; i <= 4; i++ {
		k := g.Hash()(i)
		is.True(floatEquals(0.0, bc[k]), fmt.Sprintf("Vertex %d should have BC=0 in disjoint graph", i))
	}
}

// TestBrandesDirectedDiamond tests BrandesFrom' algorithm on a directed diamond-shaped graph.
func TestBrandesDirectedDiamond(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a directed diamond:
	//
	//    1
	//   / \
	//  v   v
	//  2   3
	//   \ /
	//    v
	//    4
	//
	// Shortest paths between 1 and 4 are:
	//   1->2->4 and 1->3->4
	// Hence 2 and 3 each lie on half of those shortest paths ⇒ BC(2)=0.5, BC(3)=0.5
	g, err := simple.New(graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices
	for i := 1; i <= 4; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add directed edges
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(1, 3))
	is.NoError(g.AddEdgeWithOptions(2, 4))
	is.NoError(g.AddEdgeWithOptions(3, 4))

	bc, err := BrandesFrom(g)
	is.NoError(err)

	// Expect:
	//   BC(1) = 0.0  (source only, not "between")
	//   BC(4) = 0.0  (sink only)
	//   BC(2) = 0.5
	//   BC(3) = 0.5
	is.True(floatEquals(0.0, bc[g.Hash()(1)]), "BC(1) should be 0")
	is.True(floatEquals(0.5, bc[g.Hash()(2)]), "BC(2) should be 0.5")
	is.True(floatEquals(0.5, bc[g.Hash()(3)]), "BC(3) should be 0.5")
	is.True(floatEquals(0.0, bc[g.Hash()(4)]), "BC(4) should be 0")
}

// TestBrandesDirectedCycle tests BrandesFrom' algorithm on a directed cycle.
func TestBrandesDirectedCycle(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a small cycle: 1->2->3->1
	// We’ll see each vertex sits on at least one shortest path between other pairs.
	g, err := simple.New(graph.IntHash, graph.Directed())
	is.NoError(err)

	// Add vertices 1, 2, 3
	for i := 1; i <= 3; i++ {
		is.NoError(g.AddVertexWithOptions(i))
	}

	// Add edges for the cycle
	is.NoError(g.AddEdgeWithOptions(1, 2))
	is.NoError(g.AddEdgeWithOptions(2, 3))
	is.NoError(g.AddEdgeWithOptions(3, 1))

	bc, err := BrandesFrom(g)
	is.NoError(err)

	// Shortest paths:
	//   1->2 (no intermediate)
	//   1->3 => 1->2->3 => intermediate is 2 => BC(2) += 1
	//   2->3 (no intermediate)
	//   2->1 => 2->3->1 => intermediate is 3 => BC(3) += 1
	//   3->1 (no intermediate)
	//   3->2 => 3->1->2 => intermediate is 1 => BC(1) += 1
	//
	// => BC(1)=1, BC(2)=1, BC(3)=1
	is.True(floatEquals(1.0, bc[g.Hash()(1)]), "BC(1) should be 1")
	is.True(floatEquals(1.0, bc[g.Hash()(2)]), "BC(2) should be 1")
	is.True(floatEquals(1.0, bc[g.Hash()(3)]), "BC(3) should be 1")
}
