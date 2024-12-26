package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// TestCreate tests the creation of a new JSON object in the store
func TestCreate(t *testing.T) {
	store := NewStore()

	// Valid JSON
	validJSON := `{"name": "John", "age": 30}`
	err := store.Create("user1", validJSON)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Invalid JSON
	invalidJSON := `{"name": "John", "age": 30`
	err = store.Create("user2", invalidJSON)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, but got none")
	}
}

// TestRead tests the reading of a JSON object from the store
func TestRead(t *testing.T) {
	store := NewStore()

	// Creating a valid JSON object
	validJSON := `{"name": "John", "age": 30}`
	err := store.Create("user1", validJSON)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Reading the created JSON object
	result, err := store.Read("user1")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result != validJSON {
		t.Errorf("Expected %v, but got %v", validJSON, result)
	}

	// Trying to read a non-existent key
	_, err = store.Read("user2")
	if err == nil {
		t.Errorf("Expected error for non-existent key, but got none")
	}
}

// TestUpdate tests the updating of an existing JSON object
func TestUpdate(t *testing.T) {
	store := NewStore()

	// Creating a valid JSON object
	validJSON := `{"name": "John", "age": 30}`
	err := store.Create("user1", validJSON)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Updating the JSON object
	updatedJSON := `{"name": "John", "age": 31}`
	err = store.Update("user1", updatedJSON)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verifying the update
	result, err := store.Read("user1")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result != updatedJSON {
		t.Errorf("Expected %v, but got %v", updatedJSON, result)
	}

	// Trying to update a non-existent key
	err = store.Update("user2", updatedJSON)
	if err == nil {
		t.Errorf("Expected error for non-existent key, but got none")
	}
}

// TestDelete tests the deletion of a JSON object from the store
func TestDelete(t *testing.T) {
	store := NewStore()

	// Creating a valid JSON object
	validJSON := `{"name": "John", "age": 30}`
	err := store.Create("user1", validJSON)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Deleting the JSON object
	err = store.Delete("user1")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verifying the deletion
	_, err = store.Read("user1")
	if err == nil {
		t.Errorf("Expected error for deleted key, but got none")
	}
}

// TestPersistence tests the persistence of data in the store
func TestPersistence(t *testing.T) {
	store := NewStore()

	// Create a key-value pair
	validJSON := `{"name": "John", "age": 30}`
	err := store.Create("user1", validJSON)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Save to file
	err = store.Save()
	if err != nil {
		t.Errorf("Expected no error when saving, but got: %v", err)
	}

	// Create a new store instance and load data
	store2 := NewStore()
	err = store2.Load()
	if err != nil {
		t.Errorf("Expected no error when loading, but got: %v", err)
	}

	// Check if data is persisted correctly
	result, err := store2.Read("user1")
	if err != nil {
		t.Errorf("Expected no error when reading persisted data, but got: %v", err)
	}
	if result != validJSON {
		t.Errorf("Expected %v, but got %v", validJSON, result)
	}
}

// TestInvalidJSON tests invalid JSON scenarios
func TestInvalidJSON(t *testing.T) {
	store := NewStore()

	// Invalid JSON format
	invalidJSON := `{"name": "John", "age": }`
	err := store.Create("user1", invalidJSON)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, but got none")
	}

	// Valid JSON format but incorrect structure
	invalidStructure := `{"name": "John", "age": "invalid"}`
	err = store.Create("user2", invalidStructure)
	if err != nil {
		t.Errorf("Expected no error for valid structure, but got: %v", err)
	}
}

// TestEdgeCases tests various edge cases like empty strings or invalid keys
func TestEdgeCases(t *testing.T) {
	store := NewStore()

	// Empty JSON string
	err := store.Create("user1", "")
	if err == nil {
		t.Errorf("Expected error for empty JSON, but got none")
	}

	// Empty key
	err = store.Create("", `{"name": "John", "age": 30}`)
	if err == nil {
		t.Errorf("Expected error for empty key, but got none")
	}

	// Delete with empty key
	err = store.Delete("")
	if err == nil {
		t.Errorf("Expected error for empty key when deleting, but got none")
	}

	// Update with empty key
	err = store.Update("", `{"name": "John", "age": 30}`)
	if err == nil {
		t.Errorf("Expected error for empty key when updating, but got none")
	}
}

// Utility function to create a new store instance
func NewStore() *Store {
	return &Store{
		data:     make(map[string]string),
		filePath: "data/store.json",
	}
}
