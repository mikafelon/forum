package handlersF

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"div-01/forum/Logic/queryF"

	"golang.org/x/crypto/bcrypt"
)

var cookie_session []string

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Clear the session cookie

		userid := r.FormValue("userid")

		var cookie_session_new []string

		for _, b := range cookie_session {
			if b != userid {
				cookie_session_new = append(cookie_session_new, b)
			}
		}

		cookie_session = cookie_session_new

		fmt.Println("Nouveau tableau", cookie_session)

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		// Redirect to login page
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/login.html")
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

		valid := true

		for _, b := range cookie_session {
			if b == user.ID {
				valid = false
				break
			}
		}

		if !valid {
			http.Error(w, "You cannot be connected to the same account simultaneously", http.StatusUnauthorized)
			return
		}

		cookie_session = append(cookie_session, user.ID)
		fmt.Println("Inscription de donnée", cookie_session, user.ID)
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
