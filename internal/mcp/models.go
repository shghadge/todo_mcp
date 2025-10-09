package mcp

import (
	"encoding/json"
	"time"
)

// MCP Protocol Types

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// InitializeRequest represents the initialize request
type InitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

// InitializeResponse represents the initialize response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ClientCapabilities represents client capabilities
type ClientCapabilities struct {
	Sampling map[string]interface{} `json:"sampling,omitempty"`
}

// ServerCapabilities represents server capabilities
type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Logging   *LoggingCapability   `json:"logging,omitempty"`
}

// ToolsCapability represents tools capability
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability represents resources capability
type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// LoggingCapability represents logging capability
type LoggingCapability struct{}

// ClientInfo represents client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo represents server information
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ListToolsRequest represents the tools/list request
type ListToolsRequest struct{}

// ListToolsResponse represents the tools/list response
type ListToolsResponse struct {
	Tools []Tool `json:"tools"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolRequest represents the tools/call request
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResponse represents the tools/call response
type CallToolResponse struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// ListResourcesRequest represents the resources/list request
type ListResourcesRequest struct{}

// ListResourcesResponse represents the resources/list response
type ListResourcesResponse struct {
	Resources []Resource `json:"resources"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ReadResourceRequest represents the resources/read request
type ReadResourceRequest struct {
	URI string `json:"uri"`
}

// ReadResourceResponse represents the resources/read response
type ReadResourceResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent represents content of a resource
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

// Content represents content in a response
type Content struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// Error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// Method names
const (
	MethodInitialize     = "initialize"
	MethodInitialized    = "initialized"
	MethodListTools      = "tools/list"
	MethodCallTool       = "tools/call"
	MethodListResources  = "resources/list"
	MethodReadResource   = "resources/read"
	MethodPing           = "ping"
	MethodLoggingMessage = "logging/message"
)

// Tool names for our todo application
const (
	ToolCreateTodo = "create_todo"
	ToolGetTodo    = "get_todo"
	ToolGetTodos   = "get_todos"
	ToolUpdateTodo = "update_todo"
	ToolDeleteTodo = "delete_todo"
)

// Resource URIs for our todo application
const (
	ResourceTodosList      = "todo://todos"
	ResourceTodosPending   = "todo://todos/pending"
	ResourceTodosCompleted = "todo://todos/completed"
)

// Todo-specific request/response types for tools

// CreateTodoRequest represents parameters for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// GetTodoRequest represents parameters for getting a todo
type GetTodoRequest struct {
	ID int `json:"id"`
}

// GetTodosRequest represents parameters for getting todos
type GetTodosRequest struct {
	Status string `json:"status,omitempty"` // "pending", "completed", or empty for all
}

// UpdateTodoRequest represents parameters for updating a todo
type UpdateTodoRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"` // "pending" or "completed"
}

// DeleteTodoRequest represents parameters for deleting a todo
type DeleteTodoRequest struct {
	ID int `json:"id"`
}

// TodoResponse represents a todo in responses
type TodoResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoListResponse represents a list of todos
type TodoListResponse struct {
	Todos []TodoResponse `json:"todos"`
	Count int            `json:"count"`
}
