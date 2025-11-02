package main

import (
	"context"
	"fmt"
	"go-crud-todo-list/handler"
	"go-crud-todo-list/repository"
	"go-crud-todo-list/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize application
	if err := runApplication(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

// runApplication initializes and runs the todo application
func runApplication() error {
	fmt.Println("Go CRUD Todo List API")
	log.Println("Starting application initialization...")

	// Load configuration
	config, err := loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	log.Printf("Configuration loaded: port=%s, dataFile=%s", config.Port, config.DataFilePath)

	// Initialize data file if it doesn't exist
	if err := initializeDataFile(config.DataFilePath); err != nil {
		return fmt.Errorf("failed to initialize data file: %w", err)
	}

	// Initialize repository layer
	todoRepo := repository.NewFileBasedTodoRepository(config.DataFilePath)
	
	// Load existing data from file
	if err := todoRepo.Load(); err != nil {
		return fmt.Errorf("failed to load data from file: %w", err)
	}
	log.Println("Data loaded successfully")

	// Initialize service layer with repository dependency
	todoService := service.NewTodoService(todoRepo)
	log.Println("Service layer initialized")

	// Initialize handler layer with service dependency
	todoHandler := handler.NewTodoHandler(todoService)
	log.Println("Handler layer initialized")

	// Setup HTTP routes
	mux := todoHandler.SetupRoutes()
	log.Println("HTTP routes configured")

	// Configure HTTP server with proper timeouts
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", config.Port)
		log.Printf("API endpoints available at http://localhost:%s/todos", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Application started successfully")

	// Setup graceful shutdown
	setupGracefulShutdown(server, todoRepo)
	return nil
}

// Config holds application configuration
type Config struct {
	Port         string
	DataFilePath string
}

// loadConfiguration loads application configuration from environment variables
func loadConfiguration() (*Config, error) {
	config := &Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		DataFilePath: getEnvOrDefault("DATA_FILE", "todos.json"),
	}

	// Validate port
	if config.Port == "" {
		return nil, fmt.Errorf("port cannot be empty")
	}

	// Validate data file path
	if config.DataFilePath == "" {
		return nil, fmt.Errorf("data file path cannot be empty")
	}

	return config, nil
}

// initializeDataFile creates the data file if it doesn't exist
func initializeDataFile(filePath string) error {
	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("Data file already exists: %s", filePath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check data file status: %w", err)
	}

	// Create empty JSON structure for new file
	emptyStorage := `{
  "todos": [],
  "next_id": 1
}`

	// Create the file with initial empty structure
	if err := os.WriteFile(filePath, []byte(emptyStorage), 0644); err != nil {
		return fmt.Errorf("failed to create data file: %w", err)
	}

	log.Printf("Data file created: %s", filePath)
	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupGracefulShutdown handles graceful server shutdown on interrupt signals
func setupGracefulShutdown(server *http.Server, repo repository.TodoRepository) {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	
	// Register the channel to receive specific signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// Block until a signal is received
	sig := <-quit
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}

	// Save any pending data
	if err := repo.Save(); err != nil {
		log.Printf("Failed to save data during shutdown: %v", err)
	} else {
		log.Println("Data saved successfully during shutdown")
	}

	log.Println("Application shutdown complete")
}