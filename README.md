# Spec-Kit MCP Server (Go)

A Model Context Protocol (MCP) server written in Go that bridges Amazon Q Developer to Microsoft's Spec-Kit, enabling spec-driven development workflows through natural language interaction.

## Features

- **init_project**: Initialize new spec-kit projects
- **specify**: Create feature specifications from natural language
- **plan**: Generate implementation plans from specifications
- **implement**: Generate code from specifications
- **analyze**: Analyze project state and provide insights
- **tasks**: Break down work into actionable tasks

## Prerequisites

- Go 1.21 or later
- Python 3.8 or later
- uv (Python package manager)
- Amazon Q Developer CLI
- Spec-Kit installed globally: `uv tool install specify-cli --from git+https://github.com/github/spec-kit.git`

## Installation

1. Clone this repository:
```bash
git clone git@github.com:ahanoff/spec-kit-mcp-go.git
cd spec-kit-mcp-go
```

2. Build the MCP server:
```bash
chmod +x build.sh
./build.sh
```

3. Install spec-kit globally:
```bash
uv tool install specify-cli --from git+https://github.com/github/spec-kit.git
```

## Configuration

Add the MCP server to your Q Developer configuration at `~/.aws/amazonq/mcp.json`:

```json
{
  "mcpServers": {
    "spec-kit": {
      "type": "stdio",
      "command": "/path/to/spec-kit-mcp-go/spec-kit-mcp-go",
      "env": {
        "SPEC_KIT_WORKING_DIR": "/path/to/your/projects"
      },
      "timeout": 120000
    }
  }
}
```

## Usage

Start Q Developer and use natural language to interact with spec-kit:

```
Use the init_project tool to create a new spec-kit project called "my-app" with Claude as the AI assistant
```

```
Use the specify tool to create a specification for user authentication with email and password
```

## Development

To modify the server:

1. Make your changes to `main.go`
2. Run `go mod tidy` to update dependencies
3. Rebuild with `./build.sh`
4. Restart Q Developer

## Troubleshooting

- Ensure the binary path in `mcp.json` is correct
- Check that spec-kit is installed: `specify --version`
- Verify Go is installed: `go version`
- Test the server manually: `echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./spec-kit-mcp-go`

## License

MIT
