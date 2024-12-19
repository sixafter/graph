// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"github.com/sixafter/graph"
)

func IsSuperset[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	// Check if h is a subset of g
	return IsSubset(h, g)
}
