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
	http.HandleFunc("/token/test", func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.CheckAuthToken(w, r, db)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "%s", user)
	})
	http.HandleFunc("/token/invite", func(w http.ResponseWriter, r *http.Request) {
		auth.GenInviteToken(w, r, db)
	})
	http.HandleFunc("/token/reg", func(w http.ResponseWriter, r *http.Request) {
		auth.RegUser(w, r, db, config.NetName)
	})
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		http.DefaultServeMux))
}
