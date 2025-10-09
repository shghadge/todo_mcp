package mcp

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/shghadge/todo_mcp/internal/models"
	"github.com/shghadge/todo_mcp/internal/storage"
)

// handleCreateTodo handles the create_todo tool
func (s *MCPServer) handleCreateTodo(args map[string]interface{}) (*CallToolResponse, error) {
	// Extract and validate arguments
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: title is required and must be a non-empty string",
			}},
			IsError: true,
		}, nil
	}

	description := ""
	if desc, ok := args["description"].(string); ok {
		description = desc
	}

	// Create todo
	todo := &models.Todo{
		Title:       title,
		Description: description,
		Status:      models.StatusPending,
	}

	if err := s.storage.Create(todo); err != nil {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error creating todo: %v", err),
			}},
			IsError: true,
		}, nil
	}

	// Convert to response format
	todoResp := TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      string(todo.Status),
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	result, _ := json.MarshalIndent(todoResp, "", "  ")
	return &CallToolResponse{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Todo created successfully:\n%s", string(result)),
		}},
	}, nil
}

// handleGetTodo handles the get_todo tool
func (s *MCPServer) handleGetTodo(args map[string]interface{}) (*CallToolResponse, error) {
	// Extract and validate ID
	idInterface, ok := args["id"]
	if !ok {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: id is required",
			}},
			IsError: true,
		}, nil
	}

	var id int
	switch v := idInterface.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	case string:
		var err error
		id, err = strconv.Atoi(v)
		if err != nil {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: "Error: id must be a valid integer",
				}},
				IsError: true,
			}, nil
		}
	default:
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: id must be a number",
			}},
			IsError: true,
		}, nil
	}

	// Get todo
	todo, err := s.storage.GetByID(id)
	if err != nil {
		if err == storage.ErrTodoNotFound {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: fmt.Sprintf("Todo with ID %d not found", id),
				}},
				IsError: true,
			}, nil
		}
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error retrieving todo: %v", err),
			}},
			IsError: true,
		}, nil
	}

	// Convert to response format
	todoResp := TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      string(todo.Status),
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	result, _ := json.MarshalIndent(todoResp, "", "  ")
	return &CallToolResponse{
		Content: []Content{{
			Type: "text",
			Text: string(result),
		}},
	}, nil
}

// handleGetTodos handles the get_todos tool
func (s *MCPServer) handleGetTodos(args map[string]interface{}) (*CallToolResponse, error) {
	var todos []*models.Todo
	var err error

	// Check if status filter is provided
	if statusStr, ok := args["status"].(string); ok && statusStr != "" {
		status := models.TodoStatus(statusStr)
		if status != models.StatusPending && status != models.StatusCompleted {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: "Error: status must be 'pending' or 'completed'",
				}},
				IsError: true,
			}, nil
		}
		todos, err = s.storage.GetByStatus(status)
	} else {
		todos, err = s.storage.GetAll()
	}

	if err != nil {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error retrieving todos: %v", err),
			}},
			IsError: true,
		}, nil
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

	result, _ := json.MarshalIndent(todoListResp, "", "  ")
	return &CallToolResponse{
		Content: []Content{{
			Type: "text",
			Text: string(result),
		}},
	}, nil
}

// handleUpdateTodo handles the update_todo tool
func (s *MCPServer) handleUpdateTodo(args map[string]interface{}) (*CallToolResponse, error) {
	// Extract and validate ID
	idInterface, ok := args["id"]
	if !ok {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: id is required",
			}},
			IsError: true,
		}, nil
	}

	var id int
	switch v := idInterface.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	case string:
		var err error
		id, err = strconv.Atoi(v)
		if err != nil {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: "Error: id must be a valid integer",
				}},
				IsError: true,
			}, nil
		}
	default:
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: id must be a number",
			}},
			IsError: true,
		}, nil
	}

	// Get existing todo
	existingTodo, err := s.storage.GetByID(id)
	if err != nil {
		if err == storage.ErrTodoNotFound {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: fmt.Sprintf("Todo with ID %d not found", id),
				}},
				IsError: true,
			}, nil
		}
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error retrieving todo: %v", err),
			}},
			IsError: true,
		}, nil
	}

	// Create updated todo
	updatedTodo := *existingTodo

	// Update fields if provided
	if title, ok := args["title"].(string); ok && title != "" {
		updatedTodo.Title = title
	}

	if description, ok := args["description"].(string); ok {
		updatedTodo.Description = description
	}

	if statusStr, ok := args["status"].(string); ok && statusStr != "" {
		status := models.TodoStatus(statusStr)
		if status != models.StatusPending && status != models.StatusCompleted {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: "Error: status must be 'pending' or 'completed'",
				}},
				IsError: true,
			}, nil
		}
		updatedTodo.Status = status
	}

	// Update in storage
	if err := s.storage.Update(id, &updatedTodo); err != nil {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error updating todo: %v", err),
			}},
			IsError: true,
		}, nil
	}

	// Convert to response format
	todoResp := TodoResponse{
		ID:          updatedTodo.ID,
		Title:       updatedTodo.Title,
		Description: updatedTodo.Description,
		Status:      string(updatedTodo.Status),
		CreatedAt:   updatedTodo.CreatedAt,
		UpdatedAt:   updatedTodo.UpdatedAt,
	}

	result, _ := json.MarshalIndent(todoResp, "", "  ")
	return &CallToolResponse{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Todo updated successfully:\n%s", string(result)),
		}},
	}, nil
}

// handleDeleteTodo handles the delete_todo tool
func (s *MCPServer) handleDeleteTodo(args map[string]interface{}) (*CallToolResponse, error) {
	// Extract and validate ID
	idInterface, ok := args["id"]
	if !ok {
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: id is required",
			}},
			IsError: true,
		}, nil
	}

	var id int
	switch v := idInterface.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	case string:
		var err error
		id, err = strconv.Atoi(v)
		if err != nil {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: "Error: id must be a valid integer",
				}},
				IsError: true,
			}, nil
		}
	default:
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: "Error: id must be a number",
			}},
			IsError: true,
		}, nil
	}

	// Delete todo
	if err := s.storage.Delete(id); err != nil {
		if err == storage.ErrTodoNotFound {
			return &CallToolResponse{
				Content: []Content{{
					Type: "text",
					Text: fmt.Sprintf("Todo with ID %d not found", id),
				}},
				IsError: true,
			}, nil
		}
		return &CallToolResponse{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error deleting todo: %v", err),
			}},
			IsError: true,
		}, nil
	}

	return &CallToolResponse{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Todo with ID %d deleted successfully", id),
		}},
	}, nil
}
