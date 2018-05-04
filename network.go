package main

import (
	"database/sql"
	"fmt"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/bootstrap"
	"github.com/nashenmuck/network_server/follow"
	"github.com/nashenmuck/network_server/posts"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			dummyrequest(w, r, db)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			auth.GenAuthToken(w, r, db)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
	http.HandleFunc("/token/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			user, err := auth.CheckAuthToken(w, r, db)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Fprintf(w, "%s\n", user)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
	http.HandleFunc("/token/invite", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			auth.GenInviteToken(w, r, db)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
	http.HandleFunc("/token/reg", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			auth.RegUser(w, r, db, config.NetName)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
	http.HandleFunc("/post/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			posts.Create_post(w, r, db, config.NetName)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
    http.HandleFunc("/post/getfollowing", func (w http.ResponseWriter, r *http.Request) {
        posts.Get_followed_posts(w,r,db)
    })
    http.HandleFunc("/follow/follow", func (w http.ResponseWriter, r *http.Request) {
        follow.Follow_user(w,r,db, config.NetName)
    })
	http.HandleFunc("/post/getall", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			posts.GetAllPosts(w, r, db)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
    http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
        FrontPage(w,r,db)
    })
	log.Printf("Serving on port %s\n", config.SvcPort)
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		http.DefaultServeMux))
}
