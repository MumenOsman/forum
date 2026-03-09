package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

/*
internal/models/models.go

Role: Database interaction and domain logic.

Responsibilities:
1. Define structures representing database entities (e.g., User, Post, Comment).
2. Handle the core SQL queries for all tables (Users, Posts, Comments, Categories, Likes).
3. Offer functions to insert, update, delete, and retrieve data cleanly, abstracting DB interactions from the handlers.
4. Manage database constraints conceptually before applying queries.
*/

// ErrNoRecord is returned when a query returns no rows.
var ErrNoRecord = errors.New("models: no matching record found")

// ErrInvalidCredentials is returned when email and password don't match.
var ErrInvalidCredentials = errors.New("models: invalid credentials")

// User structure for representing a user.
type User struct {
	ID        string
	Email     string
	Username  string
	Password  string // Hashed password
	CreatedAt string
}

// Post structure for representing a forum post.
type Post struct {
	ID           string
	UserID       string
	Username     string // For easier display
	Title        string
	Content      string
	Likes        int
	Dislikes     int
	CommentCount int // For easier display
	CreatedAt    string
	Categories   []*Category // Categories associated with this post
	Comments     []*Comment  // Comments for detail view
	UserLiked    bool        // State for current user
	UserDisliked bool        // State for current user
}

// Comment structure for representing a comment on a post.
type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt string
}

// Category represents a forum category.
type Category struct {
	ID   int
	Name string
}

// AppModel wraps the database connection pool so methods can access it.
type AppModel struct {
	DB *sql.DB
}

// New returns a new instance of an AppModel.
func New(db *sql.DB) *AppModel {
	return &AppModel{DB: db}
}

// Example Stub for inserting a user:
func (m *AppModel) InsertUser(id, email, username, hashedPassword string) error {
	stmt := `INSERT INTO users (id, email, username, password) VALUES (?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, id, email, username, hashedPassword)
	return err
}

// Example Stub for getting post:
func (m *AppModel) GetPost(id string) (*Post, error) {
	// Execute SELECT query here...
	return nil, ErrNoRecord
}

// Authenticate checks if an email and password match a database record.
// It returns the user ID if the credentials are valid.
func (m *AppModel) Authenticate(email, password string) (string, error) {
	var id string
	var hashedPassword string

	stmt := "SELECT id, password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	return id, nil
}

// InsertSession adds a new session into the database.
func (m *AppModel) InsertSession(sessionID, userID string, expiresAt time.Time) error {
	stmt := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := m.DB.Exec(stmt, sessionID, userID, expiresAt)
	return err
}

// DeleteSession removes a session from the database.
func (m *AppModel) DeleteSession(sessionID string) error {
	stmt := `DELETE FROM sessions WHERE id = ?`
	_, err := m.DB.Exec(stmt, sessionID)
	return err
}

// GetUserBySession retrieves the user associated with a valid session ID.
func (m *AppModel) GetUserBySession(sessionID string) (*User, error) {
	var user User
	stmt := `
		SELECT u.id, u.email, u.username, u.password, u.created_at
		FROM users u
		INNER JOIN sessions s ON u.id = s.user_id
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP`

	err := m.DB.QueryRow(stmt, sessionID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}
