package main

import (
	"os"

	"github.com/shghadge/todo_mcp/internal/mcp"
	"github.com/shghadge/todo_mcp/internal/storage"
)

func main() {
	// Initialize file-based storage
	todoStorage := storage.NewFileStorage("todos.json")

	// Create MCP server
	server := mcp.NewMCPServer(todoStorage)

	// Process input/output via stdio
	server.ProcessInput(os.Stdin, os.Stdout)
}
