package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func Register_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL.Path", r.URL.Path)
	if r.URL.Path == "/register" {
		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			fmt.Println("Error register.go Register_handler Open")
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		defer db.Close()
		switch r.Method {
		case "GET":
			tmpl, err := template.ParseFiles("templates/register.html")
			if err != nil {
				fmt.Println("error register.go register_handler ParseFiles")
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				fmt.Println("error register.go register_handler Execute")
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
		case "POST":
			// recup ton form
			username := r.FormValue("username")
			mail := r.FormValue("mail")
			password := r.FormValue("password")
			confirm_password := r.FormValue("confirm_password")
			// check password ==
			if confirm_password != password {
				fmt.Println("reload for same password")
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
			// insertion dans la DB
			err := Register_users(db, username, mail, password)
			if err != nil {
				fmt.Println("reload for same user")
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
			// incrémente la var User Users avec notre nouveau membre
			User, err := Get_user_by_username(db, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Création de son biscuit
			Biscuit(w, r)
			// envoi sur sa page de profil
			fmt.Println("redirect from register.go register_handler to profil.go profile_handler")
			http.Redirect(w, r, "/profil?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
			return

		default:
			fmt.Println("redirect from register.go register_handler to 400 error")
			http.Redirect(w, r, "/400", http.StatusSeeOther)
			return
		}
	} else {
		fmt.Println("redirect from register.go register_handler to 404 Error")
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}
}

func Register_users(db *sql.DB, username, mail, password string) error {
	// Prepare -> C'est un Open en prévision de multiple actions. Ne pas oublier de le defer .close()
	// Prepare est une fonctionnalité SQL, le contenu des "" sera concidé comme une Query
	user, err := db.Prepare("INSERT INTO users (username, mail, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error register.go Register_users Prepare(INSERT)")
		return err
	}
	defer user.Close()
	// tu executes l'INSERT
	_, err = user.Exec(username, mail, password)
	if err != nil {
		fmt.Println("Error register.go Register_users Exec(INSERT)")
	}
	return err
}
