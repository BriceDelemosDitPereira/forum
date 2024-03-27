package handler

import (
	"database/sql"
	"fmt"
)

func Insert_categories(categories []string, db *sql.DB, title string, content string) {
	var post_id int

	qwerry, err := db.Prepare("SELECT id FROM post WHERE title = ? AND content = ?")
	if err != nil {
		fmt.Println("Error categories.go Insert_categories Prepare(SELECT)", err)
	}

	defer qwerry.Close()

	qwerry.QueryRow(title, content).Scan(&post_id)

	for i := 0; i < len(categories); i++ {
		cat, err := db.Prepare("INSERT INTO post_categories (post_id, category) VALUES (?, ?)")
		if err != nil {
			fmt.Println("Error categories.go Insert_categories Prepare(INSERT)", err)
		}
		defer cat.Close()
		// ERREUR500
		_, err = cat.Exec(post_id, categories[i])
		if err != nil {
			fmt.Println("Error biscuit_gestion.go Insert_categories Exec(INSERT)", err)
			// ERREUR500
		}
	}
}

func Get_categories_by_post_id(post_id int) (categories []string) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error categories.go Get_categories_by_post_id Open", err)
		// ERREUR500
	}
	defer db.Close()
	query, err := db.Prepare("SELECT category FROM post_categories WHERE post_id = ?")
	if err != nil {
		fmt.Println("error exe recup post_id from Get_categories_by_post_id", err)
		// ERREUR500
	}
	defer query.Close()
	rows, err := query.Query(post_id) // On fait un Query et pas un QueryRow car on a plusieurs lignes à récupérer
	if err != nil {
		fmt.Println("error exe scan from Get_categories_by_post_id", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			fmt.Println("error exe scan from Get_categories_by_post_id", err)
			return
		}
		categories = append(categories, category) // categories déjà initialisé dans la ligne de la fonction (à la fin)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("error exe rows from Get_categories_by_post_id", err)
		return

	}
	// categories contient maintenant toutes les catégories associées au post_id
	return categories
}
