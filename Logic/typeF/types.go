package typeF

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
	ID           string `json:"id"`
	Content      string `json:"content"`
	UserID       string `json:"user_id"`
	PostID       string `json:"post_id"`
	CreatedAt    string `json:"created_at"`
	Username     string `json:"username"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	UserLiked    bool   `json:"user_liked"`
	UserDisliked bool   `json:"user_disliked"`
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
