package logicF

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

func HomeHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// w: server to client / r: client to server
	if r.URL.Path != "/index" {
		print("err1")
		ErrorHandler(w, r, http.StatusNotFound)
		return
	} else {
		// Fetch data from the database
		rows, err := db.Query("SELECT * FROM post")
		if err != nil {
			print("err2")
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Prepare a slice to hold the posts
		var posts []Post

		// Loop through the rows and scan the data into the posts slice
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.IDPost, &post.PublishDate, &post.Content, &post.CategoryID, &post.UserID)
			if err != nil {
				fmt.Println("Error scanning row:", err)
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// Check for errors from iterating over rows.
		if err = rows.Err(); err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		// Create the template
		tmpl, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		// Execute the template with the posts data
		err = tmpl.Execute(w, posts)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}
