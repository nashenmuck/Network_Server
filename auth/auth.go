package auth

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func TokenToUsername(token []byte, db *sql.DB) (string, error) {
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
