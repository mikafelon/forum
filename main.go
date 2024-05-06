package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/logicF"
)

func main() {
	var err error
	// connect to the database
	db, err := sql.Open("sqlite3", "./database.sqlite")
	// fmt.Println("db in main:", db)
	// println()
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
	fmt.Printf("Starting server at port :8080\n Serving on http://localhost:8080/home.html\n")

	// fmt.Println("db before HomeHandler:", db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/home.html", func(w http.ResponseWriter, r *http.Request) {
		logicF.HomeHandler(db, w, r)
	})
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))

	http.HandleFunc("/register/", logicF.RegisterHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)

	// if the users already logged in, display home.html, if not display auth.html
}
