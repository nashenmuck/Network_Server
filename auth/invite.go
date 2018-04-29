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

func GenInviteToken(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !CheckAuthToken(w, r, db) {
		return
	}
	tokenUUID, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	newToken := tokenUUID.Bytes()
	insertStmt, err := db.Prepare("INSERT INTO regtokens(issuer, token) VALUES($1, $2)")
	defer insertStmt.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	_, err = insertStmt.Exec(r.Header.Get("User"), newToken)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	fmt.Fprintf(w, "%s", tokenUUID.String())
	return
}

func RegUser(w http.ResponseWriter, r *http.Request, db *sql.DB, servername string) {
	if !CheckRegToken(w, r, db) {
		return
	}
	var newUser networkstructs.Users
	err := newUser.Decode(w, r)
	if err != nil {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	stmt, err := db.Prepare("INSERT INTO users(username, password, isadmin, canonical_user, canonical_server) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(newUser.Username, hash, false, newUser.Username, servername)
	fmt.Fprintln(w, "Registered")
	return
}
