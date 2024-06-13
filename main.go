package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Post struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreatedAt    string    `json:"created_at"`
	Username     string    `json:"username"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	UserLiked    bool      `json:"user_liked"`
	UserDisliked bool      `json:"user_disliked"`
	Comments     []Comment `json:"comments"`
}
type Comment struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
}
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
}
type Notification struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	PostTitle string `json:"post_title"`
}
type Category struct {
	ID   string
	Name string
}

func getNotifications(userID string) ([]Notification, error) {
	query := `
        SELECT 
            notifications.id, notifications.user_id, notifications.post_id, notifications.type, notifications.created_at,
            users.username, posts.title
        FROM notifications
        JOIN users ON notifications.user_id = users.id
        JOIN posts ON notifications.post_id = posts.id
        WHERE notifications.user_id = ? AND notifications.is_read = 0
        ORDER BY notifications.created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []Notification
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.PostID, &notification.Type, &notification.CreatedAt, &notification.Username, &notification.PostTitle)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema.sql:", err)
	}
	_, err = db.Exec(string(file))
	if err != nil {
		log.Fatal("Failed to execute schema.sql:", err)
	}
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Clear the session cookie
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		// Redirect to login page
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
	}
}
func main() {
	initDatabase()
	defer db.Close()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/forum", forumHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/post", createPostHandler)
	http.HandleFunc("/post.html", serveTemplate)
	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/dislike", dislikeHandler)
	http.HandleFunc("/comments", commentsHandler)
	http.HandleFunc("/comment", commentHandler)
	http.HandleFunc("/extend-session", extendSessionHandler)
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Serve templates
	http.HandleFunc("/register.html", serveTemplate)
	http.HandleFunc("/login.html", serveTemplate)
	http.HandleFunc("/forum.html", serveTemplate)
	http.HandleFunc("/profile.html", serveTemplate)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
	if err != nil {
		userID = "guest"
	}
	// Récupérer les informations de l'utilisateur si ce n'est pas un invité
	var user User
	if userID != "guest" {
		user, err = getUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}
	} else {
		user = User{
			ID:       "guest",
			Username: "Guest",
		}
	}
	posts, err := getAllPosts(userID)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	categories, err := getCategories()
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}
	data := struct {
		User       User
		Posts      []Post
		Categories []Category
	}{
		User:       user,
		Posts:      posts,
		Categories: categories,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates" + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.URL.Path == "/post.html" {
		categories, err := getCategories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, categories)
	} else {
		tmpl.Execute(w, nil)
	}
}
func getCategories() ([]Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
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

		if userExists(input.Email) {
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
		err = insertUser(userID, input.Email, input.Username, string(hashedPassword), input.FirstName, input.LastName, createdAt)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		user := User{
			ID:        userID,
			Email:     input.Email,
			Username:  input.Username,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			CreatedAt: createdAt,
		}
		writeUserInfoToFile(user)
		// Set session cookie
		setSessionCookie(w, userID)
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/login.html")
		return
	}
	if r.Method == http.MethodPost {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		input.Email = r.FormValue("email")
		input.Password = r.FormValue("password")
		user, err := getUserByEmail(input.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		// Set session cookie
		setSessionCookie(w, user.ID)
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}
func forumHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
	if err != nil || userID == "guest" {
		// Rediriger vers la page de login si la session est expirée ou absente
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	// Récupérer les informations de l'utilisateur si ce n'est pas un invité
	user, err := getUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	searchQuery := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category_id")
	var posts []Post
	if categoryID != "" || searchQuery != "" {
		posts, err = getFilteredPosts(userID, categoryID, searchQuery)
	} else {
		posts, err = getAllPosts(userID)
	}
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	notifications, err := getNotifications(userID)
	if err != nil {
		http.Error(w, "Failed to load notifications", http.StatusInternalServerError)
		return
	}

	// Marquer les notifications comme lues
	err = markNotificationsAsRead(userID)
	if err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	categories, err := getCategories()
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	data := struct {
		User               User
		Posts              []Post
		Notifications      []Notification
		NotificationCount  int
		Categories         []Category
		SearchQuery        string
		SelectedCategoryID string
		SessionID          string
	}{
		User:               user,
		Posts:              posts,
		Notifications:      notifications,
		NotificationCount:  len(notifications),
		Categories:         categories,
		SearchQuery:        searchQuery,
		SelectedCategoryID: categoryID,
		SessionID:          r.Header.Get("Cookie"),
	}

	tmpl, err := template.ParseFiles("templates/forum.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func getFilteredPosts(userID, categoryID, searchQuery string) ([]Post, error) {
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
		query += " AND posts.category_id = ?"
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
	var posts []Post
	for rows.Next() {
		var post Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		comments, err := getComments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, nil
}
func markNotificationsAsRead(userID string) error {
	query := "UPDATE notifications SET is_read = 1 WHERE user_id = ? AND is_read = 0"
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := getUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, user)
}
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var input struct {
			Title      string `json:"title"`
			Content    string `json:"content"`
			CategoryID string `json:"category_id"`
		}
		input.Title = r.FormValue("title")
		input.Content = r.FormValue("content")
		input.CategoryID = r.FormValue("category_id")
		if input.Title == "" || input.Content == "" || input.CategoryID == "" {
			http.Error(w, "Title, Content, and Category are required", http.StatusBadRequest)
			return
		}
		postID := uuid.New().String()
		createdAt := time.Now().Format(time.RFC3339)
		err = insertPost(postID, userID, input.Title, input.Content, input.CategoryID, createdAt)
		if err != nil {
			log.Printf("Error inserting post: %v\n", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		// Redirect to forum
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	} else {
		// Display the create post page
		http.ServeFile(w, r, "templates/post.html")
	}
}
func insertPost(id, userID, title, content, categoryID, createdAt string) error {
	query := "INSERT INTO posts (id, user_id, title, content, category_id, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	log.Println("Executing query:", query)
	log.Println("With values:", id, userID, title, content, categoryID, createdAt)
	_, err := db.Exec(query, id, userID, title, content, categoryID, createdAt)
	if err != nil {
		log.Printf("Error executing insert post query: %v\n", err) // Log the SQL error
	}
	return err
}
func setSessionCookie(w http.ResponseWriter, userID string) {
	sessionID, err := createSession(userID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	expiration := time.Now().Add(5 * time.Minute) // 5 minutes session
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func extendSessionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
	if err != nil || userID == "guest" {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	sessionID, err := createSession(userID)
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

func writeUserInfoToFile(user User) {
	file, err := os.OpenFile("users.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
func userExists(email string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=? LIMIT 1)", email).Scan(&exists)
	if err != nil {
		log.Println("Failed to check if user exists:", err)
		return false
	}
	return exists
}
func insertUser(id, email, username, password, firstName, lastName, createdAt string) error {
	_, err := db.Exec("INSERT INTO users (id, email, username, password, first_name, last_name, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, email, username, password, firstName, lastName, createdAt)
	return err
}
func getUserByEmail(email string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, email, username, password, first_name, last_name, created_at FROM users WHERE email=?", email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt)
	return user, err
}
func getUserByID(id string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, email, username, first_name, last_name, created_at FROM users WHERE id=?", id).Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt)
	return user, err
}
func likeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
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
			err = insertLikeOrDislike(postID, userID, 1)
			if err == nil {
				createNotification(postID, userID, "like")
			}
		} else {
			if existingValue == 1 {
				http.Redirect(w, r, "/forum", http.StatusSeeOther)
				return
			} else if existingValue == -1 {
				_, err = db.Exec("UPDATE user_likes SET value = 1 WHERE user_id = ? AND post_id = ?", userID, postID)
				if err == nil {
					createNotification(postID, userID, "like")
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
func dislikeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
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
			err = insertLikeOrDislike(postID, userID, -1)
			if err == nil {
				createNotification(postID, userID, "dislike")
			}
		} else {
			if existingValue == -1 {
				http.Redirect(w, r, "/forum", http.StatusSeeOther)
				return
			} else if existingValue == 1 {
				_, err = db.Exec("UPDATE user_likes SET value = -1 WHERE user_id = ? AND post_id = ?", userID, postID)
				if err == nil {
					createNotification(postID, userID, "dislike")
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
func insertLikeOrDislike(postID, userID string, value int) error {
	likeID := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	query := "INSERT INTO user_likes (id, user_id, post_id, value, created_at) VALUES (?, ?, ?, ?, ?)"
	log.Println("Executing query:", query)
	log.Println("With values:", likeID, userID, postID, value, createdAt)
	_, err := db.Exec(query, likeID, userID, postID, value, createdAt)
	if err != nil {
		log.Printf("Error executing insert like/dislike query: %v\n", err)
	}
	return err
}
func createNotification(postID, userID, notificationType string) error {
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
func getAllPosts(userID string) ([]Post, error) {
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
	var posts []Post
	for rows.Next() {
		var post Post
		var likes, dislikes, userLikeValue int
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &likes, &dislikes, &userLikeValue)
		if err != nil {
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes
		post.UserLiked = userLikeValue == 1
		post.UserDisliked = userLikeValue == -1
		comments, err := getComments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, nil
}
func getComments(postID string) ([]Comment, error) {
	query := `
        SELECT 
            comments.id, comments.content, comments.user_id, comments.post_id, comments.created_at, users.username
        FROM comments
        JOIN users ON comments.user_id = users.id
        WHERE comments.post_id = ?
        ORDER BY comments.created_at DESC`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt, &comment.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
func commentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	comments, err := getComments(postID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}
	post, err := getPostByID(postID)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}
	data := struct {
		Post     Post
		Comments []Comment
	}{
		Post:     post,
		Comments: comments,
	}
	tmpl, err := template.ParseFiles("templates/comments.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
func getPostByID(postID string) (Post, error) {
	var post Post
	query := `
        SELECT 
            posts.id, posts.user_id, posts.title, posts.content, posts.created_at, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = ?`
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Username)
	return post, err
}
func commentHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSessionUserID(r)
	if err != nil || userID == "guest" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		userID, err := getSessionUserID(r)
		if err != nil || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		postID := r.FormValue("post_id")
		content := r.FormValue("comment")
		commentID := uuid.New().String()
		createdAt := time.Now().Format(time.RFC3339)
		query := "INSERT INTO comments (id, content, user_id, post_id, created_at) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, commentID, content, userID, postID, createdAt)
		if err != nil {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}
func getPostsByCategory(userID, categoryID string) ([]Post, error) {
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
	var posts []Post
	for rows.Next() {
		var post Post
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
		comments, err := getComments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, nil
}
func getSessionUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No session cookie found
		return "guest", nil
	}

	sessionID := cookie.Value
	var userID string
	var expiresAt time.Time

	err = db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return "", err
	}
	if time.Now().After(expiresAt) {
		// Session expired, destroy it
		destroySession(sessionID)
		return "guest", nil
	}

	return userID, nil
}

func createSession(userID string) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute) // 5 minutes session
	query := "INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)"
	_, err := db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func destroySession(sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getSessionUserID(r)
		if err != nil || userID == "guest" {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
