# Go CRUD Todo List API

A simple REST API for managing todos built with Go. This application provides full CRUD (Create, Read, Update, Delete) operations for todo items with file-based persistence.

## Features

- ✅ Create, read, update, and delete todos
- ✅ JSON-based REST API
- ✅ File-based data persistence
- ✅ Input validation and error handling
- ✅ Graceful shutdown
- ✅ Configurable via environment variables
- ✅ Comprehensive test coverage

## Prerequisites

- Go 1.19 or higher
- curl (for testing API endpoints)

## How to Build the App

### 1. Clone or Download the Project
```bash
# If using git
git clone <repository-url>
cd go-crud-todo-list

# Or download and extract the source code
```

### 2. Build the Executable
```bash
# Build for your current platform
go build -o todo-app

# Build for specific platforms (optional)
# For Linux
GOOS=linux GOARCH=amd64 go build -o todo-app-linux

# For Windows
GOOS=windows GOARCH=amd64 go build -o todo-app.exe

# For macOS
GOOS=darwin GOARCH=amd64 go build -o todo-app-macos
```

The executable will be created in the current directory.

## How to Run the Executable

### Basic Usage
```bash
# Run with default settings (port 8080, todos.json file)
./todo-app
```

### With Custom Configuration
```bash
# Use a different port
PORT=3000 ./todo-app

# Use a different data file
DATA_FILE=my-todos.json ./todo-app

# Use both custom port and data file
PORT=3000 DATA_FILE=my-todos.json ./todo-app
```

### Expected Output
When the app starts successfully, you'll see:
```
Go CRUD Todo List API
Starting application initialization...
Configuration loaded: port=8080, dataFile=todos.json
Data file created: todos.json
Data loaded successfully
Service layer initialized
Handler layer initialized
HTTP routes configured
Server listening on port 8080
API endpoints available at http://localhost:8080/todos
Application started successfully
```

### Stopping the App
Press `Ctrl+C` to stop the server gracefully. The app will save any pending data before shutting down.

## How to Interact with the App

The API provides the following endpoints:

### 1. Get All Todos
```bash
curl http://localhost:8080/todos
```
**Response:** Array of todo objects

### 2. Get Todo by ID
```bash
curl http://localhost:8080/todos/1
```
**Response:** Single todo object or 404 if not found

### 3. Create a New Todo
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread"}'
```
**Response:** Created todo object with assigned ID

### 4. Update an Existing Todo
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread, cheese", "completed": true}'
```
**Response:** Updated todo object

### 5. Delete a Todo
```bash
curl -X DELETE http://localhost:8080/todos/1
```
**Response:** 204 No Content on success

### Todo Object Structure
```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2023-11-02T10:30:00Z",
  "updated_at": "2023-11-02T10:30:00Z"
}
```

### Example Usage Flow
```bash
# 1. Check initial empty state
curl http://localhost:8080/todos

# 2. Create first todo
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Complete CRUD tutorial"}'

# 3. Create second todo
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Write tests", "description": "Add unit tests"}'

# 4. View all todos
curl http://localhost:8080/todos

# 5. Mark first todo as completed
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Complete CRUD tutorial", "completed": true}'

# 6. Delete second todo
curl -X DELETE http://localhost:8080/todos/2

# 7. Verify final state
curl http://localhost:8080/todos
```

## How to Test the App

### 1. Run Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 2. Manual API Testing

#### Start the Application
```bash
./todo-app
```

#### Test Each Endpoint
```bash
# Test 1: Get empty todos list
echo "Testing GET /todos (empty)..."
curl -s http://localhost:8080/todos | jq .

# Test 2: Create a todo
echo "Testing POST /todos..."
curl -s -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Todo", "description": "This is a test"}' | jq .

# Test 3: Get todos list (should have 1 item)
echo "Testing GET /todos (with data)..."
curl -s http://localhost:8080/todos | jq .

# Test 4: Get specific todo
echo "Testing GET /todos/1..."
curl -s http://localhost:8080/todos/1 | jq .

# Test 5: Update todo
echo "Testing PUT /todos/1..."
curl -s -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Todo", "description": "This is updated", "completed": true}' | jq .

# Test 6: Delete todo
echo "Testing DELETE /todos/1..."
curl -s -X DELETE http://localhost:8080/todos/1 -w "Status: %{http_code}\n"

# Test 7: Verify deletion
echo "Testing GET /todos (should be empty again)..."
curl -s http://localhost:8080/todos | jq .
```

#### Test Error Cases
```bash
# Test invalid JSON
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Invalid JSON"' # Missing closing brace

# Test missing title
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"description": "No title provided"}'

# Test non-existent todo
curl http://localhost:8080/todos/999

# Test invalid ID format
curl http://localhost:8080/todos/abc
```

### 3. Automated Test Script

Create a test script `test_api.sh`:
```bash
#!/bin/bash

# Start the app in background
./todo-app &
APP_PID=$!

# Wait for app to start
sleep 2

# Run tests
echo "Running API tests..."

# Test creating and retrieving todos
RESPONSE=$(curl -s -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Todo", "description": "Test Description"}')

if echo "$RESPONSE" | grep -q "Test Todo"; then
    echo "✅ Create todo test passed"
else
    echo "❌ Create todo test failed"
fi

# Clean up
kill $APP_PID
echo "Tests completed"
```

Make it executable and run:
```bash
chmod +x test_api.sh
./test_api.sh
```

## Configuration

| Environment Variable | Default Value | Description |
|---------------------|---------------|-------------|
| `PORT` | `8080` | Port number for the HTTP server |
| `DATA_FILE` | `todos.json` | Path to the JSON file for data persistence |

## Data Persistence

- Todos are stored in a JSON file (default: `todos.json`)
- The file is created automatically on first run
- Data is saved immediately after each operation
- Data is also saved during graceful shutdown

## Error Handling

The API returns appropriate HTTP status codes:
- `200 OK` - Successful GET/PUT operations
- `201 Created` - Successful POST operations
- `204 No Content` - Successful DELETE operations
- `400 Bad Request` - Invalid input or malformed JSON
- `404 Not Found` - Todo not found
- `405 Method Not Allowed` - Unsupported HTTP method
- `500 Internal Server Error` - Server-side errors

## Project Structure

```
.
├── main.go                 # Application entry point and server setup
├── models/
│   ├── todo.go            # Todo model and validation
│   └── storage.go         # In-memory storage management
├── repository/
│   └── todo_repository.go # Data persistence layer
├── service/
│   └── todo_service.go    # Business logic layer
├── handler/
│   └── todo_handler.go    # HTTP request handling
├── tests/
│   └── *_test.go          # Unit tests
├── todos.json             # Data file (created at runtime)
└── README.md              # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).