// Package handlers provides middleware for HTTP request processing such as authentication and logging.
package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware checks for a valid Authorization header in incoming requests and validates credentials.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// Ensure the Authorization header uses the Basic scheme
		if !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Unauthorized: Invalid Authentication Format", http.StatusUnauthorized)
			return
		}

		// Extract and decode Basic Authentication credentials
		credentials := strings.TrimPrefix(authHeader, "Basic ")
		username, password, ok := decodeBasicAuth(credentials)
		if !ok {
			http.Error(w, "Unauthorized: Invalid Authentication Format", http.StatusUnauthorized)
			return
		}

		// Validate credentials (for this example, using hardcoded username and password)
		expectedUsername := "admin"
		expectedPassword := "password123"

		if username != expectedUsername || password != expectedPassword {
			http.Error(w, "Unauthorized: Invalid Credentials", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if authentication succeeds
		next.ServeHTTP(w, r)
	})
}

// decodeBasicAuth decodes the base64-encoded username and password from the Authorization header.
func decodeBasicAuth(encoded string) (username, password string, ok bool) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// LoggingMiddleware logs details about incoming HTTP requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the HTTP method, request URI, and timestamp
		log.Printf("Received %s request for %s from %s at %s", r.Method, r.URL.Path, r.RemoteAddr, time.Now().Format(time.RFC3339))
		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
