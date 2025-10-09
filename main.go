package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shghadge/todo_mcp/internal/handlers"
	"github.com/shghadge/todo_mcp/internal/storage"
)

func main() {
	fmt.Println("Starting Todo MCP Server...")

	// Initialize file-based storage
	todoStorage := storage.NewFileStorage("/home/shubham/projects/todo_mcp/todos.json")

	// Setup routes
	router := handlers.SetupRoutes(todoStorage)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, router))
}
