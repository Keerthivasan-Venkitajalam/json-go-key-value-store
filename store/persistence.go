// Package store handles persistence to ensure data survives application restarts.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// DefaultFilePath specifies the default location of the persistent data file.
const DefaultFilePath = "./data/store.json"

// Store represents an in-memory key-value store with persistence capabilities.
type Store struct {
	data     map[string]string // In-memory data store
	filePath string            // Path to the JSON file for persistence
	mu       sync.RWMutex      // Mutex to ensure thread-safe access
}

// NewStore initializes a new Store instance with the given file path.
// If no file path is provided, it defaults to `DefaultFilePath`.
func NewStore(filePath string) *Store {
	if filePath == "" {
		filePath = DefaultFilePath
	}
	return &Store{
		data:     make(map[string]string),
		filePath: filePath,
	}
}

// Load loads the data from the JSON file into the store.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure the file exists; if not, start with an empty store.
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil
	}

	// Read the file contents.
	content, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON content into the store's data map.
	if err := json.Unmarshal(content, &s.data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// Save saves the current in-memory data to the JSON file.
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Marshal the in-memory data into JSON format.
	content, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Ensure the directory for the file exists.
	err = os.MkdirAll(filepath.Dir(s.filePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Write the JSON data to the file.
	err = os.WriteFile(s.filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Create adds a new key-value pair to the store.
func (s *Store) Create(key, value string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if key == "" {
        return errors.New("key cannot be empty")
    }

    if _, exists := s.data[key]; exists {
        return errors.New("key already exists")
    }

    if !isValidJSON(value) {
        return errors.New("invalid JSON format")
    }

    s.data[key] = value
    return nil
}

// Read retrieves the value for a given key.
func (s *Store) Read(key string) (string, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if key == "" {
        return "", errors.New("key cannot be empty")
    }

    value, exists := s.data[key]
    if !exists {
        return "", errors.New("key not found")
    }

    return value, nil
}

// Get retrieves the value for a given key.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// Update modifies the value for a given key.
func (s *Store) Update(key, value string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if key == "" {
        return errors.New("key cannot be empty")
    }

    if _, exists := s.data[key]; !exists {
        return errors.New("key not found")
    }

    if !isValidJSON(value) {
        return errors.New("invalid JSON format")
    }

    s.data[key] = value
    return nil
}

// Set sets a key-value pair in the store.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

// Delete removes a key-value pair from the store.
func (s *Store) Delete(key string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if key == "" {
        return errors.New("key cannot be empty")
    }

    if _, exists := s.data[key]; !exists {
        return errors.New("key not found")
    }

    delete(s.data, key)
    return nil
}

// Clear removes all key-value pairs from the store.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]string)
}

// // isValidJSON checks if a given string is a valid JSON object.
// func isValidJSON(data string) bool {
//     var js map[string]interface{}
//     return json.Unmarshal([]byte(data), &js) == nil
// }