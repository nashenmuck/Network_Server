package follow

import (
    "net/http"
    "encoding/json"
    "database/sql"
    "github.com/nashenmuck/network_server/auth"
    "log"
)

type resp struct {
    Username string `json:"username"`
    Server string `json:"server"`
    When string `json:"when"`
}
func Get_followed(w http.ResponseWriter, r *http.Request, db *sql.DB, nn string) {
    user, err  := auth.CheckAuthToken(w,r,db)
    if err != nil  {
        log.Println(err)
        return
    }
    stmt, err := db.Prepare("SELECT followeeid, Followee_Server, followedwhen from followers  where followerid=$1 AND follower_server=$2")
    if err != nil  {
        log.Println(err)
        return
    }
    row, err := stmt.Query(user, nn)
    if err != nil  {
        log.Println(err)
        return
    }
    w.Write([]byte("["))
    var re resp
    for row.Next() {
        err := row.Scan(&re.Username, &re.Server, &re.When)

        if err != nil  {
            log.Println(err)
            return
        }

        out, err := json.Marshal(re)
        if err != nil {
            log.Println(err)
            return
        }
        w.Write(out)
        w.Write([]byte(","))
    }
    w.Write([]byte("]"))
}
func Get_following(w http.ResponseWriter, r *http.Request, db *sql.DB, nn string) {
    user, err  := auth.CheckAuthToken(w,r,db)
    if err != nil  {
        log.Println(err)
        return
    }
    stmt, err := db.Prepare("SELECT followerid, Follower_Server, followedwhen from followers  where followeeid=$1 AND followee_server=$2")
    if err != nil  {
        log.Println(err)
        return
    }
    row, err := stmt.Query(user, nn)
    if err != nil  {
        log.Println(err)
        return
    }
    w.Write([]byte("["))
    var re resp
    for row.Next() {
        err := row.Scan(&re.Username, &re.Server, &re.When)

        if err != nil  {
            log.Println(err)
            return
        }

        out, err := json.Marshal(re)
        if err != nil {
            log.Println(err)
            return
        }
        w.Write(out)
        w.Write([]byte(","))
    }
    w.Write([]byte("]"))
}
