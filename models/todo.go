package models

import (
	"errors"
	"strings"
	"time"
)

// Todo represents a todo item with all required fields
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ValidateTitle validates the todo title according to requirements
func (t *Todo) ValidateTitle() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title is required")
	}
	if len(t.Title) > 200 {
		return errors.New("title must be 200 characters or less")
	}
	return nil
}

// ValidateDescription validates the todo description according to requirements
func (t *Todo) ValidateDescription() error {
	if len(t.Description) > 1000 {
		return errors.New("description must be 1000 characters or less")
	}
	return nil
}

// Validate performs full validation of the todo item
func (t *Todo) Validate() error {
	if err := t.ValidateTitle(); err != nil {
		return err
	}
	if err := t.ValidateDescription(); err != nil {
		return err
	}
	return nil
}

// SetTimestamps sets the creation and update timestamps
func (t *Todo) SetTimestamps() {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
}

// TodoStorage represents the storage structure for file-based persistence
type TodoStorage struct {
	Todos  []Todo `json:"todos"`
	NextID int    `json:"next_id"`
}

// NewTodoStorage creates a new TodoStorage instance with initial values
func NewTodoStorage() *TodoStorage {
	return &TodoStorage{
		Todos:  make([]Todo, 0),
		NextID: 1,
	}
}

// GenerateNextID returns the next available ID and increments the counter
func (ts *TodoStorage) GenerateNextID() int {
	id := ts.NextID
	ts.NextID++
	return id
}

// AddTodo adds a new todo to the storage and assigns it an ID
func (ts *TodoStorage) AddTodo(todo Todo) Todo {
	todo.ID = ts.GenerateNextID()
	todo.SetTimestamps()
	ts.Todos = append(ts.Todos, todo)
	return todo
}

// FindTodoByID finds a todo by its ID and returns it with its index
func (ts *TodoStorage) FindTodoByID(id int) (*Todo, int, error) {
	for i, todo := range ts.Todos {
		if todo.ID == id {
			return &ts.Todos[i], i, nil
		}
	}
	return nil, -1, errors.New("todo not found")
}

// UpdateTodo updates an existing todo in the storage
func (ts *TodoStorage) UpdateTodo(id int, updatedTodo Todo) (*Todo, error) {
	todo, index, err := ts.FindTodoByID(id)
	if err != nil {
		return nil, err
	}
	
	// Preserve original creation time and ID
	updatedTodo.ID = todo.ID
	updatedTodo.CreatedAt = todo.CreatedAt
	updatedTodo.UpdatedAt = time.Now()
	
	ts.Todos[index] = updatedTodo
	return &ts.Todos[index], nil
}

// DeleteTodo removes a todo from the storage by ID
func (ts *TodoStorage) DeleteTodo(id int) error {
	_, index, err := ts.FindTodoByID(id)
	if err != nil {
		return err
	}
	
	// Remove todo from slice
	ts.Todos = append(ts.Todos[:index], ts.Todos[index+1:]...)
	return nil
}

// GetAllTodos returns a copy of all todos in the storage
func (ts *TodoStorage) GetAllTodos() []Todo {
	// Return a copy to prevent external modification
	todos := make([]Todo, len(ts.Todos))
	copy(todos, ts.Todos)
	return todos
}