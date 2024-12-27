# JSON Key-Value Store

This project implements a JSON key-value store with a command-line interface (CLI) and HTTP API. The store supports basic operations such as adding, retrieving, updating, and deleting key-value pairs. The data is persisted to a JSON file to ensure it survives application restarts.

## Features

- Add, retrieve, update, and delete key-value pairs
- Validate JSON data
- Command-line interface (CLI) for interacting with the store
- HTTP API for remote access
- Basic authentication middleware
- Logging of actions

## Project Structure

- `main.go`: Entry point of the application
- `store/`: Contains the core logic for the key-value store and persistence
  - `store.go`: Core logic for the in-memory store
  - `persistence.go`: Persistence logic to save and load data from a JSON file
  - `validation.go`: Validation functions for keys and JSON data
  - `store_test.go`: Unit tests for the store
- `handlers/`: Contains HTTP handlers and middleware
  - `handlers.go`: HTTP handlers for the API
  - `middleware.go`: Middleware for authentication and logging
  - `handlers_test.go`: Unit tests for the handlers
- `cli/`: Contains the CLI implementation
  - `cli.go`: CLI logic for interacting with the store
- `logs/`: Directory for log files
  - `actions.log`: Log file for actions
- `data/`: Directory for the JSON data file
  - `store.json`: JSON file for storing key-value pairs
- `go.mod`: Go module file

## Prerequisites

- Go 1.23.4 or later

## Setup

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/json-go-key-value-store.git
   cd json-go-key-value-store
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Running the Application

```sh
go run main.go
```

## Logging

Actions are logged to `logs/actions.log`. Ensure the `logs` directory exists or is created by the application.
