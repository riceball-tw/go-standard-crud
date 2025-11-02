package service

import (
	"errors"
	"go-crud-todo-list/models"
	"strings"
	"testing"
	"time"
)

// MockTodoRepository is a mock implementation of TodoRepository for testing
type MockTodoRepository struct {
	todos   map[int]*models.Todo
	nextID  int
	loadErr error
	saveErr error
}

// NewMockTodoRepository creates a new mock repository
func NewMockTodoRepository() *MockTodoRepository {
	return &MockTodoRepository{
		todos:  make(map[int]*models.Todo),
		nextID: 1,
	}
}

// SetLoadError sets an error to be returned by Load method
func (m *MockTodoRepository) SetLoadError(err error) {
	m.loadErr = err
}

// SetSaveError sets an error to be returned by Save method
func (m *MockTodoRepository) SetSaveError(err error) {
	m.saveErr = err
}

// GetAll returns all todos from the mock repository
func (m *MockTodoRepository) GetAll() ([]models.Todo, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	
	todos := make([]models.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, *todo)
	}
	return todos, nil
}

// GetByID returns a specific todo by ID from the mock repository
func (m *MockTodoRepository) GetByID(id int) (*models.Todo, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	
	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}
	
	// Return a copy
	todoCopy := *todo
	return &todoCopy, nil
}

// Create adds a new todo to the mock repository
func (m *MockTodoRepository) Create(todo *models.Todo) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	
	// Validate todo
	if err := todo.Validate(); err != nil {
		return err
	}
	
	// Assign ID and timestamps
	todo.ID = m.nextID
	m.nextID++
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	
	// Store copy
	todoCopy := *todo
	m.todos[todo.ID] = &todoCopy
	
	return nil
}

// Update modifies an existing todo in the mock repository
func (m *MockTodoRepository) Update(id int, todo *models.Todo) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	
	// Check if todo exists
	existingTodo, exists := m.todos[id]
	if !exists {
		return errors.New("todo not found")
	}
	
	// Validate todo
	if err := todo.Validate(); err != nil {
		return err
	}
	
	// Preserve ID and creation time, update timestamp
	todo.ID = id
	todo.CreatedAt = existingTodo.CreatedAt
	todo.UpdatedAt = time.Now()
	
	// Store copy
	todoCopy := *todo
	m.todos[id] = &todoCopy
	
	return nil
}

// Delete removes a todo from the mock repository
func (m *MockTodoRepository) Delete(id int) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	
	if _, exists := m.todos[id]; !exists {
		return errors.New("todo not found")
	}
	
	delete(m.todos, id)
	return nil
}

// Save is a no-op for the mock repository
func (m *MockTodoRepository) Save() error {
	return m.saveErr
}

// Load is a no-op for the mock repository
func (m *MockTodoRepository) Load() error {
	return m.loadErr
}

// Test helper functions

func createTestTodo(id int, title, description string, completed bool) *models.Todo {
	now := time.Now()
	return &models.Todo{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   completed,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// TestNewTodoService tests the service constructor
func TestNewTodoService(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
}

// TestGetAllTodos tests retrieving all todos
func TestGetAllTodos(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	// Test empty repository
	todos, err := service.GetAllTodos()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("Expected 0 todos, got %d", len(todos))
	}
	
	// Add some test todos
	testTodo1 := createTestTodo(1, "Test Todo 1", "Description 1", false)
	testTodo2 := createTestTodo(2, "Test Todo 2", "Description 2", true)
	mockRepo.todos[1] = testTodo1
	mockRepo.todos[2] = testTodo2
	
	todos, err = service.GetAllTodos()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("Expected 2 todos, got %d", len(todos))
	}
}

// TestGetAllTodos_RepositoryError tests error handling in GetAllTodos
func TestGetAllTodos_RepositoryError(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	mockRepo.SetLoadError(errors.New("repository error"))
	service := NewTodoService(mockRepo)
	
	_, err := service.GetAllTodos()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to retrieve todos") {
		t.Fatalf("Expected error message to contain 'failed to retrieve todos', got %v", err)
	}
}

