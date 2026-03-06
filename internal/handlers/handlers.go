package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"literary-lions-forum/internal/models"
)

/*
internal/handlers/handlers.go

Role: HTTP request and response processing.

Responsibilities:
1. Provide methods matching the `http.HandlerFunc` signature for each route attached to the Application struct.
2. Extract and validate user input from URL parameters, query strings, and form bodies.
3. Call business logic functions (from internal/models) to fetch, create, update, or delete data based on the request.
4. Render HTML templates (from ui/html) or return JSON responses with appropriate HTTP status codes and headers based on the result.
*/

// Application holds the application-wide dependencies for the handlers.
type Application struct {
	Models *models.AppModel
}

// Home handles requests to the root URL ("/").
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	// Restrict to exactly "/" so it doesn't act as a catch-all for 404s.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Placeholder: In the future, we will fetch posts here using app.Models.GetLatestPosts()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Literary Lions Forum! (Home Page)"))
}

// PostView handles requests to view a specific post (e.g., "/post/view?id=1").
func (app *Application) PostView(w http.ResponseWriter, r *http.Request) {
	// Extract the 'id' parameter from the query string
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Placeholder: Query the database for the post
	// post, err := app.Models.GetPost(fmt.Sprintf("%d", id))
	// if err != nil { ... }

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Displaying a specific post with ID %d...", id)
}

// PostCreate handles requests to create a new post ("/post/create").
func (app *Application) PostCreate(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Placeholder for parsing form data and inserting into DB
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Saving a new post..."))
}
