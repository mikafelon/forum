package handlersF

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"div-01/forum/Logic/queryF"

	"github.com/google/uuid"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Println("Error parsing multipart form:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var input struct {
			Title      string `json:"title"`
			Content    string `json:"content"`
			CategoryID string `json:"category_id[]"`
		}
		input.Title = r.FormValue("title")
		input.Content = r.FormValue("content")
		categoryIDs := r.MultipartForm.Value["category_id[]"]

		// Validate that the content is not just whitespace
		if strings.TrimSpace(input.Content) == "" {
			http.Error(w, "Content must not be empty or consist only of whitespace", http.StatusBadRequest)
			return
		}
		if input.Title == "" || input.Content == "" || categoryIDs[0] == "" {
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
		for _, categoryP := range categoryIDs {
			// print(categoryP)
			postCatID := uuid.New().String()
			err := queryF.InsertCategoriesPost(postCatID, postID, categoryP, db)
			if err != nil {
				log.Printf("Error attributing categories for a post: %v\n", err)
				http.Error(w, "Failed to create post", http.StatusInternalServerError)
				return
			}
		}
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	} else {
		// Display the create post page
		http.ServeFile(w, r, "templates/post.html")
	}
}
