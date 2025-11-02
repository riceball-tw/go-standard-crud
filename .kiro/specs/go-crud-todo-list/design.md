# Design Document

## Overview

The Go CRUD todo list application will be a REST API server built using Go's standard library. It will provide HTTP endpoints for managing todo items with JSON-based request/response handling and file-based persistence. The application follows a clean architecture pattern with separate layers for HTTP handling, business logic, and data storage.

## Architecture

The application uses a layered architecture:

```
┌─────────────────┐
│   HTTP Layer    │  ← Handles HTTP requests/responses, routing, JSON marshaling
├─────────────────┤
│  Service Layer  │  ← Business logic, validation, orchestration
├─────────────────┤
│ Repository Layer│  ← Data access, file I/O, persistence
└─────────────────┘
```

The server will run on a configurable port (default 8080) and handle concurrent requests using Go's built-in HTTP server capabilities.

## Components and Interfaces

### Todo Model
```go
type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Repository Interface
```go
type TodoRepository interface {
    GetAll() ([]Todo, error)
    GetByID(id int) (*Todo, error)
    Create(todo *Todo) error
    Update(id int, todo *Todo) error
    Delete(id int) error
    Save() error
    Load() error
}
```

### Service Interface
```go
type TodoService interface {
    GetAllTodos() ([]Todo, error)
    GetTodoByID(id int) (*Todo, error)
    CreateTodo(title, description string) (*Todo, error)
    UpdateTodo(id int, title, description string, completed bool) (*Todo, error)
    DeleteTodo(id int) error
}
```

### HTTP Handler
```go
type TodoHandler struct {
    service TodoService
}
```

## Data Models

### Todo Structure
- **ID**: Auto-incrementing integer primary key
- **Title**: Required string field (max 200 characters)
- **Description**: Optional string field (max 1000 characters)
- **Completed**: Boolean flag indicating completion status
- **CreatedAt**: Timestamp when todo was created
- **UpdatedAt**: Timestamp when todo was last modified

### Storage Format
Data will be stored as JSON in a local file (`todos.json`):
```json
{
  "todos": [
    {
      "id": 1,
      "title": "Learn Go",
      "description": "Complete Go tutorial",
      "completed": false,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "next_id": 2
}
```

## API Endpoints

| Method | Endpoint     | Description           | Request Body | Response |
|--------|-------------|-----------------------|--------------|----------|
| GET    | /todos       | Get all todos         | None         | Array of todos |
| GET    | /todos/{id}  | Get specific todo     | None         | Single todo |
| POST   | /todos       | Create new todo       | Todo JSON    | Created todo |
| PUT    | /todos/{id}  | Update existing todo  | Todo JSON    | Updated todo |
| DELETE | /todos/{id}  | Delete todo           | None         | 204 No Content |

### Request/Response Examples

**POST /todos**
```json
Request:
{
  "title": "Buy groceries",
  "description": "Milk, bread, eggs"
}

Response (201):
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, bread, eggs",
  "completed": false,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**PUT /todos/1**
```json
Request:
{
  "title": "Buy groceries",
  "description": "Milk, bread, eggs, cheese",
  "completed": true
}

Response (200):
{
  "id": 1,
  "title": "Buy groceries", 
  "description": "Milk, bread, eggs, cheese",
  "completed": true,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:30:00Z"
}
```

## Error Handling

### HTTP Status Codes
- **200 OK**: Successful GET/PUT operations
- **201 Created**: Successful POST operations
- **204 No Content**: Successful DELETE operations
- **400 Bad Request**: Invalid request data or malformed JSON
- **404 Not Found**: Todo item not found
- **500 Internal Server Error**: Server-side errors

### Error Response Format
```json
{
  "error": "Todo not found",
  "code": 404,
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### Error Scenarios
1. **Invalid JSON**: Return 400 with parsing error details
2. **Missing required fields**: Return 400 with validation error
3. **Todo not found**: Return 404 with appropriate message
4. **File I/O errors**: Return 500 with generic error message
5. **Invalid ID format**: Return 400 with format error

## Concurrency and Thread Safety

- Use mutex locks around file operations to prevent data corruption
- Implement proper locking in the repository layer for concurrent access
- Leverage Go's built-in HTTP server concurrency handling

## Testing Strategy

### Unit Tests
- **Repository Layer**: Test file I/O operations, data persistence, CRUD operations
- **Service Layer**: Test business logic, validation, error handling
- **Handler Layer**: Test HTTP request/response handling, JSON marshaling

### Integration Tests
- **End-to-End API Tests**: Test complete request flows through all layers
- **File Persistence Tests**: Verify data survives application restarts
- **Concurrent Access Tests**: Ensure thread safety under load

### Test Data Management
- Use temporary files for testing to avoid affecting production data
- Implement test fixtures for consistent test data
- Clean up test files after test completion

## Performance Considerations

- **File I/O Optimization**: Batch writes when possible, avoid frequent file operations
- **Memory Management**: Keep todo list in memory for fast reads, persist changes asynchronously if needed
- **Request Handling**: Leverage Go's efficient HTTP server for concurrent request processing
- **JSON Processing**: Use streaming JSON for large datasets if needed

## Security Considerations

- **Input Validation**: Sanitize and validate all user inputs
- **File Access**: Restrict file operations to designated data directory
- **Error Messages**: Avoid exposing internal system details in error responses
- **Resource Limits**: Implement reasonable limits on request size and todo content length