package main

import (
	"literary-lions-forum/internal/database"
	"literary-lions-forum/internal/handlers"
	"literary-lions-forum/internal/models"
	"log"
	"net/http"
)

/*
cmd/web/main.go

Role: Entry point for the Literary Lions Forum application.

Responsibilities:
1. Parse runtime configuration (e.g., port number, DSN for database).
2. Establish a connection to the SQLite database via the internal/database package.
3. Initialize application state (e.g., models, session managers).
4. Define and register HTTP routes (combining URL paths with functions from internal/handlers).
5. Start the web server and listen for incoming HTTP requests.
*/
func main() {
	// 1. Initialize DB connection
	db, err := database.InitDB("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	// 2. Setup application dependencies
	app := models.New(db)
	_ = app // To avoid "declared and not used" error for now

	// 3. Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/post/view", handlers.PostView)

	// 4. Start server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
