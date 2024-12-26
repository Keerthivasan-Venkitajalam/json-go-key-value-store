// Package handlers implements HTTP route handlers to interact with the JSON Key-Value Store.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"json-key-value-store/store"
)

// Response represents a consistent structure for API responses.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Optional field for response data
}

// CreateKeyValueHandler handles the creation of new key-value pairs in the JSON store.
func CreateKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	key := requestData["key"]
	value := requestData["value"]
	if key == "" || value == "" {
		http.Error(w, "Key and value are required fields", http.StatusBadRequest)
		return
	}

	// Store the key-value pair
	if err := store.Create(key, value); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create key-value pair: %s", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := Response{Message: "Key-value pair created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReadKeyValueHandler retrieves a key-value pair by its key from the store.
func ReadKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing 'key' parameter", http.StatusBadRequest)
		return
	}

	// Retrieve the value
	value, err := store.Read(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Key not found: %s", err), http.StatusNotFound)
		return
	}

	// Send success response
	response := Response{Message: "Key-value pair retrieved", Data: value}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateKeyValueHandler updates the value of an existing key in the store.
func UpdateKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	key := requestData["key"]
	value := requestData["value"]
	if key == "" || value == "" {
		http.Error(w, "Key and value are required fields", http.StatusBadRequest)
		return
	}

	// Update the key-value pair
	if err := store.Update(key, value); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update key-value pair: %s", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := Response{Message: "Key-value pair updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteKeyValueHandler deletes a key-value pair from the store by its key.
func DeleteKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing 'key' parameter", http.StatusBadRequest)
		return
	}

	// Delete the key-value pair
	if err := store.Delete(key); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete key-value pair: %s", err), http.StatusNotFound)
		return
	}

	// Send success response
	response := Response{Message: "Key-value pair deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetupRoutes initializes the HTTP server routes.
func SetupRoutes() {
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/create", CreateKeyValueHandler)
	mux.HandleFunc("/read", ReadKeyValueHandler)
	mux.HandleFunc("/update", UpdateKeyValueHandler)
	mux.HandleFunc("/delete", DeleteKeyValueHandler)

	// Wrap with middleware and start the server
	wrappedMux := AuthMiddleware(LoggingMiddleware(mux))
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
