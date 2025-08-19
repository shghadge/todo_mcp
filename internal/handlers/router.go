package handlers

import (
	"fmt"
	"todo_mcp/internal/storage"

	"github.com/gorilla/mux"
)

// SetupRoutes sets up HTTP routes for the todo application
func SetupRoutes(storage storage.TodoStorage) *mux.Router {
	router := mux.NewRouter()

	// Create todo handler
	todoHandler := NewTodoHandler(storage)

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Todo routes
	api.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	api.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.GetTodo).Methods("GET")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.UpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.DeleteTodo).Methods("DELETE")

	PrintRoutes(router)
	return router
}

func PrintRoutes(router *mux.Router) {
	// Walk the router and print all registered routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		for _, method := range methods {
			fmt.Printf("  %-6s %s\n", method, pathTemplate)
		}
		return nil
	})
}
