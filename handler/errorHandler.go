package handler

import (
	"fmt"
	"net/http"
	"text/template"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		fmt.Println("Error errorHandler.go NotFoundHandler ParseFile")
		http.Redirect(w, r, "/500", http.StatusSeeOther)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error errorHandler.go NotFoundHandler Execute")
		http.Redirect(w, r, "/500", http.StatusSeeOther)
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, err := template.ParseFiles("templates/400.html")
	if err != nil {
		fmt.Println("Error BadRequest ParseFiles:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error BadRequest Execute:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func StatusInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("templates/500.html")
	if err != nil {
		fmt.Println("Error InternalServerError ParseFiles:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error InternalServerError Execute:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HandleFavicon(w http.ResponseWriter, r *http.Request) {
	// Répondre avec une réponse vide
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
