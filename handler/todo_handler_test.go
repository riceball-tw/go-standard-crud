package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-crud-todo-list/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// MockTodoService implements TodoService interface for testing
type MockTodoService struct {
	todos   []models.Todo
	nextID  int
	failGet bool
}

func NewMockTodoService() *MockTodoService {
	return &MockTodoService{
		todos:  make([]models.Todo, 0),
		nextID: 1,
	}
}

func (m *MockTodoService) GetAllTodos() ([]models.Todo, error) {
	if m.failGet {
		return nil, errors.New("service error")
	}
	return m.todos, nil
}

func (m *MockTodoService) GetTodoByID(id int) (*models.Todo, error) {
	if m.failGet {
		return nil, errors.New("service error")
	}
	for _, todo := range m.todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, errors.New("todo not found")
}

func (m *MockTodoService) CreateTodo(title, description string) (*models.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("validation failed: title is required")
	}
	if len(title) > 200 {
		return nil, errors.New("validation failed: title too long")
	}
	
	todo := models.Todo{
		ID:          m.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.nextID++
	m.todos = append(m.todos, todo)
	return &todo, nil
}

func (m *MockTodoService) UpdateTodo(id int, title, description string, completed bool) (*models.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("validation failed: title is required")
	}
	
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos[i].Title = title
			m.todos[i].Description = description
			m.todos[i].Completed = completed
			m.todos[i].UpdatedAt = time.Now()
			return &m.todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func (m *MockTodoService) DeleteTodo(id int) error {
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}

func TestGetAllTodos(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	// Add some test todos
	mockService.CreateTodo("Test Todo 1", "Description 1")
	mockService.CreateTodo("Test Todo 2", "Description 2")
	
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()
	
	handler.getAllTodos(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var todos []models.Todo
	if err := json.NewDecoder(w.Body).Decode(&todos); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestCreateTodo(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	reqBody := CreateTodoRequest{
		Title:       "New Todo",
		Description: "New Description",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.createTodo(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	var todo models.Todo
	if err := json.NewDecoder(w.Body).Decode(&todo); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if todo.Title != reqBody.Title {
		t.Errorf("Expected title %s, got %s", reqBody.Title, todo.Title)
	}
}

func TestSetupRoutes(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	mux := handler.SetupRoutes()
	
	if mux == nil {
		t.Error("Expected non-nil ServeMux")
	}
}

func TestGetAllTodos_ServiceError(t *testing.T) {
	mockService := NewMockTodoService()
	mockService.failGet = true
	handler := NewTodoHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()
	
	handler.getAllTodos(w, req)
	
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetTodoByID(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	// Create a test todo
	createdTodo, _ := mockService.CreateTodo("Test Todo", "Test Description")
	
	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	w := httptest.NewRecorder()
	
	handler.getTodoByID(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var todo models.Todo
	if err := json.NewDecoder(w.Body).Decode(&todo); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if todo.ID != createdTodo.ID {
		t.Errorf("Expected todo ID %d, got %d", createdTodo.ID, todo.ID)
	}
}

func TestGetTodoByID_NotFound(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	w := httptest.NewRecorder()
	
	handler.getTodoByID(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetTodoByID_InvalidID(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/todos/invalid", nil)
	w := httptest.NewRecorder()
	
	handler.getTodoByID(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTodo_InvalidJSON(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.createTodo(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTodo_ValidationError(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	reqBody := CreateTodoRequest{
		Title:       "", // Empty title should cause validation error
		Description: "Description",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.createTodo(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTodo(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	// Create a todo first
	mockService.CreateTodo("Original Title", "Original Description")
	
	reqBody := UpdateTodoRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		Completed:   true,
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.updateTodo(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var todo models.Todo
	if err := json.NewDecoder(w.Body).Decode(&todo); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if todo.Title != reqBody.Title {
		t.Errorf("Expected title %s, got %s", reqBody.Title, todo.Title)
	}
	if todo.Completed != reqBody.Completed {
		t.Errorf("Expected completed %t, got %t", reqBody.Completed, todo.Completed)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	reqBody := UpdateTodoRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		Completed:   true,
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/todos/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.updateTodo(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTodo(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	// Create a todo first
	mockService.CreateTodo("Test Todo", "Test Description")
	
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	w := httptest.NewRecorder()
	
	handler.deleteTodo(w, req)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
	
	// Verify todo was deleted
	if len(mockService.todos) != 0 {
		t.Errorf("Expected 0 todos after deletion, got %d", len(mockService.todos))
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	w := httptest.NewRecorder()
	
	handler.deleteTodo(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestJSONMiddleware(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)
	
	// Test POST without proper content type
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	
	middlewareHandler := handler.jsonMiddleware(handler.createTodo)
	middlewareHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}