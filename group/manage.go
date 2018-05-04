package group

import (
    "database/sql"
    "net/http"
    "encoding/json"
    "log"
    "github.com/nashenmuck/network_server/auth"
)

func Create_group(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    user, err := auth.CheckAuthToken(w, r, db)
    if err != nil {
        log.Println(err)
        return
    }
    dec := json.NewDecoder(r.Body)
    var data struct {
        Name string `json:"name"`
    }
    err = dec.Decode(&data)
    if err != nil {
        log.Println(err)
        return
    }
    stmt, err :=  db.Prepare("INSERT INTO GROUPS (owner, group_name) VALUES ($1, $2)")
    if err != nil {
        log.Println(err)
        return
    }
    defer stmt.Close()
    stmt.Query(user, data.Name)
    var id string
    db.QueryRow("SELECT group_id from groups where owner=$1 AND group_name=$2", user, data.Name).Scan(&id)
    w.Write([]byte("{\"id\":"+id+"}"))
}

func Assign_user_to_group(w http.ResponseWriter, r *http.Request, db *sql.DB, nn string) {
    _, err := auth.CheckAuthToken(w, r, db)
    if err != nil {
        log.Println(err)
        return
    }
    dec := json.NewDecoder(r.Body)
    var data struct {
        Username string `json:"username"`
        Server string `json:"server"`
        Groupid int `json:"group_id"`
    }
    err = dec.Decode(&data)
    if err != nil {
        log.Println(err)
        return
    }
    tnet := nn
    if data.Server != "" {
        tnet =  data.Server
    }
    stmt, err :=  db.Prepare("INSERT INTO group_followers (group_id, follower, follower_server) VALUES ($1, $2, $3)")
    if err != nil {
        log.Println(err)
        return
    }
    defer stmt.Close()
    //TODO: Make sure that the target username is a follower. This shouldn't matter,
    //as code that handles groups should check to  see if the user is a follower anyways
    stmt.Query(data.Groupid, data.Username, tnet)

}
