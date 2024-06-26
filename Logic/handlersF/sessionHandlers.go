package handlersF

import (
	"database/sql"
	"net/http"
	"time"

	"div-01/forumM/Logic/queryF"
	"div-01/forumM/Logic/sessionF"
)

func ExtendSessionHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil || userID == "guest" {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	sessionID, err := sessionF.CreateSession(userID, db)
	if err != nil {
		http.Error(w, "Failed to extend session", http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(5 * time.Minute)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}
