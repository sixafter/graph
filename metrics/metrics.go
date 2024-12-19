// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package metrics

import (
	"github.com/sixafter/graph"
)

// Metrics defines methods for calculating various graph metrics.
// Each method accepts one or more Interface[K, T] instances as parameters.
type Metrics[K graph.Ordered, T any] interface {
	// DegreeCentrality calculates the degree centrality for each vertex in the given graph.
	DegreeCentrality(g graph.Interface[K, T]) (map[K]float64, error)

	// BetweennessCentrality calculates the betweenness centrality for each vertex in the given graph.
	BetweennessCentrality(g graph.Interface[K, T]) (map[K]float64, error)

	// ClusteringCoefficient calculates the clustering coefficient for each vertex in the given graph.
	ClusteringCoefficient(g graph.Interface[K, T]) (map[K]float64, error)

	// ClosenessCentrality calculates the closeness centrality for each vertex in the given graph.
	ClosenessCentrality(g graph.Interface[K, T]) (map[K]float64, error)

	// EigenvectorCentrality calculates the eigenvector centrality for each vertex in the given graph.
	EigenvectorCentrality(g graph.Interface[K, T]) (map[K]float64, error)

	// PageRank calculates the PageRank for each vertex in the given graph.
	PageRank(g graph.Interface[K, T], dampingFactor float64, maxIterations int, tol float64) (map[K]float64, error)

	// GlobalClusteringCoefficient calculates the global clustering coefficient of the given graph.
	GlobalClusteringCoefficient(g graph.Interface[K, T]) (float64, error)

	// Density calculates the density of the given graph.
	Density(g graph.Interface[K, T]) (float64, error)

	// Diameter calculates the diameter of the given graph.
	Diameter(g graph.Interface[K, T]) (int, error)

	// AveragePathLength calculates the average shortest path length in the given graph.
	AveragePathLength(g graph.Interface[K, T]) (float64, error)

	// Transitivity calculates the transitivity of the given graph.
	Transitivity(g graph.Interface[K, T]) (float64, error)

	// Modularity calculates the modularity of the given graph based on a provided community structure.
	Modularity(g graph.Interface[K, T], communities map[K]int) (float64, error)
}
