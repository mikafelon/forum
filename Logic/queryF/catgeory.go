package queryF

import (
	"database/sql"
	"log"

	"div-01/forum/Logic/typeF"
)

func GetCategories(db *sql.DB) ([]typeF.Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []typeF.Category
	for rows.Next() {
		var category typeF.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func InsertCategoriesPost(id, postID, categoryID string, db *sql.DB) error {
	query := "INSERT INTO post_categories (id, post_id, category_id) VALUES (?, ?, ?)"
	log.Println("Executing query:", query)
	log.Println("With values:", id, postID, categoryID)
	_, err := db.Exec(query, id, postID, categoryID)
	if err != nil {
		log.Printf("Error executing insert post query: %v\n", err)
	}
	return err
}
