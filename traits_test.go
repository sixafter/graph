// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTraits(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Default Traits
	traits := &Traits{}
	is.False(traits.IsAcyclic, "Default Traits should not be acyclic")
	is.False(traits.IsDirected, "Default Traits should not be directed")
	is.False(traits.IsRooted, "Default Traits should not be rooted")
	is.False(traits.IsWeighted, "Default Traits should not be weighted")
	is.False(traits.PreventCycles, "Default Traits should not prevent cycles")
}

func TestAcyclicOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	Acyclic()(traits)
	is.True(traits.IsAcyclic, "Traits should be acyclic after applying Acyclic()")
	is.False(traits.PreventCycles, "Acyclic() should not enable PreventCycles")
}

func TestDirectedOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	Directed()(traits)
	is.True(traits.IsDirected, "Traits should be directed after applying Directed()")
}

func TestMultiGraphOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	MultiGraph()(traits)
	is.True(traits.IsMultiGraph, "Traits should be a multigraph after applying MultiGraph()")
}

func TestPreventCyclesOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	PreventCycles()(traits)
	is.True(traits.IsAcyclic, "PreventCycles() should make the graph acyclic")
	is.True(traits.PreventCycles, "PreventCycles() should enable cycle prevention")
}

func TestRootedOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	Rooted()(traits)
	is.True(traits.IsRooted, "Traits should be rooted after applying Rooted()")
}

func TestTreeOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	Tree()(traits)
	is.True(traits.IsRooted, "Tree() should make the graph rooted")
	is.True(traits.IsAcyclic, "Tree() should make the graph acyclic")
}

func TestWeightedOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	Weighted()(traits)
	is.True(traits.IsWeighted, "Traits should be weighted after applying Weighted()")
}

func TestCombinationOfTraits(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits := &Traits{}
	Directed()(traits)
	Weighted()(traits)
	PreventCycles()(traits)

	is.True(traits.IsDirected, "Traits should be directed")
	is.True(traits.IsWeighted, "Traits should be weighted")
	is.True(traits.IsAcyclic, "PreventCycles() should make the graph acyclic")
	is.True(traits.PreventCycles, "Traits should prevent cycles")
	is.False(traits.IsRooted, "Traits should not be rooted by default")
}

func TestEquals(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	traits1 := &Traits{
		IsAcyclic:     true,
		IsDirected:    true,
		IsRooted:      true,
		IsWeighted:    true,
		PreventCycles: true,
	}
	traits2 := &Traits{
		IsAcyclic:     true,
		IsDirected:    true,
		IsRooted:      true,
		IsWeighted:    true,
		PreventCycles: true,
	}
	is.True(traits1.Equals(traits2), "Traits should be equal")

	traits1.IsWeighted = false
	is.False(traits1.Equals(traits2), "Traits should not be equal")
}
