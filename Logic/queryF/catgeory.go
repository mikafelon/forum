package queryF

import (
	"database/sql"

	"forum/Logic/typeF"
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
