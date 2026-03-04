package handlers

import (
	"net/http"
)

/*
internal/handlers/handlers.go

Role: HTTP request and response processing.

Responsibilities:
1. Provide functions matching the `http.HandlerFunc` signature (w http.ResponseWriter, r *http.Request) for each route.
2. Extract and validate user input from URL parameters, query strings, and form bodies.
3. Call business logic functions (from internal/models) to fetch, create, update, or delete data based on the request.
4. Render HTML templates (from ui/html) or return JSON responses with appropriate HTTP status codes and headers based on the result.
*/

// Home is a stub handler for the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Literary Lions Forum!"))
}

// PostView is a stub handler for viewing a single post.
func PostView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific post..."))
}
