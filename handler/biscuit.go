package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func Biscuit(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error biscuit.go Biscuit Open", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	defer db.Close()
	token, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Error biscuit.go Biscuit uuid.NewV4()", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	now := time.Now()
	session_token := token.String()
	life_time := now.Add(2 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session_token,
		Expires: life_time,
	})
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM cookies_session WHERE user_id = ?", User.ID)
	err = row.Scan(&count)
	if err != nil {
		fmt.Println("Error biscuit.go Biscuit Scan", err)
	}
	if count > 0 {
		Update_cookie(User.ID, session_token, life_time)
	} else {
		Create_biscuit_in_db(User.ID, session_token, life_time)
	}
}
