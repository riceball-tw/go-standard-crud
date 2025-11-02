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