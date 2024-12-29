// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package topology

import (
	"fmt"
	"math"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// TarjanFrom identifies all Strongly Connected Components (SCCs) in a directed graph.
//
// In formal graph-theoretic terms, let G = (Items, E) be a directed graph where Items is the set of vertices
// and E is the set of edges. A Strongly Connected Component S ⊆ Items is a maximal subset of vertices
// such that for every pair of vertices u, v ∈ S, there is a directed path from u to v and a directed
// path from v to u. SCCs thus represent subgraphs where connectivity is mutual.
//
// This function uses TarjanFrom’s algorithm, a well-known linear-time procedure (O(|Items| + |E|)) for finding SCCs.
// TarjanFrom’s algorithm performs a single depth-first search (DFS) to compute a unique "discovery index"
// (or timestamp) for each vertex. It also maintains a "low-link" value that tracks the smallest
// discovery index of any vertex reachable from a given vertex, including itself and via back edges.
// When a vertex’s low-link value matches its own discovery index, it signifies the root of an SCC.
//
// Practical applications of SCC computation include analyzing program structure in compilers,
// detecting strongly connected regions in communication networks, and decomposing large graphs
// for further structure analysis.
//
// Parameters:
//   - g: A directed graph implementing graph.Interface[K, T], where K is the vertex key type and T
//     is the vertex data type.
//
// Returns:
//   - [][]K: A slice of slices, where each inner slice corresponds to one strongly connected component.
//   - error: An error if SCC detection fails (e.g., if the graph is undirected).
//
// Example:
//
//	components, err := TarjanFrom(graph)
//	if err != nil {
//		log.Fatal(err)
//	} else {
//		fmt.Printf("Strongly Connected Components: %v\n", components)
//	}
func TarjanFrom[K graph.Ordered, T any](g graph.Interface[K, T]) ([][]K, error) {
	if !g.Traits().IsDirected {
		return nil, graph.ErrSCCDetectionNotDirected
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", graph.ErrAdjacencyMap, err)
	}

	state := &sccState[K]{
		adjacencyMap: adjacencyMap,
		components:   make([][]K, 0),
		stack:        queue.NewStack[K](),
		visited:      make(map[K]struct{}),
		lowLink:      make(map[K]int),
		index:        make(map[K]int),
	}

	for hash := range state.adjacencyMap {
		if _, ok := state.visited[hash]; !ok {
			findSCC(hash, state)
		}
	}

	return state.components, nil
}

type sccState[K comparable] struct {
	adjacencyMap map[K]map[K]graph.Edge[K]
	stack        *queue.Stack[K]
	visited      map[K]struct{}
	lowLink      map[K]int
	index        map[K]int
	components   [][]K
	time         int
}

// findSCC identifies and extracts strongly connected components (SCCs)
// from a directedGraph graph using TarjanFrom's algorithm. It is a recursive depth-first search (DFS)
// based approach that computes the SCCs in a single traversal.
//
// A strongly connected component is a maximal subgraph where any vertex is reachable from
// any other vertex within the same component.
//
// Parameters:
//   - vertexHash: The current vertex being explored, identified by its hash value of type K.
//   - state: A pointer to the SCC computation state, which tracks the DFS Stack, visited vertices,
//     adjacency relationships, and resulting components.
//
// Algorithm Steps:
//  1. Enqueue the current vertex onto the Stack and mark it as visited.
//  2. Assign the vertex an initial index and lowLink value based on the current DFS timestamp.
//  3. Explore all adjacent vertices:
//     - If the adjacent vertex has not been visited, recursively call this function on it.
//     - If the adjacent vertex is already on the Stack, update the lowLink value of the current vertex
//     to reflect the presence of a back edge.
//  4. After visiting all adjacent vertices, check if the current vertex is the root of an SCC:
//     - If its lowLink value equals its index, the SCC is identified.
//     - Dequeue vertices from the Stack until the current vertex is reached and group them as a component.
//  5. Append the identified component to the list of SCCs.
//
// Type Parameters:
//   - K: The type of the hash used to uniquely identify vertices (must be comparable).
//
// Complexity:
//   - Time Complexity: O(Items + E), where Items is the number of vertices and E is the number of edges.
//   - Space Complexity: O(Items) for the Stack and auxiliary maps.
//
// This function is designed to be used as part of a larger SCC computation process,
// typically initiated with a full traversal over all vertices in the graph.
func findSCC[K comparable](vertexHash K, state *sccState[K]) {
	state.stack.Push(vertexHash)
	state.visited[vertexHash] = struct{}{}
	state.index[vertexHash] = state.time
	state.lowLink[vertexHash] = state.time

	state.time++

	for adjacency := range state.adjacencyMap[vertexHash] {
		if _, ok := state.visited[adjacency]; !ok {
			findSCC[K](adjacency, state)

			smallestLowLink := math.Min(
				float64(state.lowLink[vertexHash]),
				float64(state.lowLink[adjacency]),
			)
			state.lowLink[vertexHash] = int(smallestLowLink)
		} else {
			// If the adjacent vertex already is on the Stack, the edge joining
			// the current and the adjacent vertex is a back ege. Therefore, the
			// lowLink value of the vertex has to be updated to the index of the
			// adjacent vertex if it is smaller than the current lowLink value.
			if state.stack.Contains(adjacency) {
				smallestLowLink := math.Min(
					float64(state.lowLink[vertexHash]),
					float64(state.index[adjacency]),
				)
				state.lowLink[vertexHash] = int(smallestLowLink)
			}
		}
	}

	// If the lowLink value of the vertex is equal to its DFS value, this is the
	// head vertex of a strongly connected component that's shaped by the vertex
	// and all vertices on the Stack.
	if state.lowLink[vertexHash] == state.index[vertexHash] {
		var hash K
		var component []K

		for hash != vertexHash {
			hash, _ = state.stack.Pop()

			component = append(component, hash)
		}

		state.components = append(state.components, component)
	}
}
