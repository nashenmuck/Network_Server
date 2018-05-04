package posts

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/netjson"
	"log"
	"net/http"
)

func Create_post(w http.ResponseWriter, r *http.Request, db *sql.DB, nn string) {
	user, err := auth.CheckAuthToken(w, r, db)
	if err != nil {
		log.Println(err)
		return
	}
	var data netjson.Posts
	err = data.Decode(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 400)
		return
	}
	fmt.Println(data)
	defer r.Body.Close()
	var str string
	if data.IsSpecialGroup {
		str = "INSERT INTO POSTS (username, body, special_groupid, is_special_group, origin_server, groupid) values ($1, $2,  $3, TRUE, $4, null)"
	} else {
		str = "INSERT INTO POSTS (username, body, groupid, is_special_group, origin_server) values ($1, $2,  $3, FALSE, $4)"
	}
	st, err := db.Prepare(str)
	if err != nil {
		return
	}
	st.Query(user, data.Body, data.GroupId, nn)
	defer st.Close()
}
