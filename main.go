package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func dummyrequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	db.Ping()
	fmt.Fprintln(w, "We're up!")
}

func main() {
	config := dbStringConfig()
	db := dbConfig(config)
	dbmigrate(db)
	bootstrapAdminAndServer(db, config.NetName, config.NetAdmin, config.NetPass)
	log.Printf("Serving on port %s\n", config.SvcPort)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dummyrequest(w, r, db)
	})
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
