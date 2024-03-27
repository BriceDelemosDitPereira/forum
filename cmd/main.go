package main

import (
	"fmt"
	"net/http"

	database "forum/database"
	handler "forum/handler"
)

func main() {
	// Open the database
	database.Init()
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))

	http.Handle("/img_upload/", http.StripPrefix("/img_upload/", http.FileServer(http.Dir("img_upload"))))

	http.HandleFunc("/", handler.Home_handler)
	http.HandleFunc("/home_connected", handler.Home_connected_handler)
	http.HandleFunc("/login", handler.Login_handler)
	http.HandleFunc("/register", handler.Register_handler)
	http.HandleFunc("/profil", handler.Profil_handler)
	http.HandleFunc("/delete", handler.Delete_biscuit)
	http.HandleFunc("/create_post", handler.Create_post)
	http.HandleFunc("/create_comment", handler.Create_comment)
	http.HandleFunc("/like_post", handler.LikePostsHandler)
	http.HandleFunc("/404", handler.NotFoundHandler)
	http.HandleFunc("/400", handler.BadRequest)
	http.HandleFunc("/500", handler.StatusInternalServerError)
	http.HandleFunc("/favicon.ico", handler.HandleFavicon)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
