package main

import (
	"log"
	"net" // Get host, port strings
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jm-janzen/nfacer/utils/common" // Helper functions specific to this project
	"github.com/yosssi/ace"                    // HTML template engine
)

// TODO move this common obj else, once finished designing
type Data struct {
	Title string
	Other string
	Num   int
}

// Define our own simple HTTP request multiplexer
type Mux struct {
	http.Handler  // Interface to ServeHTTP(w, *r)
}

// Handlers receive and log an HTTP req, then serve our pages (using _render)
func Route(w http.ResponseWriter, r *http.Request) {

	// Get host and port values
	// FIXME will err on IPv6 addresses
	if host, port, err := net.SplitHostPort(r.Host); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return

	} else { // Switch on host
		log.Printf("host: %v, port: %v", host, port)
		switch host {
		case "jmjanzen.com", "www.jmjanzen.com":
			RenderTpl(w, r, "base")
		case "api.nullportal.com":
			w.Write([]byte("welcome to nullportal api"))
		case "nullportal.com", "www.nullportal.com":
			w.Write([]byte("welcome to nullportal"))
		}
	}
}

// Handles /other, TODO implement a more programmatic solution
func HandleOther(w http.ResponseWriter, r *http.Request) {
	RenderTpl(w, r, "other")
}

// Renders given template by name (string)
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

/*
 * TODO set up separate muxes for each host, subdomain (eg: *.nullportal.com, api.*.com, etc)
 */
func main() {
	// Change dir to project root, if not already there
	common.ChdirWebserverRoot()

	r := mux.NewRouter()
	r.HandleFunc("/", Route)

	s := r.Host("jmjanzen.com:8080").Subrouter()
	s.HandleFunc("/", Route)

	http.ListenAndServe(":8080", s)

	/*
	mux := http.NewServeMux()

	// Handle homepage, other page
	mux.HandleFunc("/", Route)
	mux.HandleFunc("/other", HandleOther)

	// Handle static resources (js, css)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", mux)
	*/
}
