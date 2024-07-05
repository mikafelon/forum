package queryF

import (
	"database/sql"
	"log"

	"div-01/forum/Logic/typeF"
)

func GetComments(postID string, db *sql.DB) ([]typeF.Comment, error) {
	query := `
        SELECT
            comments.id, comments.content, comments.user_id, comments.post_id, comments.created_at, users.username,
            COALESCE(SUM(CASE WHEN likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
        FROM comments
        JOIN users ON comments.user_id = users.id
        LEFT JOIN user_likes AS likes ON comments.id = likes.comment_id
        WHERE comments.post_id =?
        GROUP BY comments.id, users.username
        ORDER BY comments.created_at DESC`
	rows, err := db.Query(query, postID)
	if err != nil {
		log.Printf("Error querying comments: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var comments []typeF.Comment
	for rows.Next() {
		var comment typeF.Comment
		var likes, dislikes int
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt, &comment.Username, &likes, &dislikes)
		if err != nil {
			log.Printf("Error scanning comment: %v\n", err)
			return nil, err
		}
		comment.Likes = likes
		comment.Dislikes = dislikes
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetUserComments(userID string, db *sql.DB) ([]typeF.Comment, error) { // Not Needed
	query := `
        SELECT id, post_id, content, created_at
        FROM comments
        WHERE user_id =?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []typeF.Comment
	for rows.Next() {
		var comment typeF.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetUserLikedComments(userID string, db *sql.DB) ([]typeF.Comment, error) { // Not needed
	query := `
        SELECT c.id, c.content, c.user_id, c.post_id, c.created_at AS comment_created_at, users.username, l.created_at AS like_created_at
        FROM likes l
        JOIN comments c ON l.comment_id = c.id
        JOIN users ON c.user_id = users.id
        WHERE l.user_id =?
        ORDER BY l.created_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likedComments []typeF.Comment
	for rows.Next() {
		var comment typeF.Comment
		var likeCreatedAt string // Temporary variable to hold the like creation time
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt, &comment.Username, &likeCreatedAt)
		if err != nil {
			return nil, err
		}
		// Optionally, you can store likeCreatedAt somewhere relevant if needed
		likedComments = append(likedComments, comment)
	}
	return likedComments, nil
}
