// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package paths

import (
	"fmt"
	"math"

	"github.com/sixafter/graph"
	"github.com/sixafter/graph/internal/queue"
)

// DijkstraFrom computes the shortest path between a source and a target vertex
// considering the edge weights. It returns a slice of hash values of the vertices
// forming that path, including the source and target.
//
// If the target is not reachable from the source, ErrTargetNotReachable is returned.
// If there are multiple shortest paths, an arbitrary one will be returned.
//
// Example:
//
//	path, err := DijkstraFrom(graph, "A", "D")
//	if err != nil {
//		if errors.Is(err, ErrTargetNotReachable) {
//			fmt.Println("Target not reachable from source")
//		} else {
//			log.Fatal(err)
//		}
//	} else {
//		fmt.Printf("Shortest path: %v\n", path)
//	}
func DijkstraFrom[K graph.Ordered, T any](g graph.Interface[K, T], source, target K) ([]K, error) {
	if g == nil {
		return nil, graph.ErrNilInputGraph
	}

	if source == target {
		return []K{source}, nil
	}

	exists, err := g.HasVertex(source)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, graph.ErrVertexNotFound
	}

	exists, err = g.HasVertex(target)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, graph.ErrVertexNotFound
	}

	weights := make(map[K]float64)
	visited := make(map[K]bool)

	weights[source] = 0
	visited[target] = true

	q := queue.NewPriorityQueue[K]()
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", graph.ErrAdjacencyMap, err)
	}

	for hash := range adjacencyMap {
		if hash != source {
			weights[hash] = math.Inf(1)
			visited[hash] = false
		}

		q.Enqueue(hash, weights[hash])
	}

	bestPredecessors := make(map[K]K)

	for q.Len() > 0 {
		vertex, _ := q.Dequeue()
		hasInfiniteWeight := math.IsInf(weights[vertex], 1)

		for adjacency, edge := range adjacencyMap[vertex] {
			edgeWeight := edge.Properties().Weight()

			if !g.Traits().IsWeighted {
				edgeWeight = 1
			}

			weight := weights[vertex] + edgeWeight

			if weight < weights[adjacency] && !hasInfiniteWeight {
				weights[adjacency] = weight
				bestPredecessors[adjacency] = vertex
				q.SetPriority(adjacency, weight)
			}
		}
	}

	path := []K{target}
	current := target

	for current != source {
		if _, ok := bestPredecessors[current]; !ok {
			return nil, graph.ErrTargetNotReachable
		}
		current = bestPredecessors[current]
		path = append([]K{current}, path...)
	}

	return path, nil
}
