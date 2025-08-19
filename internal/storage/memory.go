package storage

import (
	"errors"
	"sync"
	"time"
	"todo_mcp/internal/models"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrTodoExists   = errors.New("todo already exists")
)

// InMemoryStorage implements TodoStorage using in-memory storage
type InMemoryStorage struct {
	todos  map[int]*models.Todo
	nextID int
	mutex  sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		todos:  make(map[int]*models.Todo),
		nextID: 1,
	}
}

// Create creates a new todo item
func (s *InMemoryStorage) Create(todo *models.Todo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	todo.ID = s.nextID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	s.todos[todo.ID] = todo
	s.nextID++

	return nil
}

// GetByID retrieves a todo by its ID
func (s *InMemoryStorage) GetByID(id int) (*models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}

	// Return a copy to avoid race conditions
	todoCopy := *todo
	return &todoCopy, nil
}

// GetAll retrieves all todos
func (s *InMemoryStorage) GetAll() ([]*models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todos := make([]*models.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		// Add a copy to avoid race conditions
		todoCopy := *todo
		todos = append(todos, &todoCopy)
	}

	return todos, nil
}

// Update updates an existing todo
func (s *InMemoryStorage) Update(id int, updatedTodo *models.Todo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return ErrTodoNotFound
	}

	// Preserve original ID and CreatedAt
	updatedTodo.ID = todo.ID
	updatedTodo.CreatedAt = todo.CreatedAt
	updatedTodo.UpdatedAt = time.Now()

	s.todos[id] = updatedTodo
	return nil
}

// Delete deletes a todo by its ID
func (s *InMemoryStorage) Delete(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.todos[id]
	if !exists {
		return ErrTodoNotFound
	}

	delete(s.todos, id)
	return nil
}

// GetByStatus retrieves todos by status
func (s *InMemoryStorage) GetByStatus(status models.TodoStatus) ([]*models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var todos []*models.Todo
	for _, todo := range s.todos {
		if todo.Status == status {
			// Add a copy to avoid race conditions
			todoCopy := *todo
			todos = append(todos, &todoCopy)
		}
	}

	return todos, nil
}
