package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/shghadge/todo_mcp/internal/models"
)

// handleTodosListResource handles the todos list resource
func (s *MCPServer) handleTodosListResource() (*ReadResourceResponse, error) {
	todos, err := s.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving todos: %w", err)
	}

	// Convert to response format
	todoListResp := TodoListResponse{
		Todos: make([]TodoResponse, len(todos)),
		Count: len(todos),
	}

	for i, todo := range todos {
		todoListResp.Todos[i] = TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      string(todo.Status),
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	result, err := json.MarshalIndent(todoListResp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling todos: %w", err)
	}

	return &ReadResourceResponse{
		Contents: []ResourceContent{
			{
				URI:      ResourceTodosList,
				MimeType: "application/json",
				Text:     string(result),
			},
		},
	}, nil
}

// handleTodosPendingResource handles the pending todos resource
func (s *MCPServer) handleTodosPendingResource() (*ReadResourceResponse, error) {
	todos, err := s.storage.GetByStatus(models.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("error retrieving pending todos: %w", err)
	}

	// Convert to response format
	todoListResp := TodoListResponse{
		Todos: make([]TodoResponse, len(todos)),
		Count: len(todos),
	}

	for i, todo := range todos {
		todoListResp.Todos[i] = TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      string(todo.Status),
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	result, err := json.MarshalIndent(todoListResp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling pending todos: %w", err)
	}

	return &ReadResourceResponse{
		Contents: []ResourceContent{
			{
				URI:      ResourceTodosPending,
				MimeType: "application/json",
				Text:     string(result),
			},
		},
	}, nil
}

// handleTodosCompletedResource handles the completed todos resource
func (s *MCPServer) handleTodosCompletedResource() (*ReadResourceResponse, error) {
	todos, err := s.storage.GetByStatus(models.StatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("error retrieving completed todos: %w", err)
	}

	// Convert to response format
	todoListResp := TodoListResponse{
		Todos: make([]TodoResponse, len(todos)),
		Count: len(todos),
	}

	for i, todo := range todos {
		todoListResp.Todos[i] = TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      string(todo.Status),
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	result, err := json.MarshalIndent(todoListResp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling completed todos: %w", err)
	}

	return &ReadResourceResponse{
		Contents: []ResourceContent{
			{
				URI:      ResourceTodosCompleted,
				MimeType: "application/json",
				Text:     string(result),
			},
		},
	}, nil
}
