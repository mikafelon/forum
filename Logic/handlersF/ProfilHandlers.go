package handlersF

import (
	"database/sql"
	"html/template"
	"net/http"

	"div-01/forum/Logic/queryF"
	"div-01/forum/Logic/typeF"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := queryF.GetUserByID(userID, db)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postUser, err := queryF.GetUserPosts(userID, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*
		commentUser, err := queryF.GetUserComments(userID, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

			likedComment, err := queryF.GetUserLikedComment(userID, db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}*/

	likedPost, err := queryF.GetUserLikedPosts(userID, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User  typeF.User
		Posts []typeF.Post
		// Comment      []typeF.Comment
		// LikedComment []typeF.Comment
		LikedPost []typeF.Post
	}{
		User:  user,
		Posts: postUser,
		// Comment: commentUser,
		// LikedComment: likedComment,
		LikedPost: likedPost,
	}

	tmpl.Execute(w, data)
}
