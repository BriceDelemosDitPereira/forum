package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Like struct {
	ID          int
	Users_id    int
	Posts_id    int
	Comments_id int
	Like_status bool
	// Like_date   time.Time
}

func LikePostsHandler(w http.ResponseWriter, r *http.Request) {
	Check_biscuit(w, r)
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println("Error like.go LikePostsHandler Open", err)

		return
	}
	defer db.Close()

	var query string
	var unused string
	var Like_status int
	var check_post_or_comment int

	r.ParseForm()
	post_id, _ := strconv.Atoi(r.FormValue("post_id"))
	comment_id, _ := strconv.Atoi(r.FormValue("comments_id"))
	like := r.FormValue("post_like")
	dislike := r.FormValue("post_dislike")

	if like == "1" {
		Like_status = 1
	} else if dislike == "-1" {
		Like_status = -1
	}

	if User.ID == 0 {
		fmt.Println("redirect from like.go LikePostsHandler to login.go login_handler")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else {
		// tester si l'utilisateur pour le post liker ou disdliker a deja liker ou disliker
		if comment_id != 0 {
			query = "SELECT Like_status FROM POST_LIKE_DISLIKE WHERE comments_id = ? AND users_id = ?"
			check_post_or_comment = comment_id
		} else {
			query = "SELECT Like_status FROM POST_LIKE_DISLIKE WHERE posts_id = ? AND users_id = ?"
			check_post_or_comment = post_id
		}

		verify, err := db.Prepare(query)
		if err != nil {
			fmt.Println("Error like.go LikePostsHandler Prepare(SELECT)")

			return
		}

		defer verify.Close()

		err = verify.QueryRow(check_post_or_comment, User.ID).Scan(&unused)

		if err != nil {
			fmt.Println("Error like.go LikePostsHandler QueryRow")
		}

		if unused != "" {
			// si deja liker ou disliker faire une update
			if comment_id != 0 {
				query = "UPDATE POST_LIKE_DISLIKE SET like_status = ? WHERE comments_id = ? AND users_id = ?"
			} else {
				query = "UPDATE POST_LIKE_DISLIKE SET like_status = ? WHERE posts_id = ? AND users_id = ?"
			}

			verify, err := db.Prepare(query)
			if err != nil {
				fmt.Println("Error like.go LikePostsHandler Prepare(UPDATE)", err)
			}
			defer verify.Close()

			_, err = verify.Exec(Like_status, check_post_or_comment, User.ID)
			if err != nil {
				fmt.Println("Error like.go LikePostsHandler Exec(UPDATE)", err)
			}

		} else {
			// sinon insert
			addLike, err := db.Prepare("INSERT INTO POST_LIKE_DISLIKE (users_id, posts_id, comments_id, Like_status) VALUES (?, ?, ?, ?)")
			if err != nil {
				fmt.Println("Error like.go LikePostsHandler Prepare(INSERT)", err)
			}

			defer addLike.Close()

			_, err = addLike.Exec(User.ID, post_id, comment_id, Like_status)
			if err != nil {
				fmt.Println("Error like.go LikePostsHandler Exec(INSERT)")
			}
		}
	}
	fmt.Println("redirect from like.go LikePostsHandler to home.go home_connected")
	http.Redirect(w, r, "/home_connected?id="+strconv.Itoa(User.ID), http.StatusSeeOther)
}

func Get_post_by_like(user_id int) (posts []Post) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	queryLike, err := db.Prepare("SELECT posts_id FROM POST_LIKE_DISLIKE WHERE users_id = ? AND like_status = 1")
	if err != nil {
		fmt.Println("error exe recup queryLike from Get_post_by_like", err)
	}
	defer queryLike.Close()
	rowsLike, err := queryLike.Query(user_id)
	if err != nil {
		fmt.Println("error exe scan queryLike from Get_post_by_like", err)
		return
	}
	defer rowsLike.Close()
	for rowsLike.Next() {
		var p Post
		err = rowsLike.Scan(&p.ID)
		if err != nil {
			fmt.Println("error exe scan from Get_post_by_like", err)
			return
		}
		queryLike := `SELECT COUNT(*) FROM POST_LIKE_DISLIKE WHERE like_status = 1 AND posts_id = ?;`
		err := db.QueryRow(queryLike, p.ID).Scan(&p.Like)
		// err = row.Scan(&p.Like)
		if err != nil {
			fmt.Print("In Get all post ", err, p.ID, p.Like, p.Dislike)
		}
		queryDislike := `SELECT COUNT(*) FROM POST_LIKE_DISLIKE WHERE like_status = -1 AND posts_id = ?;`
		err = db.QueryRow(queryDislike, p.ID).Scan(&p.Dislike)
		// err = rowDislike.Scan(&p.Like)
		if err != nil {
			fmt.Print("In Get all post ", err, p.ID, p.Like, p.Dislike)
		}
		// requete de récupération du contenu du post
		query, err := db.Prepare(`SELECT post.title, post.content, post.image, post.publication_date, 
 		users.username
 		FROM post 
 		INNER JOIN users ON users.id = post.users_id
 		WHERE users_id = ? ORDER BY publication_date DESC`)
		if err != nil {
			fmt.Println("error exe recup if from Get_post_by_user_id", err)
		}
		query.QueryRow(p.ID).Scan(&p.Title, &p.Content, &p.Image, &p.Creation_date, &p.Author)
		posts = append(posts, p)
	}
	fmt.Println("POST LIKE : ", posts)
	return posts
}
