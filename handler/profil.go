package handler

import (
	"fmt"
	"net/http"
	"text/template"
)

func Profil_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/profil" {
		biscuitOk, err := Check_biscuit(w, r)
		if err != nil {
			fmt.Println("error profil.go Profil_handler Check_biscuit")
		}
		if biscuitOk {
			posts, _ := Get_all_posts(r)
			tmpl, _ := template.New("").Funcs(template.FuncMap{
				"Get_username_by_id":        Get_username_by_id,
				"Get_categories_by_post_id": Get_categories_by_post_id,
				"Get_post_by_user_id":       Get_post_by_user_id,
				"Get_post_by_like":          Get_post_by_like,
			}).ParseFiles("templates/profil.html")
			_ = tmpl.ExecuteTemplate(w, "profil.html", struct {
				User  Users
				Posts []Post
			}{
				User:  User,
				Posts: posts,
			})
		}
	} else {
		fmt.Println("redirect from profil.go Profil_handler to 404 Error")
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}
}
