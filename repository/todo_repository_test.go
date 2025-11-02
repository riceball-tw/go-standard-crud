package repository

import (
	"fmt"
	"go-crud-todo-list/models"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// createTempFile creates a temporary file for testing
func createTempFile(t *testing.T) string {
	tempDir := t.TempDir()
	return filepath.Join(tempDir, "test_todos.json")
}

// createTestTodo creates a test todo with valid data
func createTestTodo() models.Todo {
	return models.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
	}
}

func TestNewFileBasedTodoRepository(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}

	if repo.filePath != filePath {
		t.Errorf("Expected filePath %s, got %s", filePath, repo.filePath)
	}

	if repo.storage == nil {
		t.Fatal("Expected storage to be initialized, got nil")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	err := repo.Load()
	if err != nil {
		t.Errorf("Expected no error when loading non-existent file, got %v", err)
	}

	todos, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error getting all todos, got %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected empty todos list, got %d todos", len(todos))
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	filePath := createTempFile(t)
	
	// Create empty file
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	repo := NewFileBasedTodoRepository(filePath)
	err = repo.Load()
	if err != nil {
		t.Errorf("Expected no error when loading empty file, got %v", err)
	}

	todos, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error getting all todos, got %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected empty todos list, got %d todos", len(todos))
	}
}

func TestSave_And_Load(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	// Create and save a todo
	todo := createTestTodo()
	err := repo.Create(&todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Create new repository instance and load
	repo2 := NewFileBasedTodoRepository(filePath)
	err = repo2.Load()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	todos, err := repo2.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}

	if todos[0].Title != todo.Title {
		t.Errorf("Expected title %s, got %s", todo.Title, todos[0].Title)
	}
}

func TestCreate(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	todo := createTestTodo()
	originalTitle := todo.Title

	err := repo.Create(&todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Check that ID was assigned
	if todo.ID == 0 {
		t.Error("Expected ID to be assigned, got 0")
	}

	// Check that timestamps were set
	if todo.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if todo.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	// Check that title is preserved
	if todo.Title != originalTitle {
		t.Errorf("Expected title %s, got %s", originalTitle, todo.Title)
	}

	// Verify it's in storage
	todos, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo in storage, got %d", len(todos))
	}
}

func TestCreate_NilTodo(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	err := repo.Create(nil)
	if err == nil {
		t.Error("Expected error when creating nil todo, got nil")
	}
}

func TestCreate_InvalidTodo(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	// Create todo with empty title (invalid)
	todo := models.Todo{
		Title:       "",
		Description: "Test Description",
		Completed:   false,
	}

	err := repo.Create(&todo)
	if err == nil {
		t.Error("Expected error when creating invalid todo, got nil")
	}
}

func TestGetAll(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	// Initially should be empty
	todos, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos initially, got %d", len(todos))
	}

	// Add some todos
	todo1 := createTestTodo()
	todo1.Title = "Todo 1"
	
	todo2 := createTestTodo()
	todo2.Title = "Todo 2"

	err = repo.Create(&todo1)
	if err != nil {
		t.Fatalf("Failed to create todo1: %v", err)
	}

	err = repo.Create(&todo2)
	if err != nil {
		t.Fatalf("Failed to create todo2: %v", err)
	}

	// Get all todos
	todos, err = repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestGetByID(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	// Create a todo
	todo := createTestTodo()
	err := repo.Create(&todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Get by ID
	foundTodo, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("Failed to get todo by ID: %v", err)
	}

	if foundTodo.ID != todo.ID {
		t.Errorf("Expected ID %d, got %d", todo.ID, foundTodo.ID)
	}

	if foundTodo.Title != todo.Title {
		t.Errorf("Expected title %s, got %s", todo.Title, foundTodo.Title)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("Expected error when getting non-existent todo, got nil")
	}
}

