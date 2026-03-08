package models

import (
	"database/sql"
	"errors"

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
	ID        string
	UserID    string
	Title     string
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt string
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
