package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Storage handles storing command history in SQLite
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new storage instance with database in ~/.gema/
func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	gemaDir := filepath.Join(homeDir, ".gema")
	if err := os.MkdirAll(gemaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", gemaDir, err)
	}

	dbPath := filepath.Join(gemaDir, "gema.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS command_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		input TEXT NOT NULL,
		response TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Storage{db: db}, nil
}

// StoreCommand stores a command input and response in the database
func (s *Storage) StoreCommand(input, response string) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	result, err := s.db.Exec("INSERT INTO command_history (input, response) VALUES (?, ?)", input, response)
	if err != nil {
		return fmt.Errorf("failed to store command: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expected to insert 1 row, but inserted %d", rowsAffected)
	}

	return nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	if s.db == nil {
		return nil // No connection to close
	}
	return s.db.Close()
}

// StoreCommandHistory is a facade function that handles database operations internally
func StoreCommandHistory(input, response string) error {
	// Create a new storage instance
	storage, err := NewStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	defer func() {
		if closeErr := storage.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Error closing database: %v\n", closeErr)
		}
	}()

	// Store the command
	if err := storage.StoreCommand(input, response); err != nil {
		return fmt.Errorf("failed to store command: %w", err)
	}

	return nil
}