func TestUpdate(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	// Create a todo
	todo := createTestTodo()
	err := repo.Create(&todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	originalCreatedAt := todo.CreatedAt
	originalID := todo.ID

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update the todo
	updatedTodo := models.Todo{
		Title:       "Updated Title",
		Description: "Updated Description",
		Completed:   true,
	}

	err = repo.Update(todo.ID, &updatedTodo)
	if err != nil {
		t.Fatalf("Failed to update todo: %v", err)
	}

	// Check that original creation time and ID are preserved
	if updatedTodo.ID != originalID {
		t.Errorf("Expected ID to be preserved: %d, got %d", originalID, updatedTodo.ID)
	}

	if !updatedTodo.CreatedAt.Equal(originalCreatedAt) {
		t.Errorf("Expected CreatedAt to be preserved: %v, got %v", originalCreatedAt, updatedTodo.CreatedAt)
	}

	// Check that UpdatedAt was changed
	if updatedTodo.UpdatedAt.Equal(originalCreatedAt) {
		t.Error("Expected UpdatedAt to be different from CreatedAt")
	}

	// Verify changes in storage
	foundTodo, err := repo.GetByID(originalID)
	if err != nil {
		t.Fatalf("Failed to get updated todo: %v", err)
	}

	if foundTodo.Title != "Updated Title" {
		t.Errorf("Expected updated title, got %s", foundTodo.Title)
	}

	if !foundTodo.Completed {
		t.Error("Expected todo to be completed")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	todo := createTestTodo()
	err := repo.Update(999, &todo)
	if err == nil {
		t.Error("Expected error when updating non-existent todo, got nil")
	}
}

func TestUpdate_NilTodo(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	err := repo.Update(1, nil)
	if err == nil {
		t.Error("Expected error when updating with nil todo, got nil")
	}
}

func TestDelete(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	// Create a todo
	todo := createTestTodo()
	err := repo.Create(&todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Verify it exists
	_, err = repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("Todo should exist before deletion: %v", err)
	}

	// Delete the todo
	err = repo.Delete(todo.ID)
	if err != nil {
		t.Fatalf("Failed to delete todo: %v", err)
	}

	// Verify it's gone
	_, err = repo.GetByID(todo.ID)
	if err == nil {
		t.Error("Expected error when getting deleted todo, got nil")
	}

	// Verify it's not in the list
	todos, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos after deletion, got %d", len(todos))
	}
}

func TestDelete_NotFound(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	err := repo.Delete(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent todo, got nil")
	}
}

// TestConcurrentAccess tests thread safety with concurrent operations
func TestConcurrentAccess(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewFileBasedTodoRepository(filePath)

	const numGoroutines = 10
	const todosPerGoroutine = 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*todosPerGoroutine)

	// Concurrent creates
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < todosPerGoroutine; j++ {
				todo := createTestTodo()
				todo.Title = fmt.Sprintf("Todo %d-%d", goroutineID, j)
				if err := repo.Create(&todo); err != nil {
					errors <- err
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent operation failed: %v", err)
	}

	// Verify all todos were created
	todos, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	expectedCount := numGoroutines * todosPerGoroutine
	if len(todos) != expectedCount {
		t.Errorf("Expected %d todos, got %d", expectedCount, len(todos))
	}
}

// TestFilePersistence tests that data survives repository recreation
func TestFilePersistence(t *testing.T) {
	filePath := createTempFile(t)

	// Create repository and add data
	repo1 := NewFileBasedTodoRepository(filePath)
	todo := createTestTodo()
	err := repo1.Create(&todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Create new repository instance and load data
	repo2 := NewFileBasedTodoRepository(filePath)
	err = repo2.Load()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	// Verify data persisted
	todos, err := repo2.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all todos: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 persisted todo, got %d", len(todos))
	}

	if todos[0].Title != todo.Title {
		t.Errorf("Expected persisted title %s, got %s", todo.Title, todos[0].Title)
	}
}