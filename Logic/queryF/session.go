package queryF

import (
	"database/sql"
	"net/http"
	"time"

	"div-01/forum/Logic/sessionF"
)

func GetSessionUserID(r *http.Request, db *sql.DB) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No session cookie found
		return "guest", nil
	}

	sessionID := cookie.Value
	var userID string
	var expiresAt time.Time

	err = db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE id =?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return "", err
	}
	if time.Now().After(expiresAt) {
		// Session expired, destroy it
		sessionF.DestroySession(sessionID, db)
		return "guest", nil
	}

	return userID, nil
}

func SetSessionCookie(w http.ResponseWriter, userID string, db *sql.DB) {
	sessionID, err := sessionF.CreateSession(userID, db)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	expiration := time.Now().Add(15 * time.Minute) // 15 minutes session
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
