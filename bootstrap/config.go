package bootstrap

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"time"
)

//Generate the database config string as a struct based on environment variables
func DbStringConfig() Config {
	if os.Getenv("DATABASE_URL") == "" {
		port := "8080"
		pghost := "localhost"
		pgport := "5432"
		pguser := "postgres"
		pgdb := "postgres"
		pgpass := ""
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
		} else {
			pgpass = os.Getenv("POSTGRES_PASS")
		}

		uuid, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}
		servername := uuid.String()
		if os.Getenv("NETWORK_SERVERNAME") != "" {
			servername = os.Getenv("NETWORK_SERVERNAME")
		} else {
			servername = "network-" + servername
			log.Printf("Servername not set, using %s", servername)
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
			Pass:     pgpass,
			Database: pgdb,
			NetName:  servername,
			NetAdmin: admin,
			NetPass:  password}
	} else {
		log.Println("Database url set, using")
		return Config{}
	}
}

//Use the generated config string to connect to the database and return a pointer to that connection
func DbConfig(dbconfig Config) *sql.DB {
	log.Println("Connecting to database...")
	var dbconnect string
	if os.Getenv("DATABASE_URL") != "" {
		dbconnect = os.Getenv("DATABASE_URL")
	} else if os.Getenv("POSTGRES_PASS") == "" {
		dbconnect = fmt.Sprintf("user=%s host=%s port=%s sslmode=disable dbname=%s", dbconfig.User, dbconfig.Host, dbconfig.DbPort, dbconfig.Database)
	} else {
		dbconnect = fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable dbname=%s", dbconfig.User, dbconfig.Pass, dbconfig.Host, dbconfig.DbPort, dbconfig.Database)
	}
	db, err := sql.Open("postgres", dbconnect)
	if err != nil {
		log.Fatal(err)
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
