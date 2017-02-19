package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"                   // HTTP request  router
	"github.com/jm-janzen/nfacer/utils/common" // Helper functions specific to this project
	"github.com/yosssi/ace"                    // HTML template engine
)

// Handlers receive and log an HTTP req, then serve our pages (using _render)
func HandleHome(w http.ResponseWriter, r *http.Request) {
	RenderTpl(w, r, "base")
}

func HandleOther(w http.ResponseWriter, r *http.Request) {
	RenderTpl(w, r, "other")
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dir := vars["dir"]
	file := vars["file"]

	// Set Content-Type in response header
	var contentType string
	if strings.HasSuffix(file, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(file, ".js") {
		contentType = "text/javascript"
	}
	w.Header().Add("Content-Type", contentType)

	// Log request
	log.Printf("static/%s/%s (%s)", dir, file, contentType)

	// Serve files from `/static/{css,js/'
	go http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

}

func RenderTpl(w http.ResponseWriter, r *http.Request, template string) {
	var requestedPath = r.URL.Path[1:] // Trim leading `/'

	// Print IP (sans port), requested path
	log.Println(strings.Split(r.RemoteAddr, ":")[0], requestedPath)

	// Load given template by name
	tpl, err := ace.Load("templates/"+template, "", nil)
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
	// Change dir to project root, if not already there
	common.ChdirWebserverRoot()

	r := mux.NewRouter()

	r.HandleFunc("/", HandleHome)
	r.HandleFunc("/other", HandleOther)
	r.HandleFunc("/static/{dir:(?:css|js)}/{file}", HandleStatic)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
