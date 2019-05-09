package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Request received")
	fmt.Fprintf(w, "<h1>Hello from k8s-simple-app!</h1>")
}

func main() {
	log.Print("Server started")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
