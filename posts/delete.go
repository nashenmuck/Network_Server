package posts

import (
	"database/sql"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/netjson"
	"log"
	"net/http"
)

func DeletePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	user, err := auth.CheckAuthToken(w, r, db)
	if err != nil {
		log.Println(err)
		return
	}
	var toDelete netjson.Posts
	err = toDelete.Decode(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := db.Prepare("DELETE FROM posts WHERE username=$1 AND id=$2")
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
	_, err = stmt.Exec(user, toDelete.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
}
