package main

import (
	"database/sql"
	"fmt"
	"github.com/nashenmuck/network_server/auth"
	"github.com/nashenmuck/network_server/bootstrap"
	"github.com/nashenmuck/network_server/follow"
	"github.com/nashenmuck/network_server/posts"
	"github.com/nashenmuck/network_server/group"
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
		if r.Method == "GET" && r.URL.Path == "/" {
			dummyrequest(w, r, db)
		} else if r.URL.Path != "/" {
			http.Error(w, "Not found", 404)
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
	http.HandleFunc("/post/getfollowing", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			posts.Get_followed_posts(w, r, db)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
    http.HandleFunc("/follow/getfollowed", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			follow.Get_followed(w, r, db, config.NetName)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
    http.HandleFunc("/follow/getfollowers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			follow.Get_following(w, r, db, config.NetName)
		} else {
			http.Error(w, "Invalid method", 405)
		}
	})
	http.HandleFunc("/follow/follow", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			follow.Follow_user(w, r, db, config.NetName)
		} else {
			http.Error(w, "Invalid method", 405)
		}
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
    http.HandleFunc("/group/create", func(w http.ResponseWriter, r *http.Request) {
        group.Create_group(w,r,db)
    })
	log.Printf("Serving on port %s\n", config.SvcPort)
	log.Fatal(http.ListenAndServe(":"+config.SvcPort,
		http.DefaultServeMux))
}
