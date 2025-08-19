package storage

import (
	"todo_mcp/internal/models"
)

// TodoStorage defines the interface for todo storage operations
type TodoStorage interface {
	// Create creates a new todo item
	Create(todo *models.Todo) error

	// GetByID retrieves a todo by its ID
	GetByID(id int) (*models.Todo, error)

	// GetAll retrieves all todos
	GetAll() ([]*models.Todo, error)

	// Update updates an existing todo
	Update(id int, todo *models.Todo) error

	// Delete deletes a todo by its ID
	Delete(id int) error

	// GetByStatus retrieves todos by status
	GetByStatus(status models.TodoStatus) ([]*models.Todo, error)
}
