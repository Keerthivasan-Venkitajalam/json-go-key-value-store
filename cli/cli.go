// Package cli provides the implementation for the Command-Line Interface (CLI) for interacting with the JSON Key-Value Store.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"yourproject/store"
)

// RunCLI starts the Command-Line Interface for the JSON Key-Value Store.
func RunCLI() {
	fmt.Println("Welcome to the JSON Key-Value Store CLI!")
	fmt.Println("Type 'help' for a list of commands or 'exit' to quit.")

	// Create a scanner to read user input from the terminal
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Prompt the user for input
		fmt.Print("Enter command: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading input: %v\n", err)
			} else {
				fmt.Println("Input interrupted. Exiting...")
			}
			break
		}
		input := scanner.Text()

		// Exit condition
		if strings.ToLower(input) == "exit" {
			fmt.Println("Exiting... Goodbye!")
			break
		}

		// Split the input into command and arguments
		args := strings.Fields(input)
		if len(args) == 0 {
			fmt.Println("Invalid input. Please enter a command.")
			continue
		}

		// Command processing
		switch args[0] {
		case "create":
			// Handle JSON creation
			if len(args) < 3 {
				fmt.Println("Usage: create <key> <json>")
				continue
			}
			key := args[1]
			json := strings.Join(args[2:], " ") // Combine remaining args into JSON string
			err := store.CreateJSON(key, json)
			if err != nil {
				fmt.Printf("Error creating JSON: %v\n", err)
			} else {
				fmt.Printf("JSON with key '%s' created successfully!\n", key)
			}

		case "read":
			// Handle JSON reading
			if len(args) < 2 {
				fmt.Println("Usage: read <key>")
				continue
			}
			key := args[1]
			json, err := store.ReadJSON(key)
			if err != nil {
				fmt.Printf("Error reading JSON: %v\n", err)
			} else {
				fmt.Printf("JSON for key '%s':\n%s\n", key, json)
			}

		case "update":
			// Handle JSON update
			if len(args) < 3 {
				fmt.Println("Usage: update <key> <json>")
				continue
			}
			key := args[1]
			json := strings.Join(args[2:], " ") // Combine remaining args into JSON string
			err := store.UpdateJSON(key, json)
			if err != nil {
				fmt.Printf("Error updating JSON: %v\n", err)
			} else {
				fmt.Printf("JSON with key '%s' updated successfully!\n", key)
			}

		case "delete":
			// Handle JSON deletion
			if len(args) < 2 {
				fmt.Println("Usage: delete <key>")
				continue
			}
			key := args[1]
			err := store.DeleteJSON(key)
			if err != nil {
				fmt.Printf("Error deleting JSON: %v\n", err)
			} else {
				fmt.Printf("JSON with key '%s' deleted successfully!\n", key)
			}

		case "help":
			// Display CLI usage instructions
			fmt.Println("Available commands:")
			fmt.Println("  create <key> <json>   - Create a new JSON object.")
			fmt.Println("  read <key>            - Read a JSON object.")
			fmt.Println("  update <key> <json>   - Update an existing JSON object.")
			fmt.Println("  delete <key>          - Delete a JSON object.")
			fmt.Println("  exit                  - Exit the CLI.")

		default:
			fmt.Println("Invalid command. Type 'help' for a list of available commands.")
		}
	}
}