// TestGetTodoByID tests retrieving a specific todo
func TestGetTodoByID(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	// Add test todo
	testTodo := createTestTodo(1, "Test Todo", "Test Description", false)
	mockRepo.todos[1] = testTodo
	
	// Test successful retrieval
	todo, err := service.GetTodoByID(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if todo.ID != 1 {
		t.Fatalf("Expected todo ID 1, got %d", todo.ID)
	}
	if todo.Title != "Test Todo" {
		t.Fatalf("Expected title 'Test Todo', got %s", todo.Title)
	}
}

// TestGetTodoByID_InvalidID tests invalid ID validation
func TestGetTodoByID_InvalidID(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	testCases := []int{0, -1, -100}
	
	for _, id := range testCases {
		_, err := service.GetTodoByID(id)
		if err == nil {
			t.Fatalf("Expected error for invalid ID %d, got nil", id)
		}
		if !strings.Contains(err.Error(), "invalid todo ID") {
			t.Fatalf("Expected error message to contain 'invalid todo ID', got %v", err)
		}
	}
}

// TestGetTodoByID_NotFound tests todo not found scenario
func TestGetTodoByID_NotFound(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	_, err := service.GetTodoByID(999)
	if err == nil {
		t.Fatal("Expected error for non-existent todo, got nil")
	}
	if !strings.Contains(err.Error(), "todo not found") {
		t.Fatalf("Expected error message to contain 'todo not found', got %v", err)
	}
}

// TestCreateTodo tests creating a new todo
func TestCreateTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	todo, err := service.CreateTodo("Test Todo", "Test Description")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if todo.ID != 1 {
		t.Fatalf("Expected todo ID 1, got %d", todo.ID)
	}
	if todo.Title != "Test Todo" {
		t.Fatalf("Expected title 'Test Todo', got %s", todo.Title)
	}
	if todo.Description != "Test Description" {
		t.Fatalf("Expected description 'Test Description', got %s", todo.Description)
	}
	if todo.Completed != false {
		t.Fatalf("Expected completed to be false, got %v", todo.Completed)
	}
}

// TestCreateTodo_ValidationErrors tests input validation
func TestCreateTodo_ValidationErrors(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	testCases := []struct {
		title       string
		description string
		expectedErr string
	}{
		{"", "Valid description", "title is required"},
		{"   ", "Valid description", "title is required"},
		{strings.Repeat("a", 201), "Valid description", "title must be 200 characters or less"},
		{"Valid title", strings.Repeat("a", 1001), "description must be 1000 characters or less"},
	}
	
	for _, tc := range testCases {
		_, err := service.CreateTodo(tc.title, tc.description)
		if err == nil {
			t.Fatalf("Expected error for title '%s' and description length %d, got nil", tc.title, len(tc.description))
		}
		if !strings.Contains(err.Error(), tc.expectedErr) {
			t.Fatalf("Expected error to contain '%s', got %v", tc.expectedErr, err)
		}
	}
}

// TestCreateTodo_RepositoryError tests repository error handling
func TestCreateTodo_RepositoryError(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	mockRepo.SetSaveError(errors.New("repository error"))
	service := NewTodoService(mockRepo)
	
	_, err := service.CreateTodo("Valid Title", "Valid Description")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to create todo") {
		t.Fatalf("Expected error message to contain 'failed to create todo', got %v", err)
	}
}

// TestUpdateTodo tests updating an existing todo
func TestUpdateTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	// Add existing todo
	existingTodo := createTestTodo(1, "Original Title", "Original Description", false)
	mockRepo.todos[1] = existingTodo
	
	updatedTodo, err := service.UpdateTodo(1, "Updated Title", "Updated Description", true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if updatedTodo.ID != 1 {
		t.Fatalf("Expected todo ID 1, got %d", updatedTodo.ID)
	}
	if updatedTodo.Title != "Updated Title" {
		t.Fatalf("Expected title 'Updated Title', got %s", updatedTodo.Title)
	}
	if updatedTodo.Description != "Updated Description" {
		t.Fatalf("Expected description 'Updated Description', got %s", updatedTodo.Description)
	}
	if updatedTodo.Completed != true {
		t.Fatalf("Expected completed to be true, got %v", updatedTodo.Completed)
	}
	if updatedTodo.CreatedAt != existingTodo.CreatedAt {
		t.Fatal("Expected creation timestamp to be preserved")
	}
}

