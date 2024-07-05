package handlersF

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"div-01/forum/Logic/queryF"
	"div-01/forum/Logic/typeF"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodPost {
		var input struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			Password  string `json:"password"`
		}
		input.FirstName = r.FormValue("first_name")
		input.LastName = r.FormValue("last_name")
		input.Email = r.FormValue("email")
		input.Username = r.FormValue("username")
		input.Password = r.FormValue("password")

		data := struct {
			Error string
		}{}

		if userExists(input.Email, db) {
			data.Error = "Email already registered"
			tmpl, err := template.ParseFiles("templates/register.html")
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		userID := uuid.New().String()
		createdAt := time.Now().Format(time.RFC3339)
		err = queryF.InsertUser(userID, input.Email, input.Username, string(hashedPassword), input.FirstName, input.LastName, createdAt, db)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		user := typeF.User{
			ID:        userID,
			Email:     input.Email,
			Username:  input.Username,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			CreatedAt: createdAt,
		}
		writeUserInfoToFile(user)
		// Set session cookie
		queryF.SetSessionCookie(w, userID, db)
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

func writeUserInfoToFile(user typeF.User) {
	file, err := os.OpenFile("users.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println("Failed to open file:", err)
		return
	}
	defer file.Close()
	userInfo := fmt.Sprintf("ID: %s, FirstName: %s, LastName: %s, Email: %s, Username: %s, CreatedAt: %s\n",
		user.ID, user.FirstName, user.LastName, user.Email, user.Username, user.CreatedAt)
	if _, err := file.WriteString(userInfo); err != nil {
		log.Println("Failed to write to file:", err)
	}
}

func userExists(email string, db *sql.DB) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=? LIMIT 1)", email).Scan(&exists)
	if err != nil {
		log.Println("Failed to check if user exists:", err)
		return false
	}
	return exists
}
