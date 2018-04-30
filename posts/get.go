package posts

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/netjson"
	"log"
	"net/http"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, err := auth.CheckAuthToken(w, r, db)
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := db.Prepare("SELECT username, body, origin_server, date FROM posts WHERE groupid=1")
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
