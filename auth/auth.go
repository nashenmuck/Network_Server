package auth

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nashenmuck/network_server/networkstructs"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func GenAuthToken(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userpass networkstructs.Users
	err := userpass.Decode(w, r)
	if err != nil {
		return
	}
	readStmt, err := db.Prepare("SELECT password FROM users WHERE username=$1")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	defer readStmt.Close()
	var pass []byte
	err = readStmt.QueryRow(userpass.Username).Scan(&pass)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	err = bcrypt.CompareHashAndPassword(pass, []byte(userpass.Password))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}
	tokenUUID, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	newToken := tokenUUID.Bytes()
	insertStmt, err := db.Prepare("INSERT INTO authtokens(username, token) VALUES($1, $2)")
	defer insertStmt.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	_, err = insertStmt.Exec(userpass.Username, newToken)
	fmt.Fprintf(w, "%x", newToken)
	return
}

func CheckAuthToken(w http.ResponseWriter, r *http.Request, db *sql.DB) bool {
	token := r.Header.Get("Auth-Token")
	user := r.Header.Get("User")
	if token == "" || user == "" {
		http.Error(w, "Auth not provided", 401)
		return false
	}
	stmt, err := db.Prepare("EXISTS(SELECT FROM auth_tokens WHERE username=$1 AND token=$2)")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return false
	}
	defer stmt.Close()
	row, err := stmt.Query(user, token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return false
	}
	if row == nil {
		return false
	} else {
		return true
	}
}
