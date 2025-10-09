package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/shghadge/todo_mcp/internal/storage"
)

// MCPServer represents the MCP server
type MCPServer struct {
	storage     storage.TodoStorage
	tools       map[string]ToolHandler
	resources   map[string]ResourceHandler
	initialized bool
}

// ToolHandler represents a function that handles tool calls
type ToolHandler func(args map[string]interface{}) (*CallToolResponse, error)

// ResourceHandler represents a function that handles resource requests
type ResourceHandler func() (*ReadResourceResponse, error)

// NewMCPServer creates a new MCP server
func NewMCPServer(storage storage.TodoStorage) *MCPServer {
	server := &MCPServer{
		storage:     storage,
		tools:       make(map[string]ToolHandler),
		resources:   make(map[string]ResourceHandler),
		initialized: false,
	}

	server.registerTools()
	server.registerResources()

	return server
}

// HandleRequest handles an incoming JSON-RPC request
func (s *MCPServer) HandleRequest(request *JSONRPCRequest) *JSONRPCResponse {
	var result interface{}
	var err error

	switch request.Method {
	case MethodInitialize:
		result, err = s.handleInitialize(request.Params)
	case MethodInitialized:
		// No response needed for initialized notification
		return nil
	case MethodListTools:
		result, err = s.handleListTools()
	case MethodCallTool:
		result, err = s.handleCallTool(request.Params)
	case MethodListResources:
		result, err = s.handleListResources()
	case MethodReadResource:
		result, err = s.handleReadResource(request.Params)
	case MethodPing:
		result = map[string]interface{}{"message": "pong"}
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    MethodNotFound,
				Message: fmt.Sprintf("Method not found: %s", request.Method),
			},
		}
	}

	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    InternalError,
				Message: err.Error(),
			},
		}
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// ProcessInput processes input from stdin (for stdio transport)
func (s *MCPServer) ProcessInput(input io.Reader, output io.Writer) {
	decoder := json.NewDecoder(input)
	encoder := json.NewEncoder(output)

	for {
		var request JSONRPCRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			// Send error response
			response := &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error: &JSONRPCError{
					Code:    ParseError,
					Message: "Parse error",
				},
			}
			encoder.Encode(response)
			continue
		}

		response := s.HandleRequest(&request)
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}
	}
}

// handleInitialize handles the initialize request
func (s *MCPServer) handleInitialize(params json.RawMessage) (*InitializeResponse, error) {
	var req InitializeRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid initialize request: %w", err)
	}

	s.initialized = true

	return &InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
			Logging: &LoggingCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "todo-mcp-server",
			Version: "1.0.0",
		},
	}, nil
}

// handleListTools handles the tools/list request
func (s *MCPServer) handleListTools() (*ListToolsResponse, error) {
	tools := []Tool{
		{
			Name:        ToolCreateTodo,
			Description: "Create a new todo item",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "The title of the todo item",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Optional description of the todo item",
					},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        ToolGetTodo,
			Description: "Get a specific todo item by ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "The ID of the todo item",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        ToolGetTodos,
			Description: "Get all todo items, optionally filtered by status",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Filter by status: 'pending' or 'completed' (optional)",
						"enum":        []string{"pending", "completed"},
					},
				},
			},
		},
		{
			Name:        ToolUpdateTodo,
			Description: "Update an existing todo item",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "The ID of the todo item to update",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "New title for the todo item",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New description for the todo item",
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "New status for the todo item",
						"enum":        []string{"pending", "completed"},
					},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        ToolDeleteTodo,
			Description: "Delete a todo item by ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "The ID of the todo item to delete",
					},
				},
				"required": []string{"id"},
			},
		},
	}

	return &ListToolsResponse{Tools: tools}, nil
}

// handleCallTool handles tool calls
func (s *MCPServer) handleCallTool(params json.RawMessage) (*CallToolResponse, error) {
	var req CallToolRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid call tool request: %w", err)
	}

	handler, exists := s.tools[req.Name]
	if !exists {
		return nil, fmt.Errorf("unknown tool: %s", req.Name)
	}

	return handler(req.Arguments)
}

// handleListResources handles the resources/list request
func (s *MCPServer) handleListResources() (*ListResourcesResponse, error) {
	resources := []Resource{
		{
			URI:         ResourceTodosList,
			Name:        "All Todos",
			Description: "Get all todo items",
			MimeType:    "application/json",
		},
		{
			URI:         ResourceTodosPending,
			Name:        "Pending Todos",
			Description: "Get all pending todo items",
			MimeType:    "application/json",
		},
		{
			URI:         ResourceTodosCompleted,
			Name:        "Completed Todos",
			Description: "Get all completed todo items",
			MimeType:    "application/json",
		},
	}

	return &ListResourcesResponse{Resources: resources}, nil
}

// handleReadResource handles resource read requests
func (s *MCPServer) handleReadResource(params json.RawMessage) (*ReadResourceResponse, error) {
	var req ReadResourceRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid read resource request: %w", err)
	}

	handler, exists := s.resources[req.URI]
	if !exists {
		return nil, fmt.Errorf("unknown resource: %s", req.URI)
	}

	return handler()
}

// registerTools registers all tool handlers
func (s *MCPServer) registerTools() {
	s.tools[ToolCreateTodo] = s.handleCreateTodo
	s.tools[ToolGetTodo] = s.handleGetTodo
	s.tools[ToolGetTodos] = s.handleGetTodos
	s.tools[ToolUpdateTodo] = s.handleUpdateTodo
	s.tools[ToolDeleteTodo] = s.handleDeleteTodo
}

// registerResources registers all resource handlers
func (s *MCPServer) registerResources() {
	s.resources[ResourceTodosList] = s.handleTodosListResource
	s.resources[ResourceTodosPending] = s.handleTodosPendingResource
	s.resources[ResourceTodosCompleted] = s.handleTodosCompletedResource
}
