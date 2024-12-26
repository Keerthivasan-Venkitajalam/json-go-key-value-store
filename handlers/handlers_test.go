// Package handlers provides unit tests for HTTP handlers, including authentication and logging.
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// TestLoggingMiddleware tests the logging functionality of the LoggingMiddleware.
func TestLoggingMiddleware(t *testing.T) {
    // Create a request to test with
    req, err := http.NewRequest(http.MethodGet, "/get/testuser", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    // Create a response recorder to capture the response
    rr := httptest.NewRecorder()

    // Create a simple handler for testing
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Wrap the handler with the LoggingMiddleware
    middleware := LoggingMiddleware(handler)

    // Serve the request through the middleware
    middleware.ServeHTTP(rr, req)

    // Check if the status code is correct
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
    }
}

// TestAuthMiddleware tests the authentication functionality of the AuthMiddleware.
func TestAuthMiddleware(t *testing.T) {
    tests := []struct {
        name         string
        authHeader   string
        expectedCode int
    }{
        {
            name:         "Valid Credentials",
            authHeader:   "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:password123")),
            expectedCode: http.StatusOK,
        },
        {
            name:         "Invalid Credentials",
            authHeader:   "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword")),
            expectedCode: http.StatusUnauthorized,
        },
        {
            name:         "No Credentials",
            authHeader:   "",
            expectedCode: http.StatusUnauthorized,
        },
    }

    // Iterate through test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a request to test with
            req, err := http.NewRequest(http.MethodGet, "/get/testuser", nil)
            if err != nil {
                t.Fatalf("Failed to create request: %v", err)
            }

            // Set the Authorization header if provided
            if tt.authHeader != "" {
                req.Header.Set("Authorization", tt.authHeader)
            }

            // Create a response recorder to capture the response
            rr := httptest.NewRecorder()

            // Create a simple handler for testing
            handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            })

            // Wrap the handler with the AuthMiddleware
            middleware := AuthMiddleware(handler)

            // Serve the request through the middleware
            middleware.ServeHTTP(rr, req)

            // Check if the status code matches the expected value
            if rr.Code != tt.expectedCode {
                t.Errorf("Test case '%s': Expected status code %d, but got %d", tt.name, tt.expectedCode, rr.Code)
            }
        })
    }
}