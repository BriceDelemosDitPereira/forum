package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Init() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/forum.db") // Cette ligne ouvre une connexion à la base de données en utilisant le pilote SQLite3 et la chaîne de connexion "./database/forum.db". Si la base de données n'existe pas, elle sera créée.
	if err != nil {
		fmt.Println("Error database.go Init Open", err)

		return nil
	}

	if err = db.Ping(); err != nil { // Cette ligne vérifie si la connexion à la base de données est active. Si la connexion n'est pas active, cela signifie que la base de données n'a pas été créée correctement. Dans ce cas, la fonction retourne nil pour indiquer qu'il y a eu une erreur.
		fmt.Println("Error database.go Init Ping", err)
		return nil
	}

	create_multi_tables := `CREATE TABLE IF NOT EXISTS post (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		users_id INTEGER,
		title VARCHAR(50),
		content TEXT,
		image TEXT,
		publication_date DATETIME,
		post_like_count INTEGER DEFAULT 0,
		post_dislike_count INTEGER DEFAULT 0,
		FOREIGN KEY (users_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password VARCHAR(16),
		mail TEXT UNIQUE,
		creation_date DATETIME
	);
	
	CREATE TABLE IF NOT EXISTS post_comment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		users_id INTEGER,
		post_id INTEGER,
		content TEXT,
		creation_date DATETIME,
		nb_like INTEGER,
		comment_like_count INTEGER DEFAULT 0,
		comment_dislike_count INTEGER DEFAULT 0,
		FOREIGN KEY (users_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES post(id)
	);

	CREATE TABLE IF NOT EXISTS post_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		category TEXT,
		FOREIGN KEY (post_id) REFERENCES post(id)
	);
	
	CREATE TABLE IF NOT EXISTS cookies_session (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		token TEXT,
		life_time TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS POST_LIKE_DISLIKE (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		users_id INTEGER,
		posts_id INTEGER,
		comments_id INTEGER,
		like_status INTEGER,
		like_date DATETIME,
		FOREIGN KEY (users_id) REFERENCES users(id),
		FOREIGN KEY (posts_id) REFERENCES post(id),
		FOREIGN KEY (comments_id) REFERENCES post_comment(id)
	);`

	_, err = db.Exec(create_multi_tables) // Execute les instructions SQL
	if err != nil {
		fmt.Println("Error database.go Init Exec(CREATE)", err)
		// ERREUR500
		return nil
	}

	delete_life_time := `UPDATE cookies_session SET token = NULL`
	_, err = db.Exec(delete_life_time)
	if err != nil {
		fmt.Println("Error database.go Init Exec(UPDATE)", err)
		// ERREUR500
		return nil
	}

	return db
}
