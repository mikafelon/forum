package queryF

import (
	"database/sql"
	"log"

	"div-01/forum/Logic/typeF"
)

func InsertPost(id, userID, title, content, categoryID, createdAt string, db *sql.DB) error {
	query := "INSERT INTO posts (id, user_id, title, content, category_id, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	log.Println("Executing query:", query)
	log.Println("With values:", id, userID, title, content, categoryID, createdAt)
	_, err := db.Exec(query, id, userID, title, content, categoryID, createdAt)
	if err != nil {
		log.Printf("Error executing insert post query: %v\n", err)
	}
	return err
}

func GetPostByID(postID string, db *sql.DB) (typeF.Post, error) {
	var post typeF.Post
	query := `
        SELECT
            posts.id, posts.user_id, posts.title, posts.content, posts.created_at, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = ?`
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username)
	return post, err
}

func GetPostsByCategory(userID, categoryID string, db *sql.DB) ([]typeF.Post, error) {
	query := `
        SELECT 
            posts.id, posts.user_id, posts.title, posts.content, posts.created_at, users.username,
            COALESCE(SUM(CASE WHEN user_likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN user_likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
            COALESCE(MAX(CASE WHEN user_likes.user_id = ? THEN user_likes.value END), 0) AS user_like_value
        FROM posts
        JOIN users ON posts.user_id = users.id
        LEFT JOIN user_likes ON posts.id = user_likes.post_id
        WHERE posts.category_id = ?
        GROUP BY posts.id, users.username
        ORDER BY posts.created_at DESC`
	rows, err := db.Query(query, userID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []typeF.Post
	for rows.Next() {
		var post typeF.Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		// Récupérer les commentaires pour chaque post
		comments, err := GetComments(post.ID, db)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllPosts(userID string, db *sql.DB) ([]typeF.Post, error) {
	query := `
        SELECT 
            posts.id, posts.user_id, posts.title, posts.content, posts.created_at, users.username,
            COALESCE(SUM(CASE WHEN user_likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN user_likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
            COALESCE(MAX(CASE WHEN user_likes.user_id = ? THEN user_likes.value END), 0) AS user_like_value
        FROM posts
        JOIN users ON posts.user_id = users.id
        LEFT JOIN user_likes ON posts.id = user_likes.post_id
        GROUP BY posts.id, users.username
        ORDER BY posts.created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []typeF.Post
	for rows.Next() {
		var post typeF.Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		comments, err := GetComments(post.ID, db)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, nil
}

func GetFilteredPosts(userID, categoryID, searchQuery string, db *sql.DB) ([]typeF.Post, error) {
	query := `
            SELECT 
                posts.id, posts.user_id, posts.title, posts.content, posts.created_at, users.username,
                COALESCE(SUM(CASE WHEN user_likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
                COALESCE(SUM(CASE WHEN user_likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
                COALESCE(MAX(CASE WHEN user_likes.user_id = ? THEN user_likes.value END), 0) AS user_like_value
            FROM posts
            JOIN users ON posts.user_id = users.id
            LEFT JOIN user_likes ON posts.id = user_likes.post_id
            WHERE 1=1`

	args := []interface{}{userID}

	if categoryID != "" {
		query += " AND EXISTS(SELECT 1 FROM post_categories WHERE post_categories.post_id = posts.id AND post_categories.category_id = ?)"
		args = append(args, categoryID)
	}

	if searchQuery != "" {
		query += " AND (posts.title LIKE ? OR posts.content LIKE ?)"
		args = append(args, "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	query += " GROUP BY posts.id, users.username ORDER BY posts.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []typeF.Post
	for rows.Next() {
		var post typeF.Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		comments, err := GetComments(post.ID, db)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}

	return posts, nil
}

func GetUserPosts(userID string, db *sql.DB) ([]typeF.Post, error) {
	query := `
        SELECT 
            posts.id, posts.user_id, posts.title, posts.content, posts.created_at, users.username,
            COALESCE(SUM(CASE WHEN user_likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN user_likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
            COALESCE(user_likes.value, 0) AS user_like_value -- Removed dynamic placeholder
        FROM posts
        JOIN users ON posts.user_id = users.id
        LEFT JOIN user_likes ON posts.id = user_likes.post_id AND user_likes.user_id = ?
        WHERE posts.user_id = ?
        GROUP BY posts.id, users.username
        ORDER BY posts.created_at DESC`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []typeF.Post
	for rows.Next() {
		var post typeF.Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetUserLikedPosts(userID string, db *sql.DB) ([]typeF.Post, error) {
	query := `
        SELECT posts.id, posts.title, posts.content, posts.created_at, users.username,
               COALESCE(SUM(CASE WHEN user_likes.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
               COALESCE(SUM(CASE WHEN user_likes.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
               COALESCE(user_likes.value, 0) AS user_like_value
        FROM posts
        JOIN users ON posts.user_id = users.id
        LEFT JOIN user_likes ON posts.id = user_likes.post_id AND user_likes.user_id = ?
        WHERE user_likes.user_id = ? AND user_likes.value = 1
        GROUP BY posts.id, users.username
        ORDER BY posts.created_at DESC`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likedPosts []typeF.Post
	for rows.Next() {
		var post typeF.Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		likedPosts = append(likedPosts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return likedPosts, nil
}
