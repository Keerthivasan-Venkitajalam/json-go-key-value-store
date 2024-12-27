package main

import (
    "fmt"
    "jsonkvstore/store"
)

func main() {
    // Initialize the JSON store
    jsonStore := store.NewJSONStore()

    // Example key-value pairs
    key1 := "user1"
    value1 := `{"name": "Alice", "age": 30}`

    key2 := "user2"
    value2 := `{"name": "Bob", "age": 25}`

    // Adding key-value pairs to the store
    fmt.Println("Adding key-value pairs to the store...")
    if err := jsonStore.Add(key1, value1); err != nil {
        fmt.Println("Error adding user1:", err)
    } else {
        fmt.Println("Added user1 successfully!")
    }

    if err := jsonStore.Add(key2, value2); err != nil {
        fmt.Println("Error adding user2:", err)
    } else {
        fmt.Println("Added user2 successfully!")
    }

    // Attempt to add a duplicate key
    fmt.Println("Attempting to add duplicate key...")
    if err := jsonStore.Add(key1, value1); err != nil {
        fmt.Println("Error:", err)
    }

    // Retrieve and print values
    fmt.Println("\nRetrieving values from the store...")
    if value, err := jsonStore.Get(key1); err != nil {
        fmt.Println("Error retrieving user1:", err)
    } else {
        fmt.Printf("Value for user1: %s\n", value)
    }

    if value, err := jsonStore.Get(key2); err != nil {
        fmt.Println("Error retrieving user2:", err)
    } else {
        fmt.Printf("Value for user2: %s\n", value)
    }

    // Update a key-value pair
    fmt.Println("\nUpdating a key-value pair...")
    updatedValue := `{"name": "Alice", "age": 31}`
    if err := jsonStore.Update(key1, updatedValue); err != nil {
        fmt.Println("Error updating user1:", err)
    } else {
        fmt.Println("Updated user1 successfully!")
    }

    // Delete a key-value pair
    fmt.Println("\nDeleting a key-value pair...")
    if err := jsonStore.Delete(key2); err != nil {
        fmt.Println("Error deleting user2:", err)
    } else {
        fmt.Println("Deleted user2 successfully!")
    }

    // Attempt to retrieve a deleted key
    fmt.Println("\nAttempting to retrieve a deleted key...")
    if _, err := jsonStore.Get(key2); err != nil {
        fmt.Println("Error:", err)
    }

    // Validate JSON examples
    fmt.Println("\nValidating JSON examples...")
    validJSON := `{"name": "Alice", "age": 30}`
    invalidJSON := `{"name": "Alice", age: 30}` // Missing quotes around the "age" key

    if err := store.ValidateJSON(validJSON); err != nil {
        fmt.Println("Error with valid JSON:", err)
    } else {
        fmt.Println("Valid JSON!")
    }

    if err := store.ValidateJSON(invalidJSON); err != nil {
        fmt.Println("Error with invalid JSON:", err)
    } else {
        fmt.Println("Invalid JSON!")
    }
}