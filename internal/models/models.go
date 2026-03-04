package models

/*
internal/models/models.go

Role: Database interaction and domain logic.

Responsibilities:
1. Define structures representing database entities (e.g., User, Post, Comment).
2. Handle the core SQL queries for all tables (Users, Posts, Comments, Categories, Likes).
3. Offer functions to insert, update, delete, and retrieve data cleanly, abstracting DB interactions from the handlers.
4. Manage database constraints conceptually before applying queries.
*/

// User structure for representing a user.
type User struct {
	ID       string
	Email    string
	Username string
	Password string // Hashed password
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

// ... other models (Comment, Category) go here later.
