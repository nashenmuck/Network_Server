package main

import (
	"database/sql"
	"fmt"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/bootstrap"
	"log"
	"net/http"
)

func dummyrequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	db.Ping()
	fmt.Fprintln(w, "We're up!")
}

func main() {
	config := bootstrap.DbStringConfig()
	db := bootstrap.DbConfig(config)
	bootstrap.Dbmigrate(db)
	bootstrap.BootstrapAdminAndServer(db, config.NetName, config.NetAdmin, config.NetPass)
	log.Printf("Serving on port %s\n", config.SvcPort)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dummyrequest(w, r, db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.GenAuthToken(w, r, db)
	})
	http.HandleFunc("/testtoken", func(w http.ResponseWriter, r *http.Request) {
		isAuthed := auth.CheckAuthToken(w, r, db)
		fmt.Fprintf(w, "%t", isAuthed)
	})

	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		http.DefaultServeMux))
}
