package main

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func main() {
	var err error
	// connect to the database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY, session_id TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to SQLite database.")
	// call the server function
	server()
	// if the users already logged in, display home.html, if not display auth.html

}
