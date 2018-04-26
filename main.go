package main

import (
	"database/sql"
	"fmt"
	"github.com/nashenmuck/network_server/bootstrap"
	_ "github.com/nashenmuck/network_server/json"
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
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		http.DefaultServeMux))
}
