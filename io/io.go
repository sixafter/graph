// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package io

import (
	"io"

	"github.com/sixafter/graph"
)

// Reader defines methods for reading a graph from a data source.
type Reader[K graph.Ordered, T any] interface {
	ReadGraph(r io.Reader, g graph.Interface[K, T]) error
}

// Writer defines methods for writing a graph to a data destination.
type Writer[K graph.Ordered, T any] interface {
	WriteGraph(w io.Writer, g graph.Interface[K, T]) error
}

// ReaderWriter defines methods for both reading and writing a graph.
type ReaderWriter[K graph.Ordered, T any] interface {
	Reader[K, T]
	Writer[K, T]
}
