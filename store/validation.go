// Package store handles JSON validation to ensure all inputs conform to the expected format.
// This file includes utility functions to validate JSON strings and check key-value constraints.
package store

import (
	"encoding/json" // To parse and validate JSON data
	"errors"        // For returning structured error messages
	"fmt"           // For formatted error messages
)

// ValidateJSON checks whether a string is valid JSON.
func ValidateJSON(input string) error {
	var temp interface{} // Temporary holder for parsed JSON data

	// Attempt to parse the JSON input
	if err := json.Unmarshal([]byte(input), &temp); err != nil {
		return fmt.Errorf("invalid JSON: %w", err) // Return error with context
	}

	return nil // Return nil if parsing succeeds
}

// ValidateKey ensures that the given key is non-empty and conforms to expected constraints.
func ValidateKey(key string) error {
	if key == "" {
		return errors.New("key cannot be empty") // Reject empty keys
	}

	if len(key) > 256 {
		return errors.New("key length exceeds 256 characters") // Enforce max length
	}

	return nil // Key is valid
}

// ValidateKeyValue ensures both the key and the value are valid.
func ValidateKeyValue(key string, value string) error {
	// Validate the key
	if err := ValidateKey(key); err != nil {
		return err // Return key-specific error
	}

	// Validate the value as JSON
	if err := ValidateJSON(value); err != nil {
		return errors.New("value is not a valid JSON object") // Return JSON-specific error
	}

	return nil // Both key and value are valid
}
