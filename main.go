package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
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
	NetAdmin string
	NetPass  string
}

func dummyrequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	db.Ping()
	fmt.Fprintln(w, "We're up!")
}

func config() Config {
	port := "8080"
	pghost := "localhost"
	pgport := "5432"
	pguser := "postgres"
	pgdb := "postgres"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		log.Println("Service port not found, using default 8080")
	}
	if os.Getenv("POSTGRES_HOST") != "" {
		pghost = os.Getenv("POSTGRES_HOST")
	} else {
		log.Printf("Database host not set, using default %s\n", pghost)
	}
	if os.Getenv("POSTGRES_PORT") != "" {
		pgport = os.Getenv("POSTGRES_PORT")
	} else {
		log.Printf("Database port not set, using default %s\n", pgport)
	}
	if os.Getenv("POSTGRES_USER") != "" {
		pguser = os.Getenv("POSTGRES_USER")
	} else {
		log.Printf("Database user not set, using default %s\n", pguser)
	}
	if os.Getenv("POSTGRES_DATABASE") != "" {
		pgdb = os.Getenv("POSTGRES_DATABASE")
	} else {
		log.Printf("Database db not set, using default %s\n", pgdb)
	}
	if os.Getenv("POSTGRES_PASS") == "" {
		log.Println("Database password not set, assuming empty")
	}
	admin, password := "admin", "password"
	if os.Getenv("NETWORK_ADMIN") != "" {
		admin = os.Getenv("NETWORK_ADMIN")
	} else {
		log.Printf("Network admin username not set, using \"%s\"\n", admin)
	}
	if os.Getenv("NETWORK_PASSWORD") != "" {
		password = os.Getenv("NETWORK_PASSWORD")
	} else {
		log.Printf("Network admin password not set, using \"%s\"\n", password)
	}
	return Config{
		SvcPort:  port,
		Host:     pghost,
		DbPort:   pgport,
		User:     pguser,
		Pass:     os.Getenv("POSTGRES_PASS"),
		Database: pgdb,
		NetAdmin: admin,
		NetPass:  password}
}
func dbConfig(dbconfig Config) *sql.DB {
	log.Println("Connecting to database...")
	dbconnect := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable dbname=%s", dbconfig.User, dbconfig.Pass, dbconfig.Host, dbconfig.DbPort, dbconfig.Database)
	if os.Getenv("POSTGRES_PASS") == "" {
		dbconnect = fmt.Sprintf("user=%s host=%s port=%s sslmode=disable dbname=%s", dbconfig.User, dbconfig.Host, dbconfig.DbPort, dbconfig.Database)
	}
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

func dbmigrate(db *sql.DB) {
	log.Println("Migrating SQL definitions...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://sql",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Migrate(2)
	if err != nil {
		log.Println(err)
	}
	log.Println("Migration complete!")
}

func main() {
	config := config()
	db := dbConfig(config)
	dbmigrate(db)
	log.Printf("Serving on port %s\n", config.SvcPort)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dummyrequest(w, r, db)
	})
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
