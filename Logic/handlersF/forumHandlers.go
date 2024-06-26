package handlersF

import (
	"database/sql"
	"html/template"
	"net/http"

	"div-01/forumM/Logic/queryF"
	"div-01/forumM/Logic/typeF"
)

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userID, err := queryF.GetSessionUserID(r, db)
	if err != nil || userID == "guest" {
		// Rediriger vers la page de login si la session est expirée ou absente
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	// Récupérer les informations de l'utilisateur si ce n'est pas un invité
	user, err := queryF.GetUserByID(userID, db)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	searchQuery := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category_id")
	var posts []typeF.Post
	if categoryID != "" || searchQuery != "" {
		posts, err = queryF.GetFilteredPosts(userID, categoryID, searchQuery, db)
	} else {
		posts, err = queryF.GetAllPosts(userID, db)
	}
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	notifications, err := queryF.GetNotifications(userID, db)
	if err != nil {
		http.Error(w, "Failed to load notifications", http.StatusInternalServerError)
		return
	}

	// Marquer les notifications comme lues
	err = markNotificationsAsRead(userID, db)
	if err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	categories, err := queryF.GetCategories(db)
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	data := struct {
		User               typeF.User
		Posts              []typeF.Post
		Notifications      []typeF.Notification
		NotificationCount  int
		Categories         []typeF.Category
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

func markNotificationsAsRead(userID string, db *sql.DB) error {
	query := "UPDATE notifications SET is_read = 1 WHERE user_id = ? AND is_read = 0"
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
