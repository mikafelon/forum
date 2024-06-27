package sessionF

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateSession(userID string, db *sql.DB) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(500 * time.Minute) // 5 minutes session

	// Invalidate existing sessions for the user
	_, err := db.Exec("DELETE FROM sessions WHERE user_id =?", userID)
	if err != nil {
		return "", err
	}

	query := "INSERT INTO sessions (id, user_id, expires_at) VALUES (?,?,?)"
	_, err = db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func DestroySession(sessionID string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}
