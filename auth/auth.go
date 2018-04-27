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
