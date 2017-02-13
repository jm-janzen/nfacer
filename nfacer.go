package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/yosssi/ace" // HTML template engine
)

// Handler receives and logs an HTTP req, then serves our example base page
func Handler(w http.ResponseWriter, r *http.Request) {
	var requestedPath = r.URL.Path[1:] // Trim leading `/'

	// Print IP (sans port), requested path
	log.Println(strings.Split(r.RemoteAddr, ":")[0], requestedPath)

	// Load base.ace template
	tpl, err := ace.Load("templates/base", "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply parsed template to w
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
