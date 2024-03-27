package handler

import (
	"database/sql"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID            int
	User_id       int
	Title         string
	Content       string
	Creation_date time.Time
	Author        string
	Image         string
	Like          int
	Dislike       int
	CommentsMap   map[int][]Comment
}

func Create_post(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error post.go Create_post Open", err)

		return
	}
	defer db.Close()
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["category"]
	time := time.Now()
	path_image := UploadPictures(w, r)
	if title == "" || content == "" {
		fmt.Println("redirect from post.go Create_post to home.go home_connected title|content = \"\"")
		http.Redirect(w, r, "/home_connected?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
		return
	} else {
		stmt, err := db.Prepare("INSERT INTO post (users_id, title, content, image, publication_date) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("Error post.go Create_post Prepare(INSERT)", err)

			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(User.ID, title, content, path_image, time)
		if err != nil {
			fmt.Println("Error post.go Create_post Exec(INSERT)", err)

			return
		}

		Insert_categories(categories, db, title, content)

		fmt.Println("redirect from post.go Create_post to home.go home_connected END")
		http.Redirect(w, r, "/home_connected?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
		return
	}
}

func UploadPictures(w http.ResponseWriter, r *http.Request) string {
	_, err := os.Stat("./img_upload/")
	if err != nil {
		return ""
	}
	if os.IsNotExist(err) {
		// Créer le répertoire seulement s'il n'existe pas
		errDir := os.MkdirAll("./img_upload/", os.ModePerm)
		if errDir != nil {
			return ""
		}
	}

	image_data, image_info, err := r.FormFile("image")
	if err != nil {
		return ""
	}

	if image_info.Size > 20000000 {
		w.Write([]byte("Vous ne pouvez pas upload une image de plus de 20mb"))
	}

	split_name := strings.Split(image_info.Filename, ".")
	file_type := strings.ToLower(split_name[1])

	file, err := os.Create("./img_upload/" + image_info.Filename)
	if err != nil {
		fmt.Println("Error post.go Create file", err)
		return ""
	}
	defer file.Close()

	switch file_type {
	case "png":
		img, err := png.Decode(image_data)
		if err != nil {
			fmt.Println("Error post.go Decode png", err)
			return ""
		}
		err = png.Encode(file, img)
		if err != nil {
			fmt.Println("Error post.go Encode png", err)
			return ""
		}
	case "jpeg", "jpg":
		img, err := jpeg.Decode(image_data)
		if err != nil {
			fmt.Println("Error post.go Decode jpeg", err)
			return ""
		}
		opt := jpeg.Options{
			Quality: 90,
		}
		err = jpeg.Encode(file, img, &opt)
		if err != nil {
			fmt.Println("Error post.go Encode jpeg", err)
			return ""
		}

	case "gif":
		img, err := gif.DecodeAll(image_data)
		if err != nil {
			fmt.Println("Error post.go Decode gif", err)
			return ""
		}
		err = gif.EncodeAll(file, img)
		if err != nil {
			fmt.Println("Error post.go Decode gif", err)
			return ""
		}
	}
	image_data.Close()
	return "/img_upload/" + image_info.Filename
}

func Get_all_posts(r *http.Request) ([]Post, error) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error post.go Get_all_posts Open")
		return nil, err
	}
	defer db.Close()
	r.ParseForm() // Permet de lire le formulaire
	filter_cat := r.Form["filter"]

	rows, err := db.Query("SELECT id, users_id, title, content, image, publication_date FROM post ORDER BY publication_date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.User_id, &p.Title, &p.Content, &p.Image, &p.Creation_date)
		if err != nil {
			return nil, err
		}

		queryLike := `SELECT COUNT(*) FROM POST_LIKE_DISLIKE WHERE like_status = 1 AND posts_id = ?;`
		err := db.QueryRow(queryLike, p.ID).Scan(&p.Like)
		if err != nil {
			fmt.Println("Error post.go Get_all_posts QueryRow(queryLike)")
			return nil, err
		}
		queryDislike := `SELECT COUNT(*) FROM POST_LIKE_DISLIKE WHERE like_status = -1 AND posts_id = ?;`
		err = db.QueryRow(queryDislike, p.ID).Scan(&p.Dislike)
		if err != nil {
			fmt.Println("Error post.go Get_all_posts QueryRow(queryDislike)")
			return nil, err
		}

		p.CommentsMap, err = p.Get_comment_from_post()
		if err != nil {
			return nil, err
		}
		categories := Get_categories_by_post_id(p.ID)
		if len(filter_cat) != 0 {
		outerLoop: // Permet de break la boucle juste après outerLoop (c'est forcément ce nom)
			for _, cat := range filter_cat {
				for i := 0; i < len(categories); i++ {
					if strings.Contains(cat, categories[i]) {
						posts = append(posts, p)
						break outerLoop
					}
				}
			}
		} else {
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func Get_post_by_user_id(user_id int) (posts []Post, erreur error) {
	username, _ := Get_username_by_id(user_id)
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error post.go Get_post_by_user_id Open")
		return nil, err
	}
	defer db.Close()
	query, err := db.Prepare("SELECT id,title, content, image, publication_date FROM post WHERE users_id = ? ORDER BY publication_date DESC")
	if err != nil {
		fmt.Println("error exe recup if from Get_post_by_user_id", err)
	}
	defer query.Close()
	rows, err := query.Query(user_id)
	if err != nil {
		fmt.Println("error exe scan from Get_post_by_user_id", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Creation_date)
		if err != nil {
			fmt.Println("error exe scan from Get_post_by_user_id", err)
			return
		}
		p.Author = username
		posts = append(posts, p) // categories déjà initialisé dans la ligne de la fonction (à la fin)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("error exe rows from Get_post_by_user_id", err)
		return

	}
	// categories contient maintenant toutes les catégories associées au post_id
	fmt.Println("POST : ", posts)
	return posts, nil
}
