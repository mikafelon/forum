package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"forum/Logic/handlersF"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema.sql:", err)
	}
	_, err = db.Exec(string(file))
	if err != nil {
		log.Fatal("Failed to execute schema.sql:", err)
	}
}

func main() {
	InitDatabase()
	defer db.Close()
	http.HandleFunc("/", handlersF.IndexHandler)
	http.HandleFunc("/register", handlersF.RegisterHandler)
	http.HandleFunc("/login", handlersF.LoginHandler)
	http.HandleFunc("/forum", handlersF.ForumHandler)
	http.HandleFunc("/profile", handlersF.ProfileHandler)
	http.HandleFunc("/post", handlersF.CreatePostHandler)
	http.HandleFunc("/post.html", handlersF.ServeTemplate)

	// handle the like/dislike functionnalities
	http.HandleFunc("/like", handlersF.LikeHandler)
	http.HandleFunc("/dislike", handlersF.DislikeHandler)
	http.HandleFunc("/likeComment", handlersF.LikeCommentHandler)
	http.HandleFunc("/dislikeComment", handlersF.DislikeCommentHandler)

	// handle the comments functionnalities
	http.HandleFunc("/comments", handlersF.CommentsHandler)
	http.HandleFunc("/comment", handlersF.CommentHandler)
	http.HandleFunc("/extend-session", handlersF.ExtendSessionHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve templates
	http.HandleFunc("/register.html", handlersF.ServeTemplate)
	http.HandleFunc("/login.html", handlersF.ServeTemplate)
	http.HandleFunc("/forum.html", handlersF.ServeTemplate)
	http.HandleFunc("/profile.html", handlersF.ServeTemplate)
	log.Println("Server started at :8080")
	log.Print("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
