package main

import (
	"log"
	"net/http"

	"github.com/jm-janzen/nfacer/utils/common" // Helper functions specific to this project
	"github.com/yosssi/ace"                    // HTML template engine
)

type Data struct {
	Title string
	Other string
	Num   int
}

// Handlers receive and log an HTTP req, then serve our pages (using _render)
func Route(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v %v%v", r.RemoteAddr, r.Method, r.Host, r.RequestURI)
	HandleDefault(w, r)
}

// Render whatever template is present at "/templates/bodies/{resource}.ace",
// else write verbose error to w
func HandleDefault(w http.ResponseWriter, r *http.Request) {
	var requestedPath = r.URL.Path[1:] // Trim leading `/'

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

	// Load given template by name
	tpl, err := ace.Load("templates/"+template, "", nil)
	if err != nil {

		// Invalid resource - hardcode to redirect to 404 page
		log.Println("Error:", err.Error(), "trying 404 instead")
		pageTitle, template = "not found", "404"

		// If this fails for some reason, just quit
		if tpl, err = ace.Load("templates/bodies/404", "", nil); err != nil {
			log.Println("Error:", err.Error())
			return
		}
	}

	// Print IP, URL, requested path; path to template file
	log.Println("Serving template:", "templates/bodies/"+template)

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

	// Handle favicon, home img (anon func)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/fav/favicon.ico")
	})
	mux.HandleFunc("/home.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/home.png")
	})
	mux.HandleFunc("/me.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/me.png")
	})

	// If any issue starting, log err, and exit(1)
	listenPort := ":6060"
	log.Printf("Listening on port %v", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, mux))
}
