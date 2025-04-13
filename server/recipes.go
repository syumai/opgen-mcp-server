// SPDX-FileCopyrightText: 2025 syumai
// SPDX-FileCopyrightText: 2018 AgileBits, Inc.
// SPDX-License-Identifier: Apache-2.0

// This file is based on: https://github.com/1Password/spg/blob/v0.1.0/cmd/opgen/recipes.go
package server

type charRecipe struct {
	length  int
	allow   []string
	require []string
	exclude []string
}

var defaultCharRecipe = charRecipe{
	length:  20,
	allow:   []string{"uppercase", "lowercase", "digits", "symbols"},
	exclude: []string{"ambiguous"},
}
