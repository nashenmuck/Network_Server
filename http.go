package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
)

type RS struct {
	Author string
	Post   string
	When   string
}

func FrontPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var buf bytes.Buffer
	buf.WriteString("<html><head><link rel=\"stylesheet\" href=\"https://unpkg.com/sakura.css/css/sakura.css\" type=\"text/css\"></head><body><ul>")
	stmt, err := db.Prepare("SELECT username, body, date_trunc('second', date) from posts where is_special_group = TRUE and special_groupid = 1 order by date desc limit 10")
	if err != nil {
		log.Println(err)
		return
	}
	var post RS
	row, err := stmt.Query()
	defer stmt.Close()
	for row.Next() {
		err = row.Scan(&post.Author, &post.Post, &post.When)
		if err != nil {
			log.Println(err)
			return
		}
		st := "<li><h3>" + post.Author + "</h3><p>" + post.Post + "</p><p> at " + post.When + "</p></li>"
		buf.WriteString(st)
	}
	buf.WriteString("</ul></body></html>")
	w.Write([]byte(buf.String()))
}
