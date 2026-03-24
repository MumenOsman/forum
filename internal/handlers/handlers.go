package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"literary-lions-forum/internal/auth"
	"literary-lions-forum/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	AuthenticatedUser *models.User
	User            *models.User
	ErrorMessage    string
	SuccessMessage  string
	SearchQuery     string
	Profile         *models.User
	Comments        []*models.Comment
	ErrorCode       int
	FilterAuthored  bool
	FilterLiked     bool
	FilterCategory  string
	// DM fields
	Messages       []*models.Message
	Conversations  []*models.ConversationPreview
	OtherUser      *models.User
	UnreadCount    int
}

// Application holds the application-wide dependencies for the handlers.
type Application struct {
	Models        *models.AppModel
	TemplateCache map[string]*template.Template
}

// serverError sends a generic 500 error to the user gracefully.
func (app *Application) serverError(w http.ResponseWriter, err error) {
	log.Printf("Server Error: %v", err)
	app.render(w, http.StatusInternalServerError, "error.page.tmpl", &TemplateData{
		ErrorMessage: "Internal Server Error. We are looking into it.",
		ErrorCode:    http.StatusInternalServerError,
	})
}

// clientError sends an error message and status code to the user gracefully.
func (app *Application) clientError(w http.ResponseWriter, status int, message string) {
	app.render(w, status, "error.page.tmpl", &TemplateData{
		ErrorMessage: message,
		ErrorCode:    status,
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

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		// Prevent infinite loop if error.page.tmpl itself fails
		if page == "error.page.tmpl" {
			http.Error(w, "Critical Error: Cannot render error page.", http.StatusInternalServerError)
			return
		}
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
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

func (app *Application) getAuthenticatedUser(r *http.Request) *models.User {
	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		return nil
	}
	user, err := app.Models.GetUserByID(userID)
	if err != nil {
		return nil
	}
	return user
}

// Home handles requests to the root URL ("/").
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
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
		app.serverError(w, err)
		return
	}

	categories, err := app.Models.GetAllCategories()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.page.tmpl", &TemplateData{
		Posts:             posts,
		Categories:        categories,
		IsAuthenticated:   userID != "",
		AuthenticatedUser: app.getAuthenticatedUser(r),
		SearchQuery:       searchQuery,
		FilterAuthored:    authored == "true",
		FilterLiked:       liked == "true",
		FilterCategory:    categoryID,
	})
}

// PostView handles requests to view a specific post (e.g., "/post/view?id=...").
func (app *Application) PostView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

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
			app.serverError(w, err)
		}
		return
	}

	userID := app.getAuthenticatedUserID(r)
	categories, _ := app.Models.GetAllCategories()
	app.render(w, http.StatusOK, "view.page.tmpl", &TemplateData{
		Post:              post,
		Comments:          post.Comments,
		Categories:        categories,
		IsAuthenticated:   userID != "",
		AuthenticatedUser: app.getAuthenticatedUser(r),
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
			app.serverError(w, err)
			return
		}

		app.render(w, http.StatusOK, "create.page.tmpl", &TemplateData{
			Categories:        categories,
			IsAuthenticated:   true,
			AuthenticatedUser: app.getAuthenticatedUser(r),
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
		app.serverError(w, err)
		return
	}

	err = app.Models.InsertPost(postID, userID, title, content, categories)
	if err != nil {
		app.serverError(w, err)
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

	commentID, err := app.Models.InsertComment(postID, userID, content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// AJAX Support
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		comment, err := app.Models.GetCommentByID(commentID)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"comment": comment,
			})
			return
		}
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

	if len(password) < 6 {
		app.render(w, http.StatusUnprocessableEntity, "signup.page.tmpl", &TemplateData{
			ErrorMessage: "Password must be at least 6 characters.",
		})
		return
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, err := auth.GenerateSessionID() // Reusing to generate a UUID for the user ID
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.Models.InsertUser(userID, email, username, hashedPassword)
	if err != nil {
		// Check if this is a duplicate email/username (UNIQUE constraint violation)
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			app.render(w, http.StatusConflict, "signup.page.tmpl", &TemplateData{
				ErrorMessage: "An account with that email or username already exists.",
			})
			return
		}
		app.serverError(w, err)
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
			app.serverError(w, err)
		}
		return
	}

	sessionID, err := auth.GenerateSessionID()
	if err != nil {
		app.serverError(w, err)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = app.Models.InsertSession(sessionID, userID, expiresAt)
	if err != nil {
		app.serverError(w, err)
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
		app.serverError(w, err)
		return
	}

	// AJAX Support
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		counts, err := app.Models.GetVoteCounts(targetID, targetType)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success":  true,
				"likes":    counts.Likes,
				"dislikes": counts.Dislikes,
			})
			return
		}
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
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

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

	app.render(w, http.StatusOK, "profile.page.tmpl", &TemplateData{
		User:              user,
		Posts:             posts,
		IsAuthenticated:   app.getAuthenticatedUserID(r) != "",
		AuthenticatedUser: app.getAuthenticatedUser(r),
	})
}

