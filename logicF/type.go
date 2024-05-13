package logicF

type HomeData struct {
	Posts      []Post
	Categories []Category
}

type user struct {
	id       int
	username string
	password string
	email    string
}

type Post struct {
	id         int
	userId     int
	username   string
	categoryId int
	content    string
	date       string
	like       int
	dislike    int
}

type Category struct {
	id   int
	name string
}

type erreur struct {
	Titre string
	Text  string
}
