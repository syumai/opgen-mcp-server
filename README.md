# opgen MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server implementation for password generation, based on [1Password/spg/cmd/opgen](https://github.com/1Password/spg/tree/master/cmd/opgen).

## Installation

```bash
go install github.com/syumai/opgen-mcp-server@latest
```

## Usage

```jsonc
{
    // If you have installed opgen-mcp-server
    "opgen": {
        "command": "opgen-mcp-server",
    },
    // Alternative configuration
    "opgen": {
        "command": "go",
        "args": ["run", "github.com/syumai/opgen-mcp-server@latest"],
    },
}
```

## Tools

### `generate_password_characters`

#### Configurable options

- Password length
- Allowed character sets
- Required character sets
- Excluded character sets

#### Character sets

- uppercase
- lowercase
- digits
- symbols
- ambiguous

### `generate_password_words`

#### Configurable options

- Number of words
- Word list selection:
  - words
  - syllables
- Word separator type:
  - hyphen
  - space
  - comma
  - period
  - underscore
  - digit
  - none
- Capitalization scheme:
  - none
  - first
  - all
  - random
  - one