// ProfileEdit handles GET and POST requests to edit a user's profile.
func (app *Application) ProfileEdit(w http.ResponseWriter, r *http.Request) {
	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	user, err := app.Models.GetUserByID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if r.Method == http.MethodGet {
		app.render(w, http.StatusOK, "profile_edit.page.tmpl", &TemplateData{
			User:              user,
			IsAuthenticated:   true,
			AuthenticatedUser: user,
		})
		return
	}

	if r.Method == http.MethodPost {
		// Limit the size of the uploaded file to 2MB
		err := r.ParseMultipartForm(2 << 20)
		if err != nil {
			app.clientError(w, http.StatusBadRequest, "File too large")
			return
		}

		// Update username if provided
		newUsername := r.PostFormValue("username")
		if newUsername != "" && newUsername != user.Username {
			err = app.Models.UpdateUserNameUsername(userID, newUsername)
			if err != nil {
				// Likely a duplicate username
				user, _ = app.Models.GetUserByID(userID)
				app.render(w, http.StatusOK, "profile_edit.page.tmpl", &TemplateData{
					User:              user,
					IsAuthenticated:   true,
					AuthenticatedUser: user,
					ErrorMessage:      "Username is already taken or invalid.",
				})
				return
			}
		}

		aboutMe := r.PostFormValue("about_me")
		profilePicture := user.ProfilePicture

		file, header, err := r.FormFile("profile_picture")
		if err == nil {
			defer file.Close()

			// Create a unique filename
			fileName := fmt.Sprintf("%s_%s", userID, header.Filename)
			filePath := filepath.Join("ui/static/uploads", fileName)

			dst, err := os.Create(filePath)
			if err != nil {
				app.serverError(w, err)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				app.serverError(w, err)
				return
			}

			profilePicture = "/static/uploads/" + fileName
		}

		err = app.Models.UpdateUserProfile(userID, aboutMe, profilePicture)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/user/profile?id="+userID, http.StatusSeeOther)
		return
	}

	w.Header().Set("Allow", "GET, POST")
	app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}

// PasswordChange handles POST requests to change the user's password.
func (app *Application) PasswordChange(w http.ResponseWriter, r *http.Request) {
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

	currentPassword := r.PostForm.Get("current_password")
	newPassword := r.PostForm.Get("new_password")
	confirmPassword := r.PostForm.Get("confirm_password")

	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		app.renderProfileEditWithError(w, userID, "All password fields are required.")
		return
	}

	if newPassword != confirmPassword {
		app.renderProfileEditWithError(w, userID, "New passwords do not match.")
		return
	}

	if len(newPassword) < 6 {
		app.renderProfileEditWithError(w, userID, "New password must be at least 6 characters.")
		return
	}

	// Verify current password
	hashedPassword, err := app.Models.GetUserHashedPassword(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = auth.ComparePassword(hashedPassword, currentPassword)
	if err != nil {
		app.renderProfileEditWithError(w, userID, "Current password is incorrect.")
		return
	}

	// Hash the new password
	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.Models.UpdateUserPassword(userID, newHashedPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect back to profile edit with success
	user, _ := app.Models.GetUserByID(userID)
	app.render(w, http.StatusOK, "profile_edit.page.tmpl", &TemplateData{
		User:              user,
		IsAuthenticated:   true,
		AuthenticatedUser: user,
		SuccessMessage:    "Password changed successfully.",
	})
}

// renderProfileEditWithError is a helper to re-render the profile edit page with an error.
func (app *Application) renderProfileEditWithError(w http.ResponseWriter, userID, message string) {
	user, _ := app.Models.GetUserByID(userID)
	app.render(w, http.StatusOK, "profile_edit.page.tmpl", &TemplateData{
		User:              user,
		IsAuthenticated:   true,
		AuthenticatedUser: user,
		ErrorMessage:      message,
	})
}

// AccountDelete handles POST requests to permanently delete the user's account.
func (app *Application) AccountDelete(w http.ResponseWriter, r *http.Request) {
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

	// Delete the user account
	err := app.Models.DeleteUser(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ============================================================
// PRIVATE MESSAGING HANDLERS
// ============================================================

// Inbox displays all conversations for the authenticated user.
func (app *Application) Inbox(w http.ResponseWriter, r *http.Request) {
	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	convos, err := app.Models.GetConversations(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	unread, _ := app.Models.CountUnreadMessages(userID)

	app.render(w, http.StatusOK, "inbox.page.tmpl", &TemplateData{
		Conversations:     convos,
		IsAuthenticated:   true,
		AuthenticatedUser: app.getAuthenticatedUser(r),
		UnreadCount:       unread,
	})
}

// Conversation shows the message thread between the current user and another user.
func (app *Application) Conversation(w http.ResponseWriter, r *http.Request) {
	userID := app.getAuthenticatedUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	otherUserID := r.URL.Query().Get("with")
	if otherUserID == "" || otherUserID == userID {
		http.Redirect(w, r, "/messages", http.StatusSeeOther)
		return
	}

	// Verify the other user exists
	otherUser, err := app.Models.GetUserByID(otherUserID)
	if err != nil {
		app.clientError(w, http.StatusNotFound, "User not found")
		return
	}

	// Mark incoming messages from this user as read
	_ = app.Models.MarkMessagesRead(otherUserID, userID)

	msgs, err := app.Models.GetThread(userID, otherUserID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	unread, _ := app.Models.CountUnreadMessages(userID)

	app.render(w, http.StatusOK, "conversation.page.tmpl", &TemplateData{
		Messages:          msgs,
		OtherUser:         otherUser,
		IsAuthenticated:   true,
		AuthenticatedUser: app.getAuthenticatedUser(r),
		UnreadCount:       unread,
	})
}

// SendMessage handles POST requests to send a direct message.
func (app *Application) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	receiverID := r.PostForm.Get("receiver_id")
	content := strings.TrimSpace(r.PostForm.Get("content"))

	if receiverID == "" || content == "" || receiverID == userID {
		http.Redirect(w, r, "/messages", http.StatusSeeOther)
		return
	}

	msgID, err := auth.GenerateSessionID()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.Models.InsertMessage(msgID, userID, receiverID, content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/messages/conversation?with="+receiverID, http.StatusSeeOther)
}
