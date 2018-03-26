package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func dummyrequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "We're up!")
}

func main() {
	fmt.Printf("Serving on port %s\n", os.Getenv("PORT"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dummyrequest(w, r)
	})
	log.Fatal(http.ListenAndServe(":"+(os.Getenv("PORT")),
		handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
