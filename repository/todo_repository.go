package repository

import (
	"encoding/json"
	"fmt"
	"go-crud-todo-list/models"
	"os"
	"sync"
)

// TodoRepository defines the interface for todo data persistence operations
type TodoRepository interface {
	GetAll() ([]models.Todo, error)
	GetByID(id int) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(id int, todo *models.Todo) error
	Delete(id int) error
	Save() error
	Load() error
}

// FileBasedTodoRepository implements TodoRepository using file-based persistence
type FileBasedTodoRepository struct {
	storage  *models.TodoStorage
	filePath string
	mutex    sync.RWMutex
}

// NewFileBasedTodoRepository creates a new file-based repository instance
func NewFileBasedTodoRepository(filePath string) *FileBasedTodoRepository {
	return &FileBasedTodoRepository{
		storage:  models.NewTodoStorage(),
		filePath: filePath,
	}
}

// Load reads todo data from the JSON file into memory
func (r *FileBasedTodoRepository) Load() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		// File doesn't exist, start with empty storage
		r.storage = models.NewTodoStorage()
		return nil
	}

	// Read file contents
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		r.storage = models.NewTodoStorage()
		return nil
	}

	// Unmarshal JSON data
	var storage models.TodoStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	r.storage = &storage
	return nil
}

// Save writes the current todo data to the JSON file
func (r *FileBasedTodoRepository) Save() error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Marshal storage to JSON
	data, err := json.MarshalIndent(r.storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetAll returns all todos from the repository
func (r *FileBasedTodoRepository) GetAll() ([]models.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.storage.GetAllTodos(), nil
}

// GetByID returns a specific todo by its ID
func (r *FileBasedTodoRepository) GetByID(id int) (*models.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	todo, _, err := r.storage.FindTodoByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo with ID %d not found", id)
	}

	// Return a copy to prevent external modification
	todoCopy := *todo
	return &todoCopy, nil
}

// Create adds a new todo to the repository
func (r *FileBasedTodoRepository) Create(todo *models.Todo) error {
	if todo == nil {
		return fmt.Errorf("todo cannot be nil")
	}

	// Validate the todo
	if err := todo.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Add todo to storage (this will assign ID and timestamps)
	createdTodo := r.storage.AddTodo(*todo)
	*todo = createdTodo

	// Save to file
	if err := r.saveUnsafe(); err != nil {
		return fmt.Errorf("failed to save todo: %w", err)
	}

	return nil
}

// Update modifies an existing todo in the repository
func (r *FileBasedTodoRepository) Update(id int, todo *models.Todo) error {
	if todo == nil {
		return fmt.Errorf("todo cannot be nil")
	}

	// Validate the todo
	if err := todo.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update todo in storage
	updatedTodo, err := r.storage.UpdateTodo(id, *todo)
	if err != nil {
		return fmt.Errorf("failed to update todo with ID %d: %w", id, err)
	}

	*todo = *updatedTodo

	// Save to file
	if err := r.saveUnsafe(); err != nil {
		return fmt.Errorf("failed to save updated todo: %w", err)
	}

	return nil
}

// Delete removes a todo from the repository
func (r *FileBasedTodoRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Delete todo from storage
	if err := r.storage.DeleteTodo(id); err != nil {
		return fmt.Errorf("failed to delete todo with ID %d: %w", id, err)
	}

	// Save to file
	if err := r.saveUnsafe(); err != nil {
		return fmt.Errorf("failed to save after deletion: %w", err)
	}

	return nil
}

// saveUnsafe saves data without acquiring mutex (internal use only)
func (r *FileBasedTodoRepository) saveUnsafe() error {
	// Marshal storage to JSON
	data, err := json.MarshalIndent(r.storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}