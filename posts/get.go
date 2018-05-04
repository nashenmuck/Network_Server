package posts

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/netjson"
	"log"
	"net/http"
	//    "time"
	//    "fmt"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, err := auth.CheckAuthToken(w, r, db)
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := db.Prepare("SELECT username, body, origin_server, date FROM posts WHERE special_groupid=1")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	row, err := stmt.Query()
	w.Header().Set("content-type", "application/json")
	for row.Next() {
		var post netjson.Posts
		err := row.Scan(&post.Username, &post.Body, &post.OriginServer, &post.Date)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", 500)
			return
		}
		output, err := json.Marshal(post)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", 500)
			return
		}
		w.Write(output)
	}
}

type RetPost struct {
	PostID   int    `json:"postid"`
	Username string `json:"username"`
	Body     string `json:"body"`
}

func Get_followed_posts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	user, err := auth.CheckAuthToken(w, r, db)
	if err != nil {
		log.Println(err)
		return
	}
	dec := json.NewDecoder(r.Body)
	var data struct {
		Since int `json:"since"`
	}
	err = dec.Decode(&data)
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := db.Prepare("SELECT id, username, body FROM posts WHERE date > to_timestamp($1) AND posts.username in (SELECT followeeId from followers where followerId = $2) AND (posts.is_special_group = TRUE OR $2 in (select follower from group_followers where group_id = posts.groupid))")
	if err != nil {
		log.Println(err)
		return
	}
	// tm, err := time.Parse(time.UnixDate, string(data.Since))
	row, err := stmt.Query(data.Since, user)
	w.Header().Set("content-type", "application/json")
	var post RetPost
	for row.Next() {
		err := row.Scan(&post.PostID, &post.Username, &post.Body)
		if err != nil {
			log.Println(err)
			return
		}
		out, err := json.Marshal(post)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(out)
	}
}
