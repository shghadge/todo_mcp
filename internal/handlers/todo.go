package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"todo_mcp/internal/models"
	"todo_mcp/internal/storage"

	"github.com/gorilla/mux"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	storage storage.TodoStorage
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(storage storage.TodoStorage) *TodoHandler {
	return &TodoHandler{
		storage: storage,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Validation failed", "Title is required")
		return
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.StatusPending,
	}

	if err := h.storage.Create(todo); err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create todo", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusCreated, "Todo created successfully", todo)
}

// GetTodos handles GET /todos
func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	var todos []*models.Todo
	var err error

	if status != "" {
		todoStatus := models.TodoStatus(status)
		if todoStatus != models.StatusPending && todoStatus != models.StatusCompleted {
			h.sendErrorResponse(w, http.StatusBadRequest, "Invalid status", "Status must be 'pending' or 'completed'")
			return
		}
		todos, err = h.storage.GetByStatus(todoStatus)
	} else {
		todos, err = h.storage.GetAll()
	}

	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todos", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todos retrieved successfully", todos)
}

// GetTodo handles GET /todos/{id}
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid ID", "ID must be a number")
		return
	}

	todo, err := h.storage.GetByID(id)
	if err != nil {
		if err == storage.ErrTodoNotFound {
			h.sendErrorResponse(w, http.StatusNotFound, "Todo not found", "Todo with given ID does not exist")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todo", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todo retrieved successfully", todo)
}

// UpdateTodo handles PUT /todos/{id}
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid ID", "ID must be a number")
		return
	}

	// Get existing todo
	existingTodo, err := h.storage.GetByID(id)
	if err != nil {
		if err == storage.ErrTodoNotFound {
			h.sendErrorResponse(w, http.StatusNotFound, "Todo not found", "Todo with given ID does not exist")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todo", err.Error())
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	// Update fields if provided
	updatedTodo := *existingTodo
	if req.Title != nil {
		if strings.TrimSpace(*req.Title) == "" {
			h.sendErrorResponse(w, http.StatusBadRequest, "Validation failed", "Title cannot be empty")
			return
		}
		updatedTodo.Title = *req.Title
	}
	if req.Description != nil {
		updatedTodo.Description = *req.Description
	}
	if req.Status != nil {
		if *req.Status != models.StatusPending && *req.Status != models.StatusCompleted {
			h.sendErrorResponse(w, http.StatusBadRequest, "Invalid status", "Status must be 'pending' or 'completed'")
			return
		}
		updatedTodo.Status = *req.Status
	}

	if err := h.storage.Update(id, &updatedTodo); err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to update todo", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todo updated successfully", &updatedTodo)
}

// DeleteTodo handles DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid ID", "ID must be a number")
		return
	}

	if err := h.storage.Delete(id); err != nil {
		if err == storage.ErrTodoNotFound {
			h.sendErrorResponse(w, http.StatusNotFound, "Todo not found", "Todo with given ID does not exist")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete todo", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todo deleted successfully", nil)
}

// Helper methods
func (h *TodoHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   error,
		Message: message,
	})
}

func (h *TodoHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	})
}
