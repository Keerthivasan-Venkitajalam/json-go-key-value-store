// Package store is responsible for managing all the core logic related to our key-value JSON store.
// This includes thread-safe operations for adding, retrieving, updating, and deleting JSON objects.
package store

import (
	"encoding/json" // To handle JSON operations like marshaling and unmarshaling
	"errors"        // To manage errors in a structured way
	"sync"          // To ensure thread-safe access using mutexes
)

// JSONStore is the primary data structure that holds our in-memory store.
// The `data` map stores key-value pairs where both the key and value are strings.
// A mutex (`mu`) ensures thread safety for concurrent access.
type JSONStore struct {
	data map[string]string // The in-memory key-value storage
	mu   sync.RWMutex      // Read-write mutex for thread-safe access
}

// NewJSONStore initializes and returns a new instance of JSONStore.
// This is the entry point for creating a fresh store in memory.
func NewJSONStore() *JSONStore {
	return &JSONStore{
		data: make(map[string]string), // Create an empty map for storing data
	}
}

// Add inserts a new key-value pair into the store.
// If the key already exists or the JSON is invalid, it returns an error.
func (s *JSONStore) Add(key, jsonData string) error {
	s.mu.Lock()         // Lock the mutex to prevent simultaneous writes
	defer s.mu.Unlock() // Unlock it after the operation completes

	if key == "" {
		return errors.New("key cannot be empty") // Validate key
	}

	// Check if the key already exists in the store
	if _, exists := s.data[key]; exists {
		return errors.New("key already exists") // Return an error if the key is a duplicate
	}

	// Validate the provided JSON data
	if !isValidJSON(jsonData) {
		return errors.New("invalid JSON format") // Ensure the JSON is well-formed
	}

	// Add the key-value pair to the store
	s.data[key] = jsonData
	return nil // Return nil to indicate success
}

// Get retrieves the value associated with a given key from the store.
// If the key doesn't exist, it returns an error.
func (s *JSONStore) Get(key string) (string, error) {
	s.mu.RLock()         // Use a read-lock for safe concurrent reads
	defer s.mu.RUnlock() // Release the lock when done

	if key == "" {
		return "", errors.New("key cannot be empty") // Validate key
	}

	// Check if the key exists in the store
	value, exists := s.data[key]
	if !exists {
		return "", errors.New("key not found") // Return an error if the key is missing
	}

	return value, nil // Return the associated value
}

// Update modifies the value associated with an existing key.
// It validates the new JSON and returns an error if the key doesn't exist or JSON is invalid.
func (s *JSONStore) Update(key, newJSONData string) error {
	s.mu.Lock()         // Lock the mutex for write operations
	defer s.mu.Unlock() // Unlock it afterward

	if key == "" {
		return errors.New("key cannot be empty") // Validate key
	}

	// Check if the key exists in the store
	if _, exists := s.data[key]; !exists {
		return errors.New("key not found") // Cannot update a non-existent key
	}

	// Validate the new JSON data
	if !isValidJSON(newJSONData) {
		return errors.New("invalid JSON format") // Ensure the new JSON is well-formed
	}

	// Update the key with the new value
	s.data[key] = newJSONData
	return nil // Return nil to indicate success
}

// Delete removes a key-value pair from the store.
// If the key is not found, it returns an error.
func (s *JSONStore) Delete(key string) error {
	s.mu.Lock()         // Lock the mutex to ensure safe write access
	defer s.mu.Unlock() // Unlock it after the operation

	if key == "" {
		return errors.New("key cannot be empty") // Validate key
	}

	// Check if the key exists in the store
	if _, exists := s.data[key]; !exists {
		return errors.New("key not found") // Cannot delete a key that doesn't exist
	}

	// Remove the key-value pair from the store
	delete(s.data, key)
	return nil // Return nil to indicate success
}

// isValidJSON checks if a given string is a valid JSON object.
// This is a helper function used internally for JSON validation.
func isValidJSON(data string) bool {
	var js map[string]interface{} // Create a temporary map to unmarshal JSON
	return json.Unmarshal([]byte(data), &js) == nil // Return true if unmarshaling succeeds
}
