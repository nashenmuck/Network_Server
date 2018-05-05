package bootstrap

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Config struct {
	SvcPort  string
	Host     string
	DbPort   string
	User     string
	Pass     string
	Database string
	NetName  string
	NetAdmin string
	NetPass  string
}

//Perform migrations given the files provided in the `sql` folder
func Dbmigrate(db *sql.DB) {
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

//Insert the admin user and password into the db
func BootstrapAdminAndServer(db *sql.DB, servername string, admin string, password string) {
	log.Println("Bootstrapping admin user...")
	stmt, err := db.Prepare("INSERT INTO servers(server) VALUES($1)")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(servername)
	if err != nil {
		log.Println(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	stmt, err = db.Prepare("INSERT INTO users(username, password, isadmin, canonical_user, canonical_server) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(admin, hash, true, admin, servername)
	if err != nil {
		log.Println(err)
	}
	log.Println("Admin user created!")

}
