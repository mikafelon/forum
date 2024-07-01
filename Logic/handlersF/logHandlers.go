package handlersF

import (
	"database/sql"
	"net/http"
	"time"

	"forum/Logic/queryF"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/forum.html")
		return
	}
	if r.Method == http.MethodPost {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		input.Email = r.FormValue("email")
		input.Password = r.FormValue("password")
		user, err := queryF.GetUserByEmail(input.Email, db)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		// Set session cookie
		queryF.SetSessionCookie(w, user.ID, db)
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Clear the session cookie
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		// Redirect to login page
		http.Redirect(w, r, "/forum.html", http.StatusSeeOther)
	}
}
