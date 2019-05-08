package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Request received")
	fmt.Fprintf(w, "<h1>Hello from k14s-simple-app!</h1>")
}

func main() {
	flag.Parse()
	log.Print("Server started")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
