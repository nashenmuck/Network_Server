package auth

import (
	"database/sql"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

func token_to_username(token []byte, db *sql.DB) (string, error) {
	statement, err := db.Prepare("SELECT username FROM authtokens WHERE token=$1")
	var user string
	if err != nil {
		log.Println(err)
		return user, err
	}
	defer statement.Close()
	if err := statement.QueryRow(token).Scan(&user); err != nil {
		log.Println(err)
		return user, err
	}
	if user == "" {
		return user, err
	}
	return user, nil
}

func validate_request(w http.ResponseWriter, r *http.Request, db *sql.DB) (string, error) {
	auth_token := r.Header.Get("Auth-Token")
	if auth_token == "" {
		http.Error(w, "No auth token", 401)
		return "", fmt.Errorf("No auth token")
	}
	uuid_tok, err := uuid.FromString(auth_token)
	if err != nil {
		http.Error(w, "Bad auth", 401)
		return "", err
	}
	user, err := token_to_username(uuid_tok.Bytes(), db)
	if err != nil {
		http.Error(w, "Bad auth", 401)
		return user, err
	}
	return user, nil
}
