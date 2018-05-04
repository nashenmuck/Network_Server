package posts
import (
    "net/http"
    "encoding/json"
    "database/sql"
    "github.com/nashenmuck/network_server/auth"
    "fmt"
)

type NewPost struct {
    Target int `json:"target"`
    Special bool `json:"special"`
    Post string `json:"post"`
}
func Create_post(w http.ResponseWriter, r *http.Request, db *sql.DB, nn string) {
    user , err := auth.CheckAuthToken(w,r,db)
    if err != nil {
        return
    }
    dec := json.NewDecoder(r.Body)
    var data NewPost
    err = dec.Decode(&data)
    if err != nil {
        return
    }
    fmt.Println(data)
    defer r.Body.Close()
    var str string
    if data.Special {
        fmt.Println("AAAA")
        str = "INSERT INTO POSTS (username, body, special_groupid, is_special_group, origin_server, groupid) values ($1, $2,  $3, TRUE, $4, null)"
    } else {
        fmt.Printf("BBBB")
        str = "INSERT INTO POSTS (username, body, groupid, is_special_group, origin_server) values ($1, $2,  $3, FALSE, $4)"
    }
    st, err := db.Prepare(str)
    if err != nil {
        return
    }
    st.Query(user, data.Post, data.Target, nn)
}
