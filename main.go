package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SpecKitServer struct {
	server     *mcp.Server
	workingDir string
}

func NewSpecKitServer(workingDir string) *SpecKitServer {
	impl := &mcp.Implementation{
		Name:    "spec-kit",
		Title:   "Spec-Kit MCP Server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(impl, nil)
	
	sks := &SpecKitServer{
		server:     server,
		workingDir: workingDir,
	}
	
	// Add tools
	sks.addTools()
	
	return sks
}

func (s *SpecKitServer) addTools() {
	// Add init_project tool
	initTool := &mcp.Tool{
		Name:        "init_project",
		Description: "Initialize a new spec-kit project",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the project to initialize",
				},
				"ai_assistant": map[string]interface{}{
					"type": "string",
					"enum": []string{"claude", "gemini", "copilot", "cursor"},
				},
			},
			"required": []string{"project_name"},
		},
	}
	
	s.server.AddTool(initTool, s.handleInitProject)

	// Add specify tool
	specifyTool := &mcp.Tool{
		Name:        "specify",
		Description: "Create feature specification from natural language",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"feature_description": map[string]interface{}{
					"type":        "string",
					"description": "Natural language description of the feature",
				},
			},
			"required": []string{"feature_description"},
		},
	}
	
	s.server.AddTool(specifyTool, s.handleSpecify)

	// Add plan tool
	planTool := &mcp.Tool{
		Name:        "plan",
		Description: "Generate implementation plan from specification",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"spec_file": map[string]interface{}{
					"type":        "string",
					"description": "Path to the specification file",
				},
			},
			"required": []string{"spec_file"},
		},
	}
	
	s.server.AddTool(planTool, s.handlePlan)

	// Add implement tool
	implementTool := &mcp.Tool{
		Name:        "implement",
		Description: "Generate code from specification",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"spec_file": map[string]interface{}{
					"type":        "string",
					"description": "Path to the specification file",
				},
				"output_dir": map[string]interface{}{
					"type":        "string",
					"description": "Directory to output generated code",
				},
			},
			"required": []string{"spec_file"},
		},
	}
	
	s.server.AddTool(implementTool, s.handleImplement)

	// Add analyze tool
	analyzeTool := &mcp.Tool{
		Name:        "analyze",
		Description: "Analyze project state and provide insights",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the project directory",
				},
			},
			"required": []string{"project_path"},
		},
	}
	
	s.server.AddTool(analyzeTool, s.handleAnalyze)

	// Add tasks tool
	tasksTool := &mcp.Tool{
		Name:        "tasks",
		Description: "Break down work into actionable tasks",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"spec_file": map[string]interface{}{
					"type":        "string",
					"description": "Path to the specification file",
				},
			},
			"required": []string{"spec_file"},
		},
	}
	
	s.server.AddTool(tasksTool, s.handleTasks)
}

func (s *SpecKitServer) handleInitProject(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ProjectName  string `json:"project_name"`
		AIAssistant  string `json:"ai_assistant"`
	}
	
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
			},
		}, nil
	}

	if args.ProjectName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: project_name is required"},
			},
		}, nil
	}

	if args.AIAssistant == "" {
		args.AIAssistant = "claude"
	}

	// Create project directory
	projectPath := filepath.Join(s.workingDir, args.ProjectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error creating project directory: %v", err)},
			},
		}, nil
	}

	// Run spec-kit init command
	cmd := exec.Command("specify", "init", args.ProjectName, "--ai-assistant", args.AIAssistant)
	cmd.Dir = s.workingDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error running spec-kit init: %s", string(output))},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("✅ Initialized spec-kit project '%s' with %s assistant\n\n%s", 
				args.ProjectName, args.AIAssistant, string(output))},
		},
	}, nil
}

func (s *SpecKitServer) handleSpecify(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		FeatureDescription string `json:"feature_description"`
	}
	
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
			},
		}, nil
	}

	if args.FeatureDescription == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: feature_description is required"},
			},
		}, nil
	}

	// Create a temporary spec file
	specFile := filepath.Join(s.workingDir, "temp_spec.md")
	
	// Run spec-kit specify command
	cmd := exec.Command("specify", "specify", args.FeatureDescription)
	cmd.Dir = s.workingDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error running spec-kit specify: %s", string(output))},
			},
		}, nil
	}

	// Save the output to a file
	if err := os.WriteFile(specFile, output, 0644); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error saving spec file: %v", err)},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("✅ Created specification from description\n\n**Specification:**\n%s\n\n**Saved to:** %s", 
				string(output), specFile)},
		},
	}, nil
}

func (s *SpecKitServer) handlePlan(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SpecFile string `json:"spec_file"`
	}
	
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
			},
		}, nil
	}

	if args.SpecFile == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: spec_file is required"},
			},
		}, nil
	}

	// Run spec-kit plan command
	cmd := exec.Command("specify", "plan", args.SpecFile)
	cmd.Dir = s.workingDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error running spec-kit plan: %s", string(output))},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("✅ Generated implementation plan\n\n**Plan:**\n%s", string(output))},
		},
	}, nil
}

func (s *SpecKitServer) handleImplement(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SpecFile  string `json:"spec_file"`
		OutputDir string `json:"output_dir"`
	}
	
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
			},
		}, nil
	}

	if args.SpecFile == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: spec_file is required"},
			},
		}, nil
	}

	if args.OutputDir == "" {
		args.OutputDir = filepath.Join(s.workingDir, "generated")
	}

	// Run spec-kit implement command
	cmd := exec.Command("specify", "implement", args.SpecFile, "--output", args.OutputDir)
	cmd.Dir = s.workingDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error running spec-kit implement: %s", string(output))},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("✅ Generated code from specification\n\n**Output:**\n%s\n\n**Generated files in:** %s", 
				string(output), args.OutputDir)},
		},
	}, nil
}

func (s *SpecKitServer) handleAnalyze(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ProjectPath string `json:"project_path"`
	}
	
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
			},
		}, nil
	}

	if args.ProjectPath == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: project_path is required"},
			},
		}, nil
	}

	// Run spec-kit analyze command
	cmd := exec.Command("specify", "analyze", args.ProjectPath)
	cmd.Dir = s.workingDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error running spec-kit analyze: %s", string(output))},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("✅ Analyzed project state\n\n**Analysis:**\n%s", string(output))},
		},
	}, nil
}

func (s *SpecKitServer) handleTasks(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SpecFile string `json:"spec_file"`
	}
	
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error parsing arguments: %v", err)},
			},
		}, nil
	}

	if args.SpecFile == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: spec_file is required"},
			},
		}, nil
	}

	// Run spec-kit tasks command
	cmd := exec.Command("specify", "tasks", args.SpecFile)
	cmd.Dir = s.workingDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error running spec-kit tasks: %s", string(output))},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("✅ Generated actionable tasks\n\n**Tasks:**\n%s", string(output))},
		},
	}, nil
}

func main() {
	// Get working directory from environment or use current directory
	workingDir := os.Getenv("SPEC_KIT_WORKING_DIR")
	if workingDir == "" {
		workingDir = "."
	}

	server := NewSpecKitServer(workingDir)
	
	log.SetFlags(0)
	log.Printf("Starting spec-kit MCP server in directory: %s", workingDir)
	
	// Run the server using stdio transport
	ctx := context.Background()
	if err := server.server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}