// TestUpdateTodo_InvalidID tests invalid ID validation for updates
func TestUpdateTodo_InvalidID(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	testCases := []int{0, -1, -100}
	
	for _, id := range testCases {
		_, err := service.UpdateTodo(id, "Valid Title", "Valid Description", false)
		if err == nil {
			t.Fatalf("Expected error for invalid ID %d, got nil", id)
		}
		if !strings.Contains(err.Error(), "invalid todo ID") {
			t.Fatalf("Expected error message to contain 'invalid todo ID', got %v", err)
		}
	}
}

// TestUpdateTodo_NotFound tests updating non-existent todo
func TestUpdateTodo_NotFound(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	_, err := service.UpdateTodo(999, "Valid Title", "Valid Description", false)
	if err == nil {
		t.Fatal("Expected error for non-existent todo, got nil")
	}
	if !strings.Contains(err.Error(), "todo not found") {
		t.Fatalf("Expected error message to contain 'todo not found', got %v", err)
	}
}

// TestUpdateTodo_ValidationErrors tests validation errors during update
func TestUpdateTodo_ValidationErrors(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	// Add existing todo
	existingTodo := createTestTodo(1, "Original Title", "Original Description", false)
	mockRepo.todos[1] = existingTodo
	
	testCases := []struct {
		title       string
		description string
		expectedErr string
	}{
		{"", "Valid description", "title is required"},
		{"   ", "Valid description", "title is required"},
		{strings.Repeat("a", 201), "Valid description", "title must be 200 characters or less"},
		{"Valid title", strings.Repeat("a", 1001), "description must be 1000 characters or less"},
	}
	
	for _, tc := range testCases {
		_, err := service.UpdateTodo(1, tc.title, tc.description, false)
		if err == nil {
			t.Fatalf("Expected error for title '%s' and description length %d, got nil", tc.title, len(tc.description))
		}
		if !strings.Contains(err.Error(), tc.expectedErr) {
			t.Fatalf("Expected error to contain '%s', got %v", tc.expectedErr, err)
		}
	}
}

// TestDeleteTodo tests deleting a todo
func TestDeleteTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	// Add test todo
	testTodo := createTestTodo(1, "Test Todo", "Test Description", false)
	mockRepo.todos[1] = testTodo
	
	err := service.DeleteTodo(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Verify todo was deleted
	if _, exists := mockRepo.todos[1]; exists {
		t.Fatal("Expected todo to be deleted, but it still exists")
	}
}

// TestDeleteTodo_InvalidID tests invalid ID validation for deletion
func TestDeleteTodo_InvalidID(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	testCases := []int{0, -1, -100}
	
	for _, id := range testCases {
		err := service.DeleteTodo(id)
		if err == nil {
			t.Fatalf("Expected error for invalid ID %d, got nil", id)
		}
		if !strings.Contains(err.Error(), "invalid todo ID") {
			t.Fatalf("Expected error message to contain 'invalid todo ID', got %v", err)
		}
	}
}

// TestDeleteTodo_NotFound tests deleting non-existent todo
func TestDeleteTodo_NotFound(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	err := service.DeleteTodo(999)
	if err == nil {
		t.Fatal("Expected error for non-existent todo, got nil")
	}
	if !strings.Contains(err.Error(), "todo not found") {
		t.Fatalf("Expected error message to contain 'todo not found', got %v", err)
	}
}

// TestDeleteTodo_RepositoryError tests repository error handling during deletion
func TestDeleteTodo_RepositoryError(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	// Add test todo
	testTodo := createTestTodo(1, "Test Todo", "Test Description", false)
	mockRepo.todos[1] = testTodo
	
	// Set repository error
	mockRepo.SetSaveError(errors.New("repository error"))
	
	err := service.DeleteTodo(1)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to delete todo") {
		t.Fatalf("Expected error message to contain 'failed to delete todo', got %v", err)
	}
}

// TestTrimWhitespace tests that input is properly trimmed
func TestTrimWhitespace(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	
	todo, err := service.CreateTodo("  Test Todo  ", "  Test Description  ")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if todo.Title != "Test Todo" {
		t.Fatalf("Expected trimmed title 'Test Todo', got '%s'", todo.Title)
	}
	if todo.Description != "Test Description" {
		t.Fatalf("Expected trimmed description 'Test Description', got '%s'", todo.Description)
	}
}