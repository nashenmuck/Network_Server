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
