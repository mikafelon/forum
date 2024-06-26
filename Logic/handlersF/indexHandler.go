package handlersF

import (
	"database/sql"
	"html/template"
	"net/http"

	"div-01/forumM/Logic/queryF"
	"div-01/forumM/Logic/typeF"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil {
		userID = "guest"
	}

	// Récupérer les informations de l'utilisateur si ce n'est pas un invité
	var user typeF.User
	if userID != "guest" {
		user, err = queryF.GetUserByID(userID, db)
		if err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}
	} else {
		user = typeF.User{
			ID:       "guest",
			Username: "Guest",
		}
	}
	posts, err := queryF.GetAllPosts(userID, db)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	categories, err := queryF.GetCategories(db)
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}
	data := struct {
		User       typeF.User
		Posts      []typeF.Post
		Categories []typeF.Category
	}{
		User:       user,
		Posts:      posts,
		Categories: categories,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
