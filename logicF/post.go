package logicF

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

// CustomResponseWriter wraps http.ResponseWriter to track if headers have been written.
type CustomResponseWriter struct {
	http.ResponseWriter
	headersWritten bool
}

// WriteHeader marks the headers as written.
func (w *CustomResponseWriter) WriteHeader(code int) {
	if !w.headersWritten {
		w.headersWritten = true
		w.ResponseWriter.WriteHeader(code)
	}
}

// Write ensures headers are written before the body.
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if !w.headersWritten {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func HomeHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Wrap the original ResponseWriter with our CustomResponseWriter
	rw := &CustomResponseWriter{ResponseWriter: w}

	// w: server to client / r: client to server
	if r.URL.Path != "/home.html" {
		Error(rw, http.StatusNotFound)
		return
	} else {
		// Fetch data from the database
		rows, err := db.Query("SELECT * FROM post")
		if err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}
		// username,err := db.Query("SELECT user_name FROM user WHERE post.post_user_id ")
		defer rows.Close()

		// Prepare a slice to hold the posts
		var posts []Post

		// Loop through the rows and scan the data into the posts slice
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.IDPost, &post.PublishDate, &post.Content, &post.CategoryID, &post.UserID)
			if err != nil {
				fmt.Println("Error scanning row:", err)
				Error(rw, http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// Check for errors from iterating over rows.
		if err = rows.Err(); err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}

		// Create the template
		tmpl, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}
		// Execute the template with the posts data
		err = tmpl.Execute(rw, posts)
		if err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}
	}
}
