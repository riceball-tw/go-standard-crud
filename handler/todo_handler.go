package handler

import (
	"encoding/json"
	"fmt"
	"go-crud-todo-list/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	service service.TodoService
}

// NewTodoHandler creates a new TodoHandler with the given service
func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// writeErrorResponse writes an error response with the specified status code and message
func (h *TodoHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResp := ErrorResponse{
		Error:     message,
		Code:      statusCode,
		Timestamp: time.Now(),
	}
	
	json.NewEncoder(w).Encode(errorResp)
}

// writeJSONResponse writes a JSON response with the specified status code and data
func (h *TodoHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// extractIDFromPath extracts the ID parameter from the URL path
func (h *TodoHandler) extractIDFromPath(path string) (int, error) {
	// Expected path format: /todos/{id}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] != "todos" {
		return 0, fmt.Errorf("invalid path format")
	}
	
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	
	return id, nil
}

// SetupRoutes configures the HTTP routes and returns a ServeMux
func (h *TodoHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Apply JSON middleware to all routes
	mux.HandleFunc("/todos", h.jsonMiddleware(h.todosHandler))
	mux.HandleFunc("/todos/", h.jsonMiddleware(h.todoByIDHandler))
	
	return mux
}

// jsonMiddleware adds JSON content type handling to HTTP handlers
func (h *TodoHandler) jsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set default content type for responses
		w.Header().Set("Content-Type", "application/json")
		
		// For POST and PUT requests, validate content type
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				h.writeErrorResponse(w, http.StatusBadRequest, "Content-Type must be application/json")
				return
			}
		}
		
		next(w, r)
	}
}

// todosHandler handles requests to /todos endpoint
func (h *TodoHandler) todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllTodos(w, r)
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// todoByIDHandler handles requests to /todos/{id} endpoint
func (h *TodoHandler) todoByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTodoByID(w, r)
	case http.MethodPut:
		h.updateTodo(w, r)
	case http.MethodDelete:
		h.deleteTodo(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getAllTodos handles GET /todos - returns all todos as JSON
func (h *TodoHandler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todos")
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, todos)
}

// getTodoByID handles GET /todos/{id} - returns a specific todo by ID
func (h *TodoHandler) getTodoByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	
	// Get todo from service
	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		// Check if it's a not found error
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Todo not found")
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todo")
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, todo)
}

// createTodo handles POST /todos - creates a new todo
func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	
	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	
	// Create todo using service
	todo, err := h.service.CreateTodo(req.Title, req.Description)
	if err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation failed") {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}
	
	h.writeJSONResponse(w, http.StatusCreated, todo)
}

// updateTodo handles PUT /todos/{id} - updates an existing todo
func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	
	var req UpdateTodoRequest
	
	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	
	// Update todo using service
	todo, err := h.service.UpdateTodo(id, req.Title, req.Description, req.Completed)
	if err != nil {
		// Check error type and respond accordingly
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Todo not found")
			return
		}
		if strings.Contains(err.Error(), "validation failed") {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, todo)
}

// deleteTodo handles DELETE /todos/{id} - deletes a todo by ID
func (h *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	
	// Delete todo using service
	err = h.service.DeleteTodo(id)
	if err != nil {
		// Check if it's a not found error
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Todo not found")
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}
	
	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}