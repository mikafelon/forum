package queryF

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

func InsertLikeOrDislikeComment(commentID, userID string, value int, db *sql.DB) error {
	likeID := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	query := "INSERT INTO user_likes (id, user_id, comment_id, value, created_at) VALUES (?,?,?,?,?)"
	log.Println("Executing query:", query)
	log.Println("With values:", likeID, userID, commentID, value, createdAt)
	_, err := db.Exec(query, likeID, userID, commentID, value, createdAt)
	if err != nil {
		log.Printf("Error executing insert like/dislike query: %v\n", err)
	}
	return err
}
