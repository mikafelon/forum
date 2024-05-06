package logicF

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// registerHandler handles the registration process
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	print(" 01")
	if r.URL.Path != "/register/" { // if the url isn't filter it return nothing
		Error(w, http.StatusNotFound)
		return
	} else {

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		print(" 02")

		db, err := sql.Open("sqlite3", "./database.sqlite")
		if err != nil {
			log.Fatal(err)
		}
		print(" 03")

		defer db.Close()
		print(" 04")

		// Ensure the users table exists
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, email TEXT, username TEXT, password TEXT)")
		if err != nil {
			log.Fatal(err)
		}

		// Extract form data
		email := r.FormValue("e-mail")
		username := r.FormValue("username")
		password := r.FormValue("password")
		println(email)
		println(username)
		println(password)

		// Insert the data into the 'users' table
		_, err = db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Successfully registered user %s", username)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Successfully registered user %s", username)
}

func login(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Fatal(err)
	}

	if password == storedPassword {
		fmt.Fprintf(w, "Successfully logged in as %s", username)
	} else {
		fmt.Fprint(w, "Invalid username or password")
	}
}
