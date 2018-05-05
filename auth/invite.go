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
	user, err := CheckAuthToken(w, r, db)
	if err != nil {
		log.Println(err)
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
	_, err = insertStmt.Exec(user, newToken)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	fmt.Fprintf(w, "%s\n", tokenUUID.String())
	return
}

func CheckRegToken(w http.ResponseWriter, r *http.Request, db *sql.DB) bool {
	token := r.Header.Get("Reg-Token")
	if token == "" {
		http.Error(w, "Auth not provided", 401)
		return false
	}
	stmt, err := db.Prepare("SELECT EXISTS(SELECT * FROM regtokens WHERE token=$1)")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return false
	}
	defer stmt.Close()
	var isAuthed bool
	stringTok, err := uuid.FromString(token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", 401)
		return false
	}
	err = stmt.QueryRow(stringTok.Bytes()).Scan(&isAuthed)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return false
	}
	if !isAuthed {
		log.Println(err)
		http.Error(w, "Unauthorized", 401)
		return false
	}
	return isAuthed
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
