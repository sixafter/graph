# graph

[![CI](https://github.com/sixafter/graph/workflows/ci/badge.svg)](https://github.com/sixafter/graph/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/graph)](https://img.shields.io/github/go-mod/go-version/sixafter/graph)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_graph&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_graph)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/graph)](https://github.com/sixafter/graph/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/graph.svg)](https://pkg.go.dev/github.com/sixafter/graph)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/graph)](https://goreportcard.com/report/github.com/sixafter/graph)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
![CodeQL](https://github.com/sixafter/graph/actions/workflows/codeql-analysis.yaml/badge.svg)

A Go library for creating and manipulating graph data structures.

## Features

This Go-based graph library is designed for versatility, performance, and extensibility, leveraging generics to handle various graph-related operations seamlessly. Key features include:

### Graph Interface
- **Generic Graph Model**: Define graphs with custom vertex and edge types, supported by Go's type constraints for added safety and flexibility.
- **Traits-Based Behavior**: Configure graph traits such as directionality, weighted edges, acyclic properties, multigraph capabilities, and rooted structures.
- **Basic Operations**: Add, retrieve, modify, and remove vertices and edges with ease.
- **Streaming Support**: Stream vertices and edges in paginated batches with context-aware operations.
- **Adjacency and Predecessor Maps**: Retrieve adjacency and predecessor maps for detailed graph structure insights.

### Core Functionalities
- **Directed and Undirected Graphs**: Handle both directed and undirected graphs with appropriate algorithms and constraints.
- **Set Operations**: Merge, intersect, or differentiate graphs using efficient set-based operations.
- **Graph Cloning**: Create deep copies of graphs while preserving all properties and relationships.

### Traversals and Paths
- **Traversal Algorithms**: Support for breadth-first and depth-first traversals.
- **Shortest Path**: Compute shortest paths using Dijkstra's algorithm or other customizable strategies.
- **Maximum and Minimum Spanning Trees**: Build spanning trees efficiently using Kruskal's algorithm for both weight extremes.

### Metrics and Analysis
- **Graph Metrics**:
    - Degree, Closeness, Betweenness, and Eigenvector Centralities.
    - Clustering Coefficients (local and global).
    - Graph Density, Diameter, and Average Path Length.
- **Community Detection**: Compute modularity based on a predefined community structure.
- **Transitivity and PageRank**: Analyze connectivity and rank vertices.

### Input/Output (I/O)
- **Graph Serialization**: Export and import graphs to/from various formats for interoperability.
- **Customizable Readers and Writers**: Create tailored I/O operations for graph persistence.

### Additional Utilities
- **Cycle Prevention**: Enforce acyclic properties during graph operations.
- **Weighted Operations**: Support for weighted graphs in metrics, traversal, and pathfinding.
- **Modularity**: Measure and manage community structures within graphs.

---

## Installation

### Using `go get`

To install this package, run the following command:

```sh
go get -u github.com/sixafter/graph
```

To use the this package in your Go project, import it as follows:

```go
import "github.com/sixafter/graph"
```
---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
