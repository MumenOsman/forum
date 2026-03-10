package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"literary-lions-forum/internal/auth"
	"literary-lions-forum/internal/models"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
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

// TemplateData holds the dynamic data passed to HTML templates.
type TemplateData struct {
	CurrentYear     int
	Post            *models.Post
	Posts           []*models.Post
	Categories      []*models.Category
	IsAuthenticated bool
	User            *models.User
}

// Application holds the application-wide dependencies for the handlers.
type Application struct {
	Models        *models.AppModel
	TemplateCache map[string]*template.Template
}

// Render is a helper to render HTML templates.
func (app *Application) render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist", page), http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = &TemplateData{}
	}
	data.CurrentYear = time.Now().Year()

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewTemplateCache creates a cache of parsed templates.
func NewTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		if matches, err := filepath.Glob("./ui/html/*.layout.tmpl"); err == nil && len(matches) > 0 {
			ts, err = ts.ParseGlob("./ui/html/*.layout.tmpl")
			if err != nil {
				return nil, err
			}
		}

		if matches, err := filepath.Glob("./ui/html/partials/*.partial.tmpl"); err == nil && len(matches) > 0 {
			ts, err = ts.ParseGlob("./ui/html/partials/*.partial.tmpl")
			if err != nil {
				return nil, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}

// Home handles requests to the root URL ("/").
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Placeholder for fetching posts and categories
	// posts, err := app.Models.GetLatestPosts()
	// categories, err := app.Models.GetAllCategories()

	app.render(w, http.StatusOK, "home.page.tmpl", &TemplateData{
		Posts: []*models.Post{}, // Empty for now
	})
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

// UserSignup handles requests to the user signup page ("/user/signup").
func (app *Application) UserSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.render(w, http.StatusOK, "signup.page.tmpl", nil)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, err := auth.GenerateSessionID() // Reusing to generate a UUID for the user ID
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = app.Models.InsertUser(userID, email, username, hashedPassword)
	if err != nil {
		// Possibly email/username already taken
		http.Error(w, "Error saving user properties or user already exists", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// UserLogin handles requests to the user login page ("/user/login").
func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.render(w, http.StatusOK, "login.page.tmpl", nil)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	userID, err := app.Models.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	sessionID, err := auth.GenerateSessionID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = app.Models.InsertSession(sessionID, userID, expiresAt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// UserLogout handles requests to log out the user ("/user/logout").
func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		_ = app.Models.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
