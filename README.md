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

- **Generics-Based Design**: Leverages Go generics for a flexible and type-safe graph interface supporting custom vertex and edge types.
- **Trait-Driven Configuration**: Supports traits such as directed/undirected, weighted, acyclic, rooted, and multigraph properties.
- **Comprehensive Graph Operations**: Provides efficient algorithms for CRUD operations, set operations, adjacency/predecessor maps, and graph cloning.
- **Traversal and Pathfinding**: Implements breadth-first, depth-first, and shortest-path algorithms. Supports minimum and maximum spanning tree computation.
- **Graph Metrics**: Offers centrality measures (degree, closeness, betweenness, eigenvector), clustering coefficients, density, diameter, and average path length.
- **Community and Ranking Analysis**: Includes modularity and PageRank calculations for advanced graph analysis.
- **Cycle Management**: Prevents cycles in acyclic graphs during edge additions.
- **Streaming Support**: Enables paginated streaming of vertices and edges with context management for cancellation and resumption.
- **Customizable Input/Output**: Supports flexible graph serialization and custom reader/writer implementations.
- **Concurrency Safe**: Designed for thread-safe operations in multi-threaded environments.
- **Lightweight and Efficient**: Optimized for high performance with minimal overhead.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library.
- **Supports `io.Reader` Interface**:
  - Graph Serialization: Export and import graphs to/from various formats for interoperability. 
  - Customizable Readers and Writers: Create tailored I/O operations for graph persistence.

---

## Installation

### Using `go get`

To install this package, run the following command:

```sh
go get -u github.com/sixafter/graph
```

To use the package in your Go project, import it as follows:

```go
import "github.com/sixafter/graph"
```
---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
