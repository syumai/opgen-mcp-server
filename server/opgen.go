// SPDX-FileCopyrightText: 2025 syumai
// SPDX-FileCopyrightText: 2018 AgileBits, Inc.
// SPDX-License-Identifier: Apache-2.0

// This file is based on: https://github.com/1Password/spg/blob/v0.1.0/cmd/opgen/opgen.go
package server

import (
	"go.1password.io/spg"
)

var ccMap = map[string]spg.CTFlag{
	"uppercase": spg.Uppers,
	"lowercase": spg.Lowers,
	"digits":    spg.Digits,
	"symbols":   spg.Symbols,
	"ambiguous": spg.Ambiguous,
}

var separatorMap = map[string]spg.SFFunction{
	"hyphen":     createSeparatorFunc("-"),
	"space":      createSeparatorFunc(" "),
	"comma":      createSeparatorFunc(","),
	"period":     createSeparatorFunc("."),
	"underscore": createSeparatorFunc("_"),
	"digit":      spg.SFDigits1,
	"none":       spg.SFNone,
}

var capitalizeMap = map[string]spg.CapScheme{
	"none":   spg.CSNone,
	"first":  spg.CSFirst,
	"all":    spg.CSAll,
	"random": spg.CSRandom,
	"one":    spg.CSOne,
}

func createSeparatorFunc(value string) spg.SFFunction {
	return func() (string, spg.FloatE) {
		return value, 0
	}
}
