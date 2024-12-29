// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"fmt"
	"math"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// BrandesFrom computes the betweenness centrality for all vertices in the provided graph g.
// It automatically selects between unweighted and weighted implementations based on the graph's traits.
//
// For undirected graphs, the final centrality scores are halved to account for bidirectional path counting.
//
// The function returns an error if the input graph is nil or if it contains negative-weight edges,
// as Brandes' algorithm does not support graphs with negative weights.
//
// The resulting betweenness centrality scores are returned as a map where each key corresponds
// to a vertex in the graph, and the value represents its centrality score.
//
// Example:
//
//	g, err := simple.New(graph.IntHash, graph.Directed())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Add vertices and edges to g
//	bc, err := BrandesFrom(g)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for vertex, centrality := range bc {
//	    fmt.Printf("Vertex %v has betweenness centrality %f\n", vertex, centrality)
//	}
func BrandesFrom[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	if g.Traits().IsWeighted {
		return brandesWeighted(g)
	}
	return brandesUnweighted(g)
}

// brandesUnweighted computes betweenness centrality for unweighted graphs.
func brandesUnweighted[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Fetch adjacency for quick neighbor lookups
	adj, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("brandes: failed to retrieve adjacency map: %w", err)
	}

	// Get all vertices in the graph.
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("brandes: failed to retrieve vertices: %w", err)
	}

	// Initialize betweenness centrality scores to 0.
	BC := make(map[K]float64, len(vertices))
	for _, v := range vertices {
		k := g.Hash()(v.Value()) // Derive key using the hash function
		BC[k] = 0.0
	}

	// Iterate over each vertex as the source.
	for _, sourceVertex := range vertices {
		s := g.Hash()(sourceVertex.Value())

		// Stack S
		S := queue.NewStack[K]()

		// Predecessors
		P := make(map[K][]K, len(vertices))

		// Distance from source
		dist := make(map[K]float64, len(vertices))

		// Number of shortest paths
		sigma := make(map[K]float64, len(vertices))

		// Initialize distances and sigma
		for _, v := range vertices {
			k := g.Hash()(v.Value())
			dist[k] = math.Inf(1) // Initialize to infinity
			sigma[k] = 0.0
			P[k] = []K{}
		}
		dist[s] = 0.0
		sigma[s] = 1.0

		// Queue for BFS
		q := queue.NewQueue(queue.WithMode(queue.ModeFIFO))
		q.Enqueue(s)

		for !q.IsEmpty() {
			vAny, _ := q.Dequeue()
			v := vAny.(K)
			S.Push(v)

			for wK := range adj[v] { // adj[v] is map[K]Edge[K]
				if dist[wK] == math.Inf(1) {
					dist[wK] = dist[v] + 1.0
					q.Enqueue(wK)
				}
				if math.Abs(dist[wK]-(dist[v]+1.0)) < 1e-14 {
					sigma[wK] += sigma[v]
					P[wK] = append(P[wK], v)
				}
			}
		}

		// Accumulation phase
		delta := make(map[K]float64, len(vertices))
		for _, v := range vertices {
			k := g.Hash()(v.Value())
			delta[k] = 0.0
		}

		for !S.IsEmpty() {
			w, _ := S.Pop()
			for _, v := range P[w] {
				if sigma[w] != 0.0 {
					delta[v] += (sigma[v] / sigma[w]) * (1.0 + delta[w])
				}
			}
			if w != s {
				BC[w] += delta[w]
			}
		}
	}

	// If the graph is undirected, divide centralities by 2.0
	if !g.Traits().IsDirected {
		for k := range BC {
			BC[k] /= 2.0
		}
	}

	return BC, nil
}

// brandesWeighted computes betweenness centrality for weighted graphs.
// It returns an error if any edge has a negative weight.
func brandesWeighted[K graph.Ordered, T any](g graph.Interface[K, T]) (map[K]float64, error) {
	// Fetch adjacency for quick neighbor lookups
	adj, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("brandes: failed to retrieve adjacency map: %w", err)
	}

	// Quick negative-weight check
	for _, neighbors := range adj {
		for _, edge := range neighbors {
			if edge.Properties().Weight() < 0 {
				return nil, fmt.Errorf("brandes: negative-weight edge detected; not supported")
			}
		}
	}

	// Get all vertices in the graph.
	vertices, err := g.Vertices()
	if err != nil {
		return nil, fmt.Errorf("brandes: failed to retrieve vertices: %w", err)
	}

	// Initialize betweenness centrality scores to 0.
	BC := make(map[K]float64, len(vertices))
	for _, v := range vertices {
		k := g.Hash()(v.Value()) // Derive key using the hash function
		BC[k] = 0.0
	}

	// Iterate over each vertex as the source.
	for _, sourceVertex := range vertices {
		s := g.Hash()(sourceVertex.Value())

		// Stack S
		S := queue.NewStack[K]()

		// Predecessors
		P := make(map[K][]K, len(vertices))

		// Distance from source
		dist := make(map[K]float64, len(vertices))

		// Number of shortest paths
		sigma := make(map[K]float64, len(vertices))

		// Initialize distances and sigma
		for _, v := range vertices {
			k := g.Hash()(v.Value())
			dist[k] = math.Inf(1) // Initialize to infinity
			sigma[k] = 0.0
			P[k] = []K{}
		}
		dist[s] = 0.0
		sigma[s] = 1.0

		// Priority Queue for Dijkstra-like traversal
		pq := queue.NewPriorityQueue[K]()
		pq.Enqueue(s, 0.0)

		for pq.Len() > 0 {
			v, errPQ := pq.Dequeue()
			if errPQ != nil {
				return nil, fmt.Errorf("brandes: priority queue dequeue error: %w", errPQ)
			}
			S.Push(v)

			for wK, e := range adj[v] { // adj[v] is map[K]Edge[K]
				weight := e.Properties().Weight()
				alt := dist[v] + weight

				if alt < dist[wK] {
					dist[wK] = alt
					sigma[wK] = 0.0
					P[wK] = []K{}
					pq.SetPriority(wK, alt)
					pq.Enqueue(wK, alt)
				}
				if math.Abs(alt-dist[wK]) < 1e-14 {
					sigma[wK] += sigma[v]
					P[wK] = append(P[wK], v)
				}
			}
		}

		// Accumulation phase
		delta := make(map[K]float64, len(vertices))
		for _, v := range vertices {
			k := g.Hash()(v.Value())
			delta[k] = 0.0
		}

		for !S.IsEmpty() {
			w, _ := S.Pop()
			for _, v := range P[w] {
				if sigma[w] != 0.0 {
					delta[v] += (sigma[v] / sigma[w]) * (1.0 + delta[w])
				}
			}
			if w != s {
				BC[w] += delta[w]
			}
		}
	}

	// If the graph is undirected, divide centralities by 2.0
	if !g.Traits().IsDirected {
		for k := range BC {
			BC[k] /= 2.0
		}
	}

	return BC, nil
}
