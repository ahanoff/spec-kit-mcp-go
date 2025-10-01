package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func testServer() {
	// Create server
	server := NewSpecKitServer(".")
	
	// Create in-memory transport for testing
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	
	// Start server in a goroutine
	go func() {
		ctx := context.Background()
		if err := server.server.Run(ctx, serverTransport); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Create client to test the server
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	ctx := context.Background()
	
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		log.Fatalf("Client connection error: %v", err)
	}
	
	// Test listing tools
	tools, err := clientSession.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		log.Fatalf("List tools error: %v", err)
	}
	
	log.Printf("Available tools: %d", len(tools.Tools))
	for _, tool := range tools.Tools {
		log.Printf("- %s: %s", tool.Name, tool.Description)
	}
	
	// Test calling a tool (this will fail because spec-kit is not installed, but we can see the error handling)
	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "init_project",
		Arguments: map[string]any{
			"project_name": "test-project",
			"ai_assistant": "claude",
		},
	})
	if err != nil {
		log.Printf("Tool call error: %v", err)
	} else {
		log.Printf("Tool call result: %+v", result)
	}
}
