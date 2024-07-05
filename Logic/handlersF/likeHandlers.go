package handlersF

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"div-01/forum/Logic/queryF"

	"github.com/google/uuid"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
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
			err = queryF.InsertLikeOrDislike(postID, userID, 1, db)
			if err == nil {
				createNotification(postID, "like", db)
			}
		} else {
			if existingValue == 1 {
				http.Redirect(w, r, "/forum", http.StatusSeeOther)
				return
			} else if existingValue == -1 {
				_, err = db.Exec("UPDATE user_likes SET value = 1 WHERE user_id = ? AND post_id = ?", userID, postID)
				if err == nil {
					createNotification(postID, "like", db)
				}
			}
		}
		if err != nil {
			log.Printf("Error inserting/updating like: %v\n", err)
			http.Error(w, "Failed to like post", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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
			err = queryF.InsertLikeOrDislikeComment(commentID, userID, 1, db)
			if err == nil {
				createNotificationForComment(commentID, "like", db)
			}
		} else {
			if existingValue == 1 {
				http.Redirect(w, r, "/forum", http.StatusSeeOther)
				return
			} else if existingValue == -1 {
				_, err = db.Exec("UPDATE user_likes SET value = 1 WHERE user_id =? AND comment_id =?", userID, commentID)
				if err == nil {
					createNotificationForComment(commentID, "like", db)
				}
			}
			if err != nil {
				log.Printf("Error inserting/updating like: %v\n", err)
				http.Error(w, "Failed to like comment", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/forum", http.StatusSeeOther)
		}
	}
}

func createNotification(postID, notificationType string, db *sql.DB) error {
	notificationID := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	// Get the post's user ID
	var postUserID string
	err := db.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&postUserID)
	if err != nil {
		return err
	}
	query := "INSERT INTO notifications (id, user_id, post_id, type, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, notificationID, postUserID, postID, notificationType, createdAt)
	if err != nil {
		log.Printf("Error inserting notification: %v\n", err)
	}
	return err
}

func createNotificationForComment(commentID, notificationType string, db *sql.DB) error {
	notificationID := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	// Assuming comments table has a user_id field to identify the author of the comment
	var commentUserID string
	err := db.QueryRow("SELECT user_id FROM comments WHERE id =?", commentID).Scan(&commentUserID)
	if err != nil {
		return err
	}
	query := "INSERT INTO notifications (id, user_id, comment_id, type, created_at) VALUES (?,?,?,?,?)"
	_, err = db.Exec(query, notificationID, commentUserID, commentID, notificationType, createdAt)
	if err != nil {
		log.Printf("Error inserting notification: %v\n", err)
	}
	return err
}
