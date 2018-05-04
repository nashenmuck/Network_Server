package follow

import (
	"database/sql"
	"encoding/json"
	"github.com/nashenmuck/network_server/auth"
	"net/http"
)

type FollowUserStruct struct {
	Username string `json:"username"`
	Server   string `json:"server"`
}

func Follow_user(w http.ResponseWriter, r *http.Request, db *sql.DB, nn string) {
	user, err := auth.CheckAuthToken(w, r, db)
	if err != nil {
		return
	}
	dec := json.NewDecoder(r.Body)
	var dat FollowUserStruct
	err = dec.Decode(&dat)
	if err != nil {
		return
	}
	defer r.Body.Close()
	stmt, err := db.Prepare("INSERT INTO followers (followeeId, followee_server, followerId, follower_server) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return
	}
	defer stmt.Close()
	tsn := nn
	if dat.Server != "" {
		tsn = dat.Server
	}
	stmt.Query(dat.Username, tsn, user, nn)
}
