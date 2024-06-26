package queryF

import (
	"database/sql"

	"forum/Logic/typeF"
)

func GetNotifications(userID string, db *sql.DB) ([]typeF.Notification, error) {
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
	var notifications []typeF.Notification
	for rows.Next() {
		var notification typeF.Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.PostID, &notification.Type, &notification.CreatedAt, &notification.Username, &notification.PostTitle)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
