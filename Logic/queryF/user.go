package queryF

import (
	"database/sql"

	"forum/Logic/typeF"
)

func InsertUser(id, email, username, password, firstName, lastName, createdAt string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (id, email, username, password, first_name, last_name, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, email, username, password, firstName, lastName, createdAt)
	return err
}

func GetUserByEmail(email string, db *sql.DB) (typeF.User, error) {
	var user typeF.User
	err := db.QueryRow("SELECT id, email, username, password, first_name, last_name, created_at FROM users WHERE email=?", email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt)
	return user, err
}

func GetUserByID(id string, db *sql.DB) (typeF.User, error) {
	var user typeF.User
	err := db.QueryRow("SELECT id, email, username, first_name, last_name, created_at FROM users WHERE id=?", id).Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt)
	return user, err
}
