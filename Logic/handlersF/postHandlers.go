package handlersF

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"forum/Logic/queryF"

	"github.com/google/uuid"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var input struct {
			Title      string `json:"title"`
			Content    string `json:"content"`
			CategoryID string `json:"category_id"`
		}
		input.Title = r.FormValue("title")
		input.Content = r.FormValue("content")
		input.CategoryID = r.FormValue("category_id")
		if input.Title == "" || input.Content == "" || input.CategoryID == "" {
			http.Error(w, "Title, Content, and Category are required", http.StatusBadRequest)
			return
		}
		postID := uuid.New().String()
		createdAt := time.Now().Format(time.RFC3339)
		err = queryF.InsertPost(postID, userID, input.Title, input.Content, input.CategoryID, createdAt, db)
		if err != nil {
			log.Printf("Error inserting post: %v\n", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	} else {
		// Display the create post page
		http.ServeFile(w, r, "templates/post.html")
	}
}
