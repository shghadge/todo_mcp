package main

import (
	"fmt"
	"log"
	"net/http"
	"todo_mcp/internal/handlers"
	"todo_mcp/internal/storage"
)

func main() {
	fmt.Println("Starting Todo MCP Server...")

	// Initialize in-memory storage
	todoStorage := storage.NewInMemoryStorage()

	// Setup routes
	router := handlers.SetupRoutes(todoStorage)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("API endpoints:")

	log.Fatal(http.ListenAndServe(port, router))
}
