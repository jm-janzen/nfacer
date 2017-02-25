package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/jm-janzen/nfacer/utils/common" // Helper functions specific to this project
	"github.com/yosssi/ace"                    // HTML template engine
)

// TODO move this common obj else, once finished designing
type Data struct {
	Title string
	Other string
	Num   int
}

// Handlers receive and log an HTTP req, then serve our pages (using _render)
func HandleHome(w http.ResponseWriter, r *http.Request) {
	RenderTpl(w, r, "base")
}

func HandleOther(w http.ResponseWriter, r *http.Request) {
	RenderTpl(w, r, "other")
}

func RenderTpl(w http.ResponseWriter, r *http.Request, template string) {
	var requestedPath = r.URL.Path[1:] // Trim leading `/'

	// Print IP, URL, requested path
	log.Println(strings.Split(r.RemoteAddr, ":")[0], strings.Split(r.Host, ":")[0], requestedPath)

	// Load given template by name
	tpl, err := ace.Load("templates/"+template, "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load our Data obj
	data := Data{Title: "jm - " + template}

	// Apply parsed template to w, passing in our Data obj
	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Change dir to project root, if not already there
	common.ChdirWebserverRoot()

	// Handle homepage, other page
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/other", HandleOther)

	// Handle static resources (js, css)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}
