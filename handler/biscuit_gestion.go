package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Create_biscuit_in_db(user_id int, token string, life_time time.Time) error {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Create_biscuit_in_db Open")
		return err
	}
	defer db.Close()
	biscuit := `INSERT INTO cookies_session (user_id, token, life_time) VALUES (?, ?, ?)`
	_, err = db.Exec(biscuit, user_id, token, life_time)
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Create_biscuit_in_db Exec(INSERT)", err)
		return err
	}
	return nil
}

func Update_cookie(user_id int, token string, life_time time.Time) error {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Update_cookie Open")
		return err
	}
	defer db.Close()
	biscuit := `UPDATE cookies_session SET token = ?, life_time = ?  WHERE user_id = ?`
	_, err = db.Exec(biscuit, token, life_time, user_id)
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Update_cookie Exec(UPDATE)")
		return err
	}
	return nil
}

func Check_biscuit(w http.ResponseWriter, r *http.Request) (bool, error) {
	// Ouvrez une connexion à la base de données SQLite
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Check_biscuit Open")
		return false, err
	}
	defer db.Close()
	// Récupérez le token de cookie depuis la requête HTTP
	biscuit, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Check_biscuit r.Cookie(\"session_token\")")
		Delete_biscuit(w, r)
		return false, err
	}
	sessionToken := biscuit.Value
	// Récupérez l'ID utilisateur et la date d'expiration depuis la base de données
	var user_id int
	var expirationDate time.Time
	row := db.QueryRow("SELECT user_id, life_time FROM cookies_session WHERE token = ?", sessionToken)
	err = row.Scan(&user_id, &expirationDate)

	if err != nil {
		fmt.Println("Error biscuit_gestion.go Check_biscuit Scan")
		return false, err
	}
	return true, nil
}

func Delete_biscuit(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Delete_biscuit Open", err)

		return
	}
	defer db.Close()
	biscuit, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		fmt.Println("Error biscuit_gestion.go Update_cookie r.Cookie(\"session_token\")", err)
		return
	}
	biscuit.MaxAge = -1
	http.SetCookie(w, biscuit)
	delete_life_time := `UPDATE cookies_session SET life_time = NULL WHERE user_id = ?`
	_, err = db.Exec(delete_life_time, User.ID)
	if err != nil {
		fmt.Println("Error biscuit_gestion.go Update_cookie Exec(UPDATE)", err)
		return
	}
	fmt.Println("redirect from biscuit_gestion.go Delete_cookie to home.go home_handler")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
