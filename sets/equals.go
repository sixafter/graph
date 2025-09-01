// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package sets

import (
	"github.com/sixafter/graph"
)

func Equals[K graph.Ordered, T any](g, h graph.Interface[K, T]) (bool, error) {
	if nil == g || nil == h {
		return false, nil
	}

	if !g.Traits().Equals(h.Traits()) {
		return false, graph.ErrGraphTypeMismatch
	}

	// Check if g is a subset of h
	subset, err := IsSubset(g, h)
	if err != nil {
		return false, err
	}
	if !subset {
		return false, nil
	}

	// Check if h is a subset of g
	subset, err = IsSubset(h, g)
	if err != nil {
		return false, err
	}
	if !subset {
		return false, nil
	}

	return true, nil
}
