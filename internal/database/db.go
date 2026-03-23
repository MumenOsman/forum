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
	db, err := sql.Open("sqlite3", dsn+"?_pragma=foreign_keys(1)")
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
		about_me TEXT DEFAULT '',
		profile_picture TEXT DEFAULT '',
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

	// Seed demo data for new databases
	SeedDemoData(db)

	// Apply migrations for existing users who might lack the new columns
	ApplyMigrations(db)

	return nil
}

// SeedDemoData inserts a demo user and sample post on first run.
// Uses INSERT OR IGNORE so existing data is never overwritten.
func SeedDemoData(db *sql.DB) {
	// Demo user — test@test.com / 123456
	// Password hash is bcrypt of "123456"
	const demoUserID = "demo-seed-user-001"
	const demoPasswordHash = "$2a$10$ejeGRPm7B7DabU/Hist3nOt5TO8A34T/pJxNBsEY56IxCJ4rQAHBG"
	_, err := db.Exec(`INSERT OR IGNORE INTO users (id, email, username, password, about_me)
		VALUES (?, ?, ?, ?, ?)`,
		demoUserID,
		"test@test.com",
		"Elara",
		demoPasswordHash,
		"A devoted reader and occasional scribbler of notes in book margins.",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to seed demo user: %v\n", err)
		return
	}

	// Demo post in Fiction category
	const demoPostID = "demo-seed-post-001"
	_, err = db.Exec(`INSERT OR IGNORE INTO posts (id, user_id, title, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		demoPostID,
		demoUserID,
		"The Name of the Wind — A Reader's First Encounter",
		`I picked up Patrick Rothfuss's "The Name of the Wind" on a whim after seeing it on a recommendation list, and I am genuinely struggling to put it down.

The story follows Kvothe — a legendary figure whose name alone inspires both awe and fear — as he sits in a quiet inn and begins recounting his life to a scribe. What unfolds is a tale of a prodigiously gifted child who grows up among traveling performers, loses everything, and claws his way into a prestigious academy of magic.

What strikes me most is the prose. Rothfuss writes with a rare kind of precision — every sentence feels deliberate, like it was weighed before being placed on the page. The magic system (called "Sympathy") is refreshingly grounded: it runs on physical laws, requires mental energy, and has real consequences when you get it wrong.

The framing narrative of Kvothe sitting in his inn, older and quieter than the legend suggests, also adds a layer of melancholy that I didn't expect from a fantasy novel. You know things went wrong somewhere. The tension of waiting to find out how is quietly unbearable.

If you enjoy fantasy that feels literary rather than purely adventurous, I'd strongly recommend giving this one a chapter or two. It earns every word.`,
		"2026-03-01 10:00:00",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to seed demo post: %v\n", err)
		return
	}

	// Link post to Fiction category
	var fictionID int
	_ = db.QueryRow(`SELECT id FROM categories WHERE name = 'Fiction'`).Scan(&fictionID)
	if fictionID > 0 {
		_, _ = db.Exec(`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			demoPostID, fictionID)
	}
}

// ApplyMigrations adds columns to existing tables if they don't exist.
func ApplyMigrations(db *sql.DB) {
	// Add about_me if it doesn't exist
	_, _ = db.Exec("ALTER TABLE users ADD COLUMN about_me TEXT DEFAULT ''")
	// Add profile_picture if it doesn't exist
	_, _ = db.Exec("ALTER TABLE users ADD COLUMN profile_picture TEXT DEFAULT ''")
	// Add likes/dislikes to comments if they don't exist
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN likes INTEGER DEFAULT 0")
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN dislikes INTEGER DEFAULT 0")
}

// initCategories is a helper to ensure at least some categories exist.
func initCategories(db *sql.DB) {
	categories := []string{"General", "Fiction", "Non-Fiction", "Sci-Fi", "Mystery", "Romance", "Historical Fiction", "Biography", "Fantasy"}
	for _, cat := range categories {
		query := "INSERT OR IGNORE INTO categories (name) VALUES (?)"
		_, err := db.Exec(query, cat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to insert category %s: %v\n", cat, err)
		}
	}
}
