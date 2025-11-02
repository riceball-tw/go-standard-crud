package service

import (
	"errors"
	"fmt"
	"go-crud-todo-list/models"
	"go-crud-todo-list/repository"
	"strings"
)

// TodoService defines the interface for todo business logic operations
type TodoService interface {
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id int) (*models.Todo, error)
	CreateTodo(title, description string) (*models.Todo, error)
	UpdateTodo(id int, title, description string, completed bool) (*models.Todo, error)
	DeleteTodo(id int) error
}

// TodoServiceImpl implements the TodoService interface
type TodoServiceImpl struct {
	repository repository.TodoRepository
}

// NewTodoService creates a new TodoService instance with the given repository
func NewTodoService(repo repository.TodoRepository) TodoService {
	return &TodoServiceImpl{
		repository: repo,
	}
}

// validateTodoInput validates input parameters for todo creation and updates
func (s *TodoServiceImpl) validateTodoInput(title, description string) error {
	// Validate title
	if strings.TrimSpace(title) == "" {
		return errors.New("title is required and cannot be empty")
	}
	if len(title) > 200 {
		return errors.New("title must be 200 characters or less")
	}

	// Validate description
	if len(description) > 1000 {
		return errors.New("description must be 1000 characters or less")
	}

	return nil
}

// GetAllTodos retrieves all todos from the repository
func (s *TodoServiceImpl) GetAllTodos() ([]models.Todo, error) {
	todos, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todos: %w", err)
	}
	return todos, nil
}

// GetTodoByID retrieves a specific todo by its ID
func (s *TodoServiceImpl) GetTodoByID(id int) (*models.Todo, error) {
	if id <= 0 {
		return nil, errors.New("invalid todo ID: ID must be a positive integer")
	}

	todo, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}
	return todo, nil
}

// CreateTodo creates a new todo with the provided title and description
func (s *TodoServiceImpl) CreateTodo(title, description string) (*models.Todo, error) {
	// Validate input
	if err := s.validateTodoInput(title, description); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create new todo
	todo := &models.Todo{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Completed:   false,
	}

	// Save to repository
	if err := s.repository.Create(todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

// UpdateTodo updates an existing todo with new values
func (s *TodoServiceImpl) UpdateTodo(id int, title, description string, completed bool) (*models.Todo, error) {
	if id <= 0 {
		return nil, errors.New("invalid todo ID: ID must be a positive integer")
	}

	// Validate input
	if err := s.validateTodoInput(title, description); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if todo exists
	existingTodo, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}

	// Create updated todo with new values
	updatedTodo := &models.Todo{
		ID:          existingTodo.ID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Completed:   completed,
		CreatedAt:   existingTodo.CreatedAt, // Preserve original creation time
	}

	// Update in repository
	if err := s.repository.Update(id, updatedTodo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return updatedTodo, nil
}

// DeleteTodo removes a todo by its ID
func (s *TodoServiceImpl) DeleteTodo(id int) error {
	if id <= 0 {
		return errors.New("invalid todo ID: ID must be a positive integer")
	}

	// Check if todo exists before attempting deletion
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("todo not found: %w", err)
	}

	// Delete from repository
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}