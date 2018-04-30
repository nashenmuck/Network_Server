package auth

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nashenmuck/network_server/netjson"
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
	var newUser netjson.Users
	err := newUser.Decode(w, r)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		log.Println(err)
		return
	}
	insertStmt, err := db.Prepare("INSERT INTO users(username, password, isadmin, canonical_user, canonical_server) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		http.Error(w, "Internal server error", 500)
		log.Println(err)
		return
	}
	defer insertStmt.Close()
	_, err = insertStmt.Exec(newUser.Username, hash, false, newUser.Username, servername)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		log.Println(err)
		return
	}
	deleteStmt, err := db.Prepare("DELETE FROM regtokens WHERE token=$1")
	if err != nil {
		http.Error(w, "Internal server error", 500)
		log.Println(err)
		return
	}
	defer deleteStmt.Close()
	stringTok, err := uuid.FromString(r.Header.Get("Reg-Token"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	_, err = deleteStmt.Exec(stringTok.Bytes())
	if err != nil {
		http.Error(w, "Internal server error", 500)
		log.Println(err)
		return
	}
	fmt.Fprintln(w, "Registered")
	return
}
