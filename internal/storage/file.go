package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shghadge/todo_mcp/internal/models"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrTodoExists   = errors.New("todo already exists")
)

// FileStorage implements TodoStorage using JSON file storage
type FileStorage struct {
	filePath string
	mutex    sync.RWMutex
}

// NewFileStorage creates a new file-based storage instance
func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

// loadTodos loads todos from the JSON file
func (f *FileStorage) loadTodos() (map[int]*models.Todo, int, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	// Check if file exists
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		// File doesn't exist, return empty map
		return make(map[int]*models.Todo), 1, nil
	}

	// Read file
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var todos map[int]*models.Todo
	if len(data) == 0 {
		// Empty file, return empty map
		return make(map[int]*models.Todo), 1, nil
	}

	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Find the highest ID to determine next ID
	nextID := 1
	for id := range todos {
		if id >= nextID {
			nextID = id + 1
		}
	}

	return todos, nextID, nil
}

// saveTodos saves todos to the JSON file
func (f *FileStorage) saveTodos(todos map[int]*models.Todo) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(f.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(f.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Create creates a new todo item
func (f *FileStorage) Create(todo *models.Todo) error {
	todos, nextID, err := f.loadTodos()
	if err != nil {
		return err
	}

	todo.ID = nextID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	todos[todo.ID] = todo

	return f.saveTodos(todos)
}

// GetByID retrieves a todo by its ID
func (f *FileStorage) GetByID(id int) (*models.Todo, error) {
	todos, _, err := f.loadTodos()
	if err != nil {
		return nil, err
	}

	todo, exists := todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}

	// Return a copy to avoid race conditions
	todoCopy := *todo
	return &todoCopy, nil
}

// GetAll retrieves all todos
func (f *FileStorage) GetAll() ([]*models.Todo, error) {
	todos, _, err := f.loadTodos()
	if err != nil {
		return nil, err
	}

	result := make([]*models.Todo, 0, len(todos))
	for _, todo := range todos {
		// Add a copy to avoid race conditions
		todoCopy := *todo
		result = append(result, &todoCopy)
	}

	return result, nil
}

// Update updates an existing todo
func (f *FileStorage) Update(id int, updatedTodo *models.Todo) error {
	todos, _, err := f.loadTodos()
	if err != nil {
		return err
	}

	todo, exists := todos[id]
	if !exists {
		return ErrTodoNotFound
	}

	// Preserve original ID and CreatedAt
	updatedTodo.ID = todo.ID
	updatedTodo.CreatedAt = todo.CreatedAt
	updatedTodo.UpdatedAt = time.Now()

	todos[id] = updatedTodo
	return f.saveTodos(todos)
}

// Delete deletes a todo by its ID
func (f *FileStorage) Delete(id int) error {
	todos, _, err := f.loadTodos()
	if err != nil {
		return err
	}

	_, exists := todos[id]
	if !exists {
		return ErrTodoNotFound
	}

	delete(todos, id)
	return f.saveTodos(todos)
}

// GetByStatus retrieves todos by status
func (f *FileStorage) GetByStatus(status models.TodoStatus) ([]*models.Todo, error) {
	todos, _, err := f.loadTodos()
	if err != nil {
		return nil, err
	}

	var result []*models.Todo
	for _, todo := range todos {
		if todo.Status == status {
			// Add a copy to avoid race conditions
			todoCopy := *todo
			result = append(result, &todoCopy)
		}
	}

	return result, nil
}
