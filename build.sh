#!/bin/bash

# Build the MCP server
echo "Building spec-kit MCP server..."
go build -o spec-kit-mcp-go main.go

# Make it executable
chmod +x spec-kit-mcp-go

echo "✅ Build complete! Binary: ./spec-kit-mcp-go"
echo "Install spec-kit globally: npm install -g @github/spec-kit"
