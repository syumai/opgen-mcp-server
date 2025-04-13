// SPDX-FileCopyrightText: 2025 syumai
// SPDX-License-Identifier: Apache-2.0
package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.1password.io/spg"
	"golang.design/x/clipboard"
)

func NewServer() *server.MCPServer {
	server := server.NewMCPServer(
		"opgen",
		"0.0.2",
		server.WithLogging(),
		server.WithRecovery(),
	)

	opgenCharactersTool := mcp.NewTool("generate_password_characters",
		mcp.WithDescription("Generate secure passwords using characters and copy to the clipboard. This is the general way to generate passwords."),
		mcp.WithNumber("length",
			mcp.Description("Password length for characters recipe"),
			mcp.DefaultNumber(float64(defaultCharRecipe.length)),
		),
		mcp.WithArray("allow",
			mcp.Description("Allowed character classes"),
			mcp.Items(map[string]any{
				"enum": []string{"uppercase", "lowercase", "digits", "symbols", "ambiguous"},
			}),
			mcp.AdditionalProperties(map[string]any{
				"default": defaultCharRecipe.allow,
			}),
		),
		mcp.WithArray("require",
			mcp.Description("Required character classes"),
			mcp.Items(map[string]any{
				"enum": []string{"uppercase", "lowercase", "digits", "symbols", "ambiguous"},
			}),
			mcp.AdditionalProperties(map[string]any{
				"default": defaultCharRecipe.require,
			}),
		),
		mcp.WithArray("exclude",
			mcp.Description("Excluded character classes"),
			mcp.Items(map[string]any{
				"enum": []string{"uppercase", "lowercase", "digits", "symbols", "ambiguous"},
			}),
			mcp.AdditionalProperties(map[string]any{
				"default": defaultCharRecipe.exclude,
			}),
		),
		// opgen MCP server does not support entropy
		// mcp.WithBoolean("entropy",
		// 	mcp.Description("Show the entropy of the password recipe"),
		// 	mcp.DefaultBool(false),
		// ),
	)
	server.AddTool(opgenCharactersTool, charGeneratorHandler)

	opgenWordsTool := mcp.NewTool("generate_password_words",
		mcp.WithDescription("Generate secure passwords using words and copy to the clipboard."),
		mcp.WithNumber("size",
			mcp.Description("Number of elements for words recipe"),
			mcp.DefaultNumber(4),
		),
		mcp.WithString("list",
			mcp.Description("Wordlist to use"),
			mcp.Enum("words", "syllables"),
			mcp.DefaultString("words"),
		),
		// opgen MCP server does not support wordlist file
		// mcp.WithString("file",
		// 	mcp.Description("Wordlist file to use"),
		// ),
		mcp.WithString("separator",
			mcp.Description("Word separator type"),
			mcp.Enum("hyphen", "space", "comma", "period", "underscore", "digit", "none"),
			mcp.DefaultString("hyphen"),
		),
		mcp.WithString("capitalize",
			mcp.Description("Capitalization scheme"),
			mcp.Enum("none", "first", "all", "random", "one"),
			mcp.DefaultString("none"),
		),
		// opgen MCP server does not support entropy
		// mcp.WithBoolean("entropy",
		// 	mcp.Description("Show the entropy of the password recipe"),
		// 	mcp.DefaultBool(false),
		// ),
	)
	server.AddTool(opgenWordsTool, wlGeneratorHandler)

	return server
}

func copyPasswordToClipboard(pwd *spg.Password) (*mcp.CallToolResult, error) {
	if err := clipboard.Init(); err != nil {
		return nil, err
	}
	clipboard.Write(clipboard.FmtText, []byte(pwd.String()))
	return mcp.NewToolResultText("Password generated and copied to clipboard"), nil
}

func convertCharacterClasses(classes []string, defaults []string) spg.CTFlag {
	var ccFlags spg.CTFlag
	if len(classes) == 0 {
		classes = defaults
	}

	for _, c := range classes {
		if ccFlag, ok := ccMap[c]; ok {
			ccFlags |= ccFlag
		}
	}

	return ccFlags
}

func slicesMap[T any, R any](s []T, fn func(item T) R) []R {
	result := make([]R, len(s))
	for i, item := range s {
		result[i] = fn(item)
	}
	return result
}

func assertAnyToString(s any) string {
	return s.(string)
}

func charGeneratorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	length := request.Params.Arguments["length"].(float64)
	var allow []string
	if allowAny, ok := request.Params.Arguments["allow"].([]any); ok {
		allow = slicesMap(allowAny, assertAnyToString)
	}
	var require []string
	if requireAny, ok := request.Params.Arguments["require"].([]any); ok {
		require = slicesMap(requireAny, assertAnyToString)
	}
	var exclude []string
	if excludeAny, ok := request.Params.Arguments["exclude"].([]any); ok {
		exclude = slicesMap(excludeAny, assertAnyToString)
	}

	recipe := spg.NewCharRecipe(int(length))
	recipe.Allow = convertCharacterClasses(allow, defaultCharRecipe.allow)
	recipe.Require = convertCharacterClasses(require, defaultCharRecipe.require)
	recipe.Exclude = convertCharacterClasses(exclude, defaultCharRecipe.exclude)

	pwd, err := recipe.Generate()
	if err != nil {
		return nil, err
	}

	return copyPasswordToClipboard(pwd)
}

func convertWordList(list string) (*spg.WordList, error) {
	var words []string
	switch list {
	case "words":
		words = spg.AgileWords
	case "syllables":
		words = spg.AgileSyllables
	default:
		return nil, fmt.Errorf("invalid wordlist: %s", list)
	}

	wordList, err := spg.NewWordList(words)
	if err != nil {
		return nil, fmt.Errorf("error setting up wordlist: %v", err)
	}
	return wordList, nil
}

func wlGeneratorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	size := request.Params.Arguments["size"].(float64)
	list := request.Params.Arguments["list"].(string)
	separator := request.Params.Arguments["separator"].(string)
	capitalize := request.Params.Arguments["capitalize"].(string)

	wordList, err := convertWordList(list)
	if err != nil {
		return nil, err
	}

	recipe := spg.NewWLRecipe(int(size), wordList)
	recipe.SeparatorFunc = separatorMap[separator]
	recipe.Capitalize = capitalizeMap[capitalize]

	pwd, err := recipe.Generate()
	if err != nil {
		return nil, err
	}

	return copyPasswordToClipboard(pwd)
}
