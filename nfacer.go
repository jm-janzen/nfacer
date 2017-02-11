package main

import (
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {

	// Log IP, URL requested
	log.Println(r.RemoteAddr, r.URL.Path)

	// Serve file at path
	http.ServeFile(w, r, r.URL.Path)
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}
