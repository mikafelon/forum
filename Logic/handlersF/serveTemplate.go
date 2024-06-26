package handlersF

import (
	"database/sql"
	"html/template"
	"net/http"

	"forum/Logic/queryF"
)

func ServeTemplate(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	tmpl, err := template.ParseFiles("templates" + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.URL.Path == "/post.html" {
		categories, err := queryF.GetCategories(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, categories)
	} else {
		tmpl.Execute(w, nil)
	}
}
