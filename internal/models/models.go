package models

import (
	"database/sql"
	"errors"
	"strconv"
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
	ID             string
	Email          string
	Username       string
	Password       string // Hashed password
	AboutMe        string
	ProfilePicture string
	CreatedAt      time.Time
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
	CreatedAt    time.Time
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
	Username  string // For easier display
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt time.Time
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

// InsertPost inserts a new post into the database.
func (m *AppModel) InsertPost(postID, userID, title, content string, categoryIDs []int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO posts (id, user_id, title, content) VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(stmt, postID, userID, title, content)
	if err != nil {
		return err
	}

	for _, catID := range categoryIDs {
		catStmt := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
		_, err = tx.Exec(catStmt, postID, catID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetFilteredPosts retrieves posts based on optional filters.
func (m *AppModel) GetFilteredPosts(categoryID, authoredBy, likedBy, searchQuery string) ([]*Post, error) {
	stmt := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.likes, p.dislikes, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id`

	var args []interface{}
	var joins []string
	var where []string

	if categoryID != "" {
		joins = append(joins, "INNER JOIN post_categories pc ON p.id = pc.post_id")
		where = append(where, "pc.category_id = ?")
		args = append(args, categoryID)
	}

	if likedBy != "" {
		joins = append(joins, "INNER JOIN likes_dislikes ld ON p.id = ld.target_id")
		where = append(where, "ld.user_id = ? AND ld.target_type = 'post' AND ld.vote_type = 1")
		args = append(args, likedBy)
	}

	for _, j := range joins {
		stmt += " " + j
	}

	if authoredBy != "" {
		where = append(where, "p.user_id = ?")
		args = append(args, authoredBy)
	}

	if searchQuery != "" {
		where = append(where, "(p.title LIKE '%' || ? || '%' COLLATE NOCASE OR p.content LIKE '%' || ? || '%' COLLATE NOCASE)")
		args = append(args, searchQuery, searchQuery)
	}

	if len(where) > 0 {
		stmt += " WHERE "
		for i, w := range where {
			if i > 0 {
				stmt += " AND "
			}
			stmt += w
		}
	}

	stmt += ` ORDER BY p.created_at DESC`

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.CommentCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetAllCategories retrieves all categories.
func (m *AppModel) GetAllCategories() ([]*Category, error) {
	stmt := `SELECT id, name FROM categories ORDER BY name ASC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetPostByID retrieves a specific post and its associated comments.
func (m *AppModel) GetPostByID(postID string) (*Post, error) {
	stmt := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.likes, p.dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comment_count,
		       p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`

	var p Post
	err := m.DB.QueryRow(stmt, postID).Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.CommentCount, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	// Fetch comments
	cStmt := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.likes, c.dislikes, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`

	rows, err := m.DB.Query(cStmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		p.Comments = append(p.Comments, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &p, nil
}

// InsertComment adds a comment to a specific post and returns the new comment's ID.
func (m *AppModel) InsertComment(postID, userID, content string) (string, error) {
	commentID := strconv.FormatInt(time.Now().UnixNano(), 36)

	stmt := `INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, commentID, postID, userID, content)
	if err != nil {
		return "", err
	}

	return commentID, nil
}

// GetCommentByID retrieves a single comment by its ID.
func (m *AppModel) GetCommentByID(id string) (*Comment, error) {
	stmt := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.likes, c.dislikes, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`

	var c Comment
	err := m.DB.QueryRow(stmt, id).Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &c, nil
}

// GetVoteCounts returns the total likes and dislikes for a given target.
func (m *AppModel) GetVoteCounts(targetID, targetType string) (struct{ Likes, Dislikes int }, error) {
	var counts struct{ Likes, Dislikes int }
	
	// Since we cache likes/dislikes in the posts and comments tables, we can just read them from there.
	var stmt string
	if targetType == "post" {
		stmt = `SELECT likes, dislikes FROM posts WHERE id = ?`
	} else {
		stmt = `SELECT likes, dislikes FROM comments WHERE id = ?`
	}

	err := m.DB.QueryRow(stmt, targetID).Scan(&counts.Likes, &counts.Dislikes)
	return counts, err
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
		SELECT u.id, u.email, u.username, u.password, u.about_me, u.profile_picture, u.created_at
		FROM users u
		INNER JOIN sessions s ON u.id = s.user_id
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP`

	err := m.DB.QueryRow(stmt, sessionID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.AboutMe,
		&user.ProfilePicture,
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

// GetUserByID retrieves a specific user's details.
func (m *AppModel) GetUserByID(userID string) (*User, error) {
	stmt := `SELECT id, email, username, about_me, profile_picture, created_at FROM users WHERE id = ?`
	var user User
	err := m.DB.QueryRow(stmt, userID).Scan(&user.ID, &user.Email, &user.Username, &user.AboutMe, &user.ProfilePicture, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile updates a user's about_me and profile_picture.
func (m *AppModel) UpdateUserProfile(userID, aboutMe, profilePicture string) error {
	stmt := `UPDATE users SET about_me = ?, profile_picture = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, aboutMe, profilePicture, userID)
	return err
}

// UpdateUserNameUsername updates a user's username.
func (m *AppModel) UpdateUserNameUsername(userID, username string) error {
	stmt := `UPDATE users SET username = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, username, userID)
	return err
}

// UpdateUserPassword updates a user's hashed password.
func (m *AppModel) UpdateUserPassword(userID, newHashedPassword string) error {
	stmt := `UPDATE users SET password = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, newHashedPassword, userID)
	return err
}

// GetUserHashedPassword retrieves the hashed password for a user.
func (m *AppModel) GetUserHashedPassword(userID string) (string, error) {
	var hashedPassword string
	stmt := `SELECT password FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, userID).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

// DeleteUser permanently deletes a user and all their associated data.
func (m *AppModel) DeleteUser(userID string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete sessions
	_, err = tx.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	// Delete votes by this user
	_, err = tx.Exec(`DELETE FROM likes_dislikes WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	// Delete comments by this user
	_, err = tx.Exec(`DELETE FROM comments WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	// Delete post categories for this user's posts
	_, err = tx.Exec(`DELETE FROM post_categories WHERE post_id IN (SELECT id FROM posts WHERE user_id = ?)`, userID)
	if err != nil {
		return err
	}

	// Delete posts by this user
	_, err = tx.Exec(`DELETE FROM posts WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = tx.Exec(`DELETE FROM users WHERE id = ?`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// InsertOrUpdateVote registers a user's vote on a specific post or comment.
func (m *AppModel) InsertOrUpdateVote(userID, targetID, targetType string, voteType int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Check if the user has already voted on this specific target
	var existingVote int
	stmtCheck := `SELECT vote_type FROM likes_dislikes WHERE user_id = ? AND target_id = ? AND target_type = ?`
	err = tx.QueryRow(stmtCheck, userID, targetID, targetType).Scan(&existingVote)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User hasn't voted yet = Insert the new vote
			stmtInsert := `INSERT INTO likes_dislikes (user_id, target_id, target_type, vote_type) VALUES (?, ?, ?, ?)`
			_, err = tx.Exec(stmtInsert, userID, targetID, targetType, voteType)
			if err != nil {
				return err
			}

			// Increment the target
			incrementTarget(tx, targetID, targetType, voteType, 1)
		} else {
			return err
		}
	} else {
		if existingVote == voteType {
			// User clicked the exact same button they already pressed = Unlike/Undislike
			stmtDelete := `DELETE FROM likes_dislikes WHERE user_id = ? AND target_id = ? AND target_type = ?`
			_, err = tx.Exec(stmtDelete, userID, targetID, targetType)
			if err != nil {
				return err
			}

			// Decrement the target
			incrementTarget(tx, targetID, targetType, existingVote, -1)
		} else {
			// User clicked the opposite button = Swap their vote
			stmtUpdate := `UPDATE likes_dislikes SET vote_type = ? WHERE user_id = ? AND target_id = ? AND target_type = ?`
			_, err = tx.Exec(stmtUpdate, voteType, userID, targetID, targetType)
			if err != nil {
				return err
			}

			// Decrement the old vote, Increment the new vote
			incrementTarget(tx, targetID, targetType, existingVote, -1)
			incrementTarget(tx, targetID, targetType, voteType, 1)
		}
	}

	return tx.Commit()
}

// incrementTarget is a helper to update the total counts on the raw posts or comments table
func incrementTarget(tx *sql.Tx, targetID, targetType string, voteType, increment int) {
	table := "posts"
	if targetType == "comment" {
		table = "comments"
	}

	column := "likes"
	if voteType == -1 {
		column = "dislikes"
	}

	stmt := `UPDATE ` + table + ` SET ` + column + ` = ` + column + ` + ? WHERE id = ?`
	tx.Exec(stmt, increment, targetID)
}
