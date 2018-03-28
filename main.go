package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	SvcPort  string
	Host     string
	DbPort   string
	User     string
	Pass     string
	Database string
}

func dummyrequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "We're up!")
}

func config() Config {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		log.Println("Service port not found, using default 8080")
	}
	if os.Getenv("POSTGRES_HOST") == "" {
		log.Fatal("Database host not set!")
	}
	if os.Getenv("POSTGRES_PORT") == "" {
		log.Fatal("Database port not set!")
	}
	if os.Getenv("POSTGRES_USER") == "" {
		log.Fatal("Database user not set!")
	}
	if os.Getenv("POSTGRES_PASS") == "" {
		log.Printf("Database password not set, assuming empty\n")
	}
	return Config{
		SvcPort:  port,
		Host:     os.Getenv("POSTGRES_HOST"),
		DbPort:   os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Pass:     os.Getenv("POSTGRES_PASS"),
		Database: "network"}
}
func dbConfig(dbconfig Config) *sql.DB {
	log.Println("Connecting to database...")
	dbconnect := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable", dbconfig.User, dbconfig.Pass, dbconfig.Host, dbconfig.DbPort)
	db, err := sql.Open("postgres", dbconnect)
	if err != nil {
		log.Println(err)
	}
	var dbError error
	for attempts := 1; attempts <= 20; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Println(dbError)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		log.Fatal(dbError)
	}
	defer log.Println("Connected to database!")
	return db
}
func main() {
	config := config()
	_ = dbConfig(config)
	log.Printf("Serving on port %s\n", config.SvcPort)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dummyrequest(w, r)
	})
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
