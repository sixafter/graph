// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentityFunctionForStrings(t *testing.T) {
	t.Parallel()
	as := assert.New(t)

	input := "testVertex"
	expectedHash := "testVertex"

	hash := StringHash(input)

	as.Equal(expectedHash, hash, "StringHash should return the same string as the hash")
}

func TestDifferentStringsReturnDifferentHashes(t *testing.T) {
	t.Parallel()
	as := assert.New(t)

	input1 := "vertex1"
	input2 := "vertex2"

	hash1 := StringHash(input1)
	hash2 := StringHash(input2)

	as.NotEqual(hash1, hash2, "Different strings should return different hashes")
}

func TestEmptyStringHash(t *testing.T) {
	t.Parallel()
	as := assert.New(t)

	input := ""
	expectedHash := ""

	hash := StringHash(input)

	as.Equal(expectedHash, hash, "StringHash of an empty string should return an empty string")
}

func TestIdentityFunctionForIntegers(t *testing.T) {
	t.Parallel()
	as := assert.New(t)

	input := 42
	expectedHash := 42

	hash := IntHash(input)

	as.Equal(expectedHash, hash, "IntHash should return the same integer as the hash")
}

func TestDifferentIntegersReturnDifferentHashes(t *testing.T) {
	t.Parallel()
	as := assert.New(t)

	input1 := 1
	input2 := 2

	hash1 := IntHash(input1)
	hash2 := IntHash(input2)

	as.NotEqual(hash1, hash2, "Different integers should return different hashes")
}

func TestZeroHash(t *testing.T) {
	t.Parallel()
	as := assert.New(t)

	input := 0
	expectedHash := 0

	hash := IntHash(input)

	as.Equal(expectedHash, hash, "IntHash of zero should return zero")
}
