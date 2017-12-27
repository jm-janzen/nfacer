package main

import (
	"log"
	"net" // Get host, port strings
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
func Route(w http.ResponseWriter, r *http.Request) {

	// Get host and port values
	// FIXME will err on IPv6 addresses
	if host, port, err := net.SplitHostPort(r.Host); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return

	} else { // Switch on host
		log.Printf("host: %v, port: %v", host, port)
		HandleDefault(w, r)
	}
}

// Render whatever template is present at "/templates/bodies/{resource}.ace",
// else write verbose error to w
func HandleDefault(w http.ResponseWriter, r *http.Request) {
	var requestedPath = r.URL.Path[1:] // Trim leading `/'
	log.Println("Requested Path: " + requestedPath)

	// TODO eventually replace with AJAX
	prefix := "bodies/"

	// Default to rendering home template
	if requestedPath == "" {
		requestedPath = "home"
	}

	var pageTitle string = requestedPath

	RenderTpl(w, r, prefix+requestedPath, pageTitle)

}

// Renders given template by name (string)
// FIXME handle errors better once dev settles
func RenderTpl(w http.ResponseWriter, r *http.Request, template string, pageTitle string) {
	var requestedPath = r.URL.Path[1:] // Trim leading `/'

	// Print IP, URL, requested path; path to template file
	log.Println(strings.Split(r.RemoteAddr, ":")[0], strings.Split(r.Host, ":")[0], requestedPath)
	log.Println("Serving template:", "templates/"+template)

	// Load given template by name
	tpl, err := ace.Load("templates/"+template, "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error:", err.Error())
		return
	}

	// Load our Data obj
	data := Data{Title: "jm - " + pageTitle}

	// Apply parsed template to w, passing in our Data obj
	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error:", err.Error())
		return
	}
}

func main() {
	// Change dir to project root, if not already there
	common.ChdirWebserverRoot()

	mux := http.NewServeMux()

	// Handle homepage, others
	mux.HandleFunc("/", Route)

	// Handle static resources (js, css)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle favicon (anon func)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/fav/favicon.ico")
	})

	// If any issue starting, log err, and exit(1)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
