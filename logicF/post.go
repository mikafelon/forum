package logicF

import (
	"database/sql"
	"fmt"
	"log"
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
		// Fetch posts from the database
		rows, err := db.Query("SELECT * FROM post")
		if err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Prepare a slice to hold the posts
		var posts []Post

		// Loop through the rows and scan the data into the posts slice
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.id, &post.date, &post.content, &post.categoryId, &post.userId)
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

		// Fetch categories from the database
		rows, err = db.Query("SELECT * FROM category")
		if err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Prepare a slice to hold the categories
		var categories []Category

		// Loop through the rows and scan the data into the categories slice
		for rows.Next() {
			var category Category
			err := rows.Scan(&category.id, &category.name)
			if err != nil {
				fmt.Println("Error scanning row:", err)
				Error(rw, http.StatusInternalServerError)
				return
			}
			categories = append(categories, category)
		}

		// Check for errors from iterating over rows.
		if err = rows.Err(); err != nil {
			Error(rw, http.StatusInternalServerError)
			return
		}

		// Create the HomeData struct
		homeData := HomeData{
			Posts:      posts,
			Categories: categories,
		}

		log.Println("Parsing home.html template")
		tmpl, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			log.Println("Error parsing home.html:", err)
			Error(rw, http.StatusInternalServerError)
			return
		}

		log.Println("Executing template with homeData")
		err = tmpl.Execute(rw, homeData)
		if err != nil {
			log.Println("Error executing template:", err)
			Error(rw, http.StatusInternalServerError)
			return
		}
	}
}
