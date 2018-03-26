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
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	fmt.Printf("Serving on port %s\n", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dummyrequest(w, r)
	})
	log.Fatal(http.ListenAndServe(":"+port,
		handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
