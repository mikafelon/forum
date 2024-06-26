package handlersF

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"div-01/forumM/Logic/queryF"
	"div-01/forumM/Logic/typeF"

	"github.com/google/uuid"
)

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	comments, err := queryF.GetComments(postID, db)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}
	post, err := queryF.GetPostByID(postID, db)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}
	data := struct {
		Post     typeF.Post
		Comments []typeF.Comment
	}{
		Post:     post,
		Comments: comments,
	}
	tmpl, err := template.ParseFiles("templates/comments.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil || userID == "guest" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		userID, err := queryF.GetSessionUserID(r, db)
		if err != nil || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		postID := r.FormValue("post_id")
		content := r.FormValue("comment")
		commentID := uuid.New().String()
		createdAt := time.Now().Format(time.RFC3339)
		query := "INSERT INTO comments (id, content, user_id, post_id, created_at) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, commentID, content, userID, postID, createdAt)
		if err != nil {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}
