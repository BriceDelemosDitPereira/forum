package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	ID            int
	User_id       int
	Post_id       int
	Content       string
	Like          int
	Dislike       int
	Creation_date time.Time
}

func Create_comment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/create_comment" {
		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			fmt.Println("Error comment.go Create_comment Open", err)

			return
		}
		defer db.Close()
		content := r.FormValue("content")
		post_id := r.FormValue("post_id")
		time := time.Now()
		insert_post, err := db.Prepare("INSERT INTO post_comment (users_id, post_id, content, creation_date) VALUES (?, ?, ?, ?)")
		if err != nil {
			fmt.Println("Error comment.go Create_comment Prepare(INSERT)", err)

			return
		}
		defer insert_post.Close()

		_, err = insert_post.Exec(User.ID, post_id, content, time)
		if err != nil {
			fmt.Println("Error comment.go Create_comment Exec(INSERT)", err)

			return
		}
		fmt.Println("redirect from comment.go Create_comment to home.go home_connected")
		http.Redirect(w, r, "/home_connected?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
		return
	} else {
		fmt.Println("redirect from comment.go Create_comment to 404 error")
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
}

func (p Post) Get_comment_from_post() (commentsMap map[int][]Comment, erreur error) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error comment.go Get_comment_from_post Open")
		return nil, err
	}
	defer db.Close()

	commentsMap = make(map[int][]Comment)

	rows, err := db.Query("SELECT id, users_id, content, creation_date FROM post_comment WHERE post_id=?", p.ID)
	if err != nil {
		fmt.Println("Error comment.go Get_comment_from_post Query(SELECT)", err)
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.ID, &c.User_id, &c.Content, &c.Creation_date)
		if err != nil {
			fmt.Println("Error comment.go Get_comment_from_post Scan", err)
			return nil, err
		}

		queryLike := `SELECT COUNT(*) FROM POST_LIKE_DISLIKE WHERE like_status = 1 AND comments_id = ?;`
		err = db.QueryRow(queryLike, c.ID).Scan(&c.Like)
		if err != nil {
			fmt.Println("Error comment.go Get_comment_from_post QueryRow(SELECT)", err)
		}
		queryDislike := `SELECT COUNT(*) FROM POST_LIKE_DISLIKE WHERE like_status = -1 AND comments_id = ?;`
		err = db.QueryRow(queryDislike, c.ID).Scan(&c.Dislike)

		if err != nil {
			fmt.Println("Error comment.go Get_comment_from_post QueryRow(SELECT)", err)
			return nil, err
		}

		commentsMap[c.Post_id] = append(commentsMap[c.Post_id], c)
	}
	return commentsMap, nil
}

/*
Lorsqu'on définit une méthode en Go, on peut définir si elle doit recevoir un pointeur sur une valeur ou directement une valeur en utilisant *Post ou Post.
Dans le premier cas (p *Post), on utilise un pointeur sur une valeur Post. Cela signifie que la méthode peut modifier directement la valeur pointée sans avoir à la renvoyer.
Dans le deuxième cas (p Post), on utilise directement une valeur Post. Cela signifie que la méthode travaille sur une copie de la valeur. Si la méthode modifie la valeur, cela n'affectera pas la valeur originale.
Dans notre cas, on a une méthode qui crée et renvoie une nouvelle valeur commentsMap. Il n'est donc pas nécessaire d'utiliser un pointeur sur une valeur Post. On peut utiliser directement une valeur Post.
C'est pourquoi il n'y a pas de différence entre utiliser p *Post ou p Post dans notre cas. Les deux versions fonctionneront de la même manière, car la méthode ne modifie pas la valeur Post originale.
*/
