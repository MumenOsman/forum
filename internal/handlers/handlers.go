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
	ErrorMessage    string
	SearchQuery     string
	Profile         *models.User
	Comments        []*models.Comment
	ErrorCode       int
}

// Application holds the application-wide dependencies for the handlers.
type Application struct {
	Models        *models.AppModel
	TemplateCache map[string]*template.Template
}

// serverError sends a generic 500 error to the user gracefully.
func (app *Application) serverError(w http.ResponseWriter, _ error) {
	app.render(w, http.StatusInternalServerError, "error.page.tmpl", &TemplateData{
		ErrorMessage: "Internal Server Error. We are looking into it.",
	})
}

// clientError sends an error message and status code to the user gracefully.
func (app *Application) clientError(w http.ResponseWriter, status int, message string) {
	app.render(w, status, "error.page.tmpl", &TemplateData{
		ErrorMessage: message,
	})
}

// notFound sends a 404 error gracefully
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound, "The requested resource could not be found.")
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

// getAuthenticatedUserID is a helper that reads the session cookie to return the UserID, or empty string.
func (app *Application) getAuthenticatedUserID(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	user, err := app.Models.GetUserBySession(cookie.Value)
	if err != nil {
		return ""
	}
	return user.ID
}

// Home handles requests to the root URL ("/").
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	userID := app.getAuthenticatedUserID(r)

	searchQuery := r.URL.Query().Get("q")
	categoryID := r.URL.Query().Get("category")
	authored := r.URL.Query().Get("authored")
	liked := r.URL.Query().Get("liked")

	var authoredBy, likedBy string
	if authored == "true" {
		if userID == "" {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		authoredBy = userID
	}
	if liked == "true" {
		if userID == "" {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		likedBy = userID
	}

	posts, err := app.Models.GetFilteredPosts(categoryID, authoredBy, likedBy, searchQuery)
	if err != nil {
		app.serverError(w, nil)
		return
	}

	categories, err := app.Models.GetAllCategories()
	if err != nil {
		app.serverError(w, nil)
		return
	}

	app.render(w, http.StatusOK, "home.page.tmpl", &TemplateData{
		Posts:           posts,
		Categories:      categories,
		IsAuthenticated: userID != "",
	})
}

// PostView handles requests to view a specific post (e.g., "/post/view?id=...").
func (app *Application) PostView(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		app.notFound(w)
		return
	}

	post, err := app.Models.GetPostByID(idStr)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, nil)
		}
		return
	}

	userID := app.getAuthenticatedUserID(r)
	categories, _ := app.Models.GetAllCategories()
	app.render(w, http.StatusOK, "view.page.tmpl", &TemplateData{
		Post:            post,
		Categories:      categories,
		IsAuthenticated: userID != "",
	})
}

// PostCreate handles requests to create a new post ("/post/create").
func (app *Application) PostCreate(w http.ResponseWriter, r *http.Request) {
	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		categories, err := app.Models.GetAllCategories()
		if err != nil {
			app.serverError(w, nil)
			return
		}

		app.render(w, http.StatusOK, "create.page.tmpl", &TemplateData{
			Categories:      categories,
			IsAuthenticated: true,
		})
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	categoryStrings := r.PostForm["categories"]

	var categories []int
	for _, str := range categoryStrings {
		catID, _ := strconv.Atoi(str)
		categories = append(categories, catID)
	}

	if title == "" || content == "" || len(categories) == 0 {
		app.clientError(w, http.StatusBadRequest, "Bad request: title, content, and at least 1 category required")
		return
	}

	postID, err := auth.GenerateSessionID() // secure UUID for post ID
	if err != nil {
		app.serverError(w, nil)
		return
	}

	err = app.Models.InsertPost(postID, userID, title, content, categories)
	if err != nil {
		app.serverError(w, nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CommentCreate handles requests to create a new comment ("/comment/create").
func (app *Application) CommentCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	postID := r.PostForm.Get("post_id")
	content := r.PostForm.Get("content")

	if postID == "" || content == "" {
		app.clientError(w, http.StatusBadRequest, "Bad request: post_id and content required")
		return
	}

	commentID, err := auth.GenerateSessionID() // secure UUID for comment ID
	if err != nil {
		app.serverError(w, nil)
		return
	}

	err = app.Models.InsertComment(commentID, postID, userID, content)
	if err != nil {
		app.serverError(w, nil)
		return
	}

	http.Redirect(w, r, "/post/view?id="+postID, http.StatusSeeOther)
}

// UserSignup handles requests to the user signup page ("/user/signup").
func (app *Application) UserSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.render(w, http.StatusOK, "signup.page.tmpl", nil)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	username := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		app.serverError(w, nil)
		return
	}

	userID, err := auth.GenerateSessionID() // Reusing to generate a UUID for the user ID
	if err != nil {
		app.serverError(w, nil)
		return
	}

	err = app.Models.InsertUser(userID, email, username, hashedPassword)
	if err != nil {
		// Possibly email/username already taken
		app.clientError(w, http.StatusInternalServerError, "Error saving user properties or user already exists")
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
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	userID, err := app.Models.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.render(w, http.StatusUnauthorized, "login.page.tmpl", &TemplateData{
				ErrorMessage: "Invalid email or password",
			})
		} else {
			app.serverError(w, nil)
		}
		return
	}

	sessionID, err := auth.GenerateSessionID()
	if err != nil {
		app.serverError(w, nil)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = app.Models.InsertSession(sessionID, userID, expiresAt)
	if err != nil {
		app.serverError(w, nil)
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
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
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

// VoteHandler handles requests to upvote or downvote a post or comment ("/vote").
func (app *Application) VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	targetID := r.PostForm.Get("target_id")
	targetType := r.PostForm.Get("target_type")
	voteTypeStr := r.PostForm.Get("vote_type")

	if targetID == "" || (targetType != "post" && targetType != "comment") || (voteTypeStr != "1" && voteTypeStr != "-1") {
		app.clientError(w, http.StatusBadRequest, "Bad Request: invalid vote parameters")
		return
	}

	voteType, err := strconv.Atoi(voteTypeStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	err = app.Models.InsertOrUpdateVote(userID, targetID, targetType, voteType)
	if err != nil {
		app.serverError(w, nil)
		return
	}

	// Redirect back to the page the user came from (or fallback to Home)
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// UserProfile handles requests to view a user's profile ("/user/profile?id=...").
func (app *Application) UserProfile(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		app.notFound(w)
		return
	}

	user, err := app.Models.GetUserByID(idStr)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	posts, err := app.Models.GetFilteredPosts("", idStr, "", "")
	if err != nil {
		app.serverError(w, err)
		return
	}

	currentUserID := app.getAuthenticatedUserID(r)
	categories, _ := app.Models.GetAllCategories()
	app.render(w, http.StatusOK, "profile.page.tmpl", &TemplateData{
		User:            user,
		Posts:           posts,
		Categories:      categories,
		IsAuthenticated: currentUserID != "",
	})
}
