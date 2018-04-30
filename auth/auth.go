package auth
import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "uuid"
)

func token_to_username(token []bytes, db *sql.DB) (string, error) {
    statement, err := db.Prepare("SELECT username FROM authtokens WHERE token=$1")
    if err != nil {
        log.Println(err)
        return nil, true;
    }
    defer statement.Close()
    var user []bytes
    if err := statement.QueryRow(token).scan(&bytes); err != nil {
        log.Println(err)
        return nil, true
    }
    if user == nil {
        return nil, true
    }
    return user, nil
}

func validate_request(w http.ResponseWriter, r *http.Request, db *sql.DB) (string, error) {
    auth_token = r.Header.Get("Auth-Token")
    if auth_token == "" {
        http.Error(w, "No auth token", 401)
        return nil, true
    }
    uuid_tok = uuid.FromString(auth_token)
    user, err = token_to_username(uuid_tok.Bytes())
    if err != nil {
        http.Error(w, "Bad auth", 401)
        return nil, true
    }
    return user, nil
}
