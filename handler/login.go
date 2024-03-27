package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type Users struct {
	ID       int
	Username string
	Mail     string
	Password string
}

func Login_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			fmt.Println("Error login.go Login_handler Open", err)
			return
		}
		defer db.Close()
		switch r.Method {
		case "GET":
			biscuitOk, err := Check_biscuit(w, r)
			if err != nil {
				fmt.Println("error login.go Login_handler Check_biscuit")
			}
			if biscuitOk && User.ID != 0 {
				fmt.Println("redirect from login.go Login_handler to profil.go profil_handler")

				http.Redirect(w, r, "/profil?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
				return
			} else {
				tmpl, err := template.ParseFiles("templates/login.html")
				if err != nil {
					return
				}
				err = tmpl.Execute(w, nil)
				if err != nil {
					return
				}
			}
		case "POST":
			username := r.FormValue("username")
			password := r.FormValue("password")
			user, err := Get_user_by_username(db, username)
			if err != nil || user == nil {
				fmt.Println("reload for nil user|err")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else if user.Password != password {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			Biscuit(w, r)
			fmt.Println("redirect from login.go Login_handler to home.go home_connected")
			http.Redirect(w, r, "/home_connected?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
			return
		}
	} else {
		fmt.Println("redirect from login.go Login_handler to 404 Error")
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
}

var User Users

// check si le username renseigné par le form, est présent dans les username de la db
func Get_user_by_username(db *sql.DB, username string) (*Users, error) {
	row := db.QueryRow("SELECT id, username, password, mail FROM users WHERE username = ? OR mail = ?", username, username)
	err := row.Scan(&User.ID, &User.Username, &User.Password, &User.Mail)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &User, nil
}

// Fonction qui ...
func Get_username_by_id(user_id int) (username string, err error) {
	// log.Println("Get_username_by_id called with user_id:", user_id)
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error login.go Get_username_by_id Open", err)
		return "", err
	}
	defer db.Close()
	var name string
	query, err := db.Prepare("SELECT username FROM users WHERE id = ?")
	if err != nil {
		fmt.Println("error for recup username in Get_username_by_id", err)
	}
	defer query.Close()
	row := query.QueryRow(user_id)
	err = row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, err
}
