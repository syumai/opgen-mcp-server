// SPDX-FileCopyrightText: 2025 syumai
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
	opgenserver "github.com/syumai/opgen-mcp-server/server"
)

func main() {
	srv := opgenserver.NewServer()
	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
}
