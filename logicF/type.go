package logicF

type user struct {
	ID       int
	Username string
	password string
	email    string
}

type Post struct {
	IDPost      int
	UserID      int
	Username    string
	CategoryID  int
	Content     string
	PublishDate string
	Like        int
	Dislike     int
}

type erreur struct {
	Titre string
	Text  string
}
