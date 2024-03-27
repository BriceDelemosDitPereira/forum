package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func Home_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL.Path", r.URL.Path)
	if r.URL.Path == "/" || r.URL.Path == "/favicon.ico" {
		if r.Method == "GET" {
			db, err := sql.Open("sqlite3", "./database/forum.db")
			if err != nil {
				fmt.Println("Error home.go Home_handler Open", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}

			defer db.Close()

			biscuitOk, err := Check_biscuit(w, r)
			if err != nil {
				fmt.Println("error home.go Home_handler Check_biscuit")
			}
			if biscuitOk {
				fmt.Println("redirect to home_connected")
				http.Redirect(w, r, "/home_connected?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
				return
			} else {
				posts, err := Get_all_posts(r)
				if err != nil {
					fmt.Println("Error home.go Home_handler Get_all_posts")
					http.Redirect(w, r, "/500", http.StatusSeeOther)
					return
				}
				tmpl, err := template.New("").Funcs(template.FuncMap{
					"Get_username_by_id":        Get_username_by_id,
					"Get_categories_by_post_id": Get_categories_by_post_id,
				}).ParseFiles("templates/home.html")
				if err != nil {
					fmt.Println("Error home.go Home_handler template.New().Funcs().ParseFiles")
					http.Redirect(w, r, "/500", http.StatusSeeOther)
					return
				}
				err = tmpl.ExecuteTemplate(w, "home.html", struct {
					User  Users
					Posts []Post
				}{
					User:  User,
					Posts: posts,
				})
				if err != nil {
					fmt.Println("Error home.go Home_handler ExecuteTemplate()")
					http.Redirect(w, r, "/500", http.StatusSeeOther)
					return
				}
			}

		} else {
			fmt.Println("400 error home.go home_handler")
			http.Redirect(w, r, "/400", http.StatusSeeOther)
			return
		}
	} else {
		fmt.Println("404 error home.go home_handler", r.URL.Path)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
}

func Home_connected_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path == "/home_connected" || r.URL.Path == "/create_post" || r.URL.Path == "/create_comment" || r.URL.Path == "/like_post" || r.URL.Path == "/register" {
			biscuitOk, err := Check_biscuit(w, r)
			if err != nil {
				fmt.Println("error home.go Home_connected_handler Check_biscuit")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			if biscuitOk {
				posts, err := Get_all_posts(r)
				if err != nil {
					fmt.Println("Error home.go Home_connected_handler Get_all_posts")
					http.Redirect(w, r, "/500", http.StatusSeeOther)
					return
				}
				tmpl, err := template.New("").Funcs(template.FuncMap{
					"Get_username_by_id":        Get_username_by_id,
					"Get_categories_by_post_id": Get_categories_by_post_id,
				}).ParseFiles("templates/home_connected.html")
				if err != nil {
					fmt.Println("Error home.go Home_connected_handler template.New().Funcs().ParseFiles")
					http.Redirect(w, r, "/500", http.StatusSeeOther)
					return
				}
				err = tmpl.ExecuteTemplate(w, "home_connected.html", struct {
					User  Users
					Posts []Post
				}{
					User:  User,
					Posts: posts,
				})
				if err != nil {
					fmt.Println("Error home.go Home_connected_handler ExecuteTemplate()")
					http.Redirect(w, r, "/500", http.StatusSeeOther)
					return

				}
			}
		} else {
			fmt.Println("404 error home.go home_connected")
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}
	} else {
		fmt.Println("400 error home.go home_connected")
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}
}
