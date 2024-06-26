package handlersF

import (
	"database/sql"
	"log"
	"net/http"

	"forum/Logic/queryF"
)

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil || userID == "guest" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		postID := r.FormValue("post_id")
		var existingValue int
		err := db.QueryRow("SELECT value FROM user_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingValue)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("Error querying user_likes: %v\n", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			err = insertLikeOrDislike(postID, userID, -1, db)
			if err == nil {
				createNotification(postID, "dislike", db)
			}
		} else {
			if existingValue == -1 {
				http.Redirect(w, r, "/forum", http.StatusSeeOther)
				return
			} else if existingValue == 1 {
				_, err = db.Exec("UPDATE user_likes SET value = -1 WHERE user_id = ? AND post_id = ?", userID, postID)
				if err == nil {
					createNotification(postID, "dislike", db)
				}
			}
		}
		if err != nil {
			log.Printf("Error inserting/updating dislike: %v\n", err)
			http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil || userID == "guest" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		commentID := r.FormValue("comment_id")
		var existingValue int
		err := db.QueryRow("SELECT value FROM user_likes WHERE user_id =? AND comment_id =?", userID, commentID).Scan(&existingValue)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("Error querying user_likes: %v\n", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			err = queryF.InsertLikeOrDislikeComment(commentID, userID, -1, db)
			if err == nil {
				createNotificationForComment(commentID, "dislike", db)
			}
		} else {
			if existingValue == -1 {
				http.Redirect(w, r, "/forum", http.StatusSeeOther)
				return
			} else if existingValue == 1 {
				_, err = db.Exec("UPDATE user_likes SET value = -1 WHERE user_id =? AND comment_id =?", userID, commentID)
				if err == nil {
					createNotificationForComment(commentID, "dislike", db)
				}
			}
		}
		if err != nil {
			log.Printf("Error inserting/updating dislike: %v\n", err)
			http.Error(w, "Failed to dislike comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}
