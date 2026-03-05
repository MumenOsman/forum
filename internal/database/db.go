package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

/*
internal/database/db.go

Role: Initialize and configure the database connection.

Responsibilities:
1. Establish the connection to the SQLite3 database file.
2. Read the connection details from the environment (or defaults).
3. Ensure initial connection is verified (e.g., using db.Ping()).
4. Provide access to the database object for the application models to utilize.
5. Create and run schema definition queries if tables don't exist upon startup.
*/

// InitDB attempts to connect to a given dsn (data source name).
// Example: InitDB("forum.db")
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Verify the connection is active
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

// InitSchema creates the necessary tables if they don't already exist.
func InitSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS post_categories (
		post_id TEXT NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS likes_dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		target_id TEXT NOT NULL,
		target_type TEXT CHECK(target_type IN ('post', 'comment')) NOT NULL,
		vote_type INTEGER CHECK(vote_type IN (1, -1)) NOT NULL,
		UNIQUE(user_id, target_id, target_type),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("error executing schema: %w", err)
	}

	// Populate default categories if empty
	initCategories(db)

	return nil
}

// initCategories is a helper to ensure at least some categories exist.
func initCategories(db *sql.DB) {
	categories := []string{"General", "Fiction", "Non-Fiction", "Sci-Fi", "Mystery", "Romance"}
	for _, cat := range categories {
		query := "INSERT OR IGNORE INTO categories (name) VALUES (?)"
		_, err := db.Exec(query, cat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to insert category %s: %v\n", cat, err)
		}
	}
}
