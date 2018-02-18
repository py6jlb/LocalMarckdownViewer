package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var (
	postTemplate  = template.Must(template.ParseFiles(path.Join("assets", "templates", "layout.html"), path.Join("assets", "templates", "post.html")))
	errorTemplate = template.Must(template.ParseFiles(path.Join("assets", "templates", "layout.html"), path.Join("assets", "templates", "error.html")))
	notes         = newNotesCollection()
)

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	mux := mux.NewRouter()
	staticHandler := http.StripPrefix("/static", noDirListing(http.FileServer(http.Dir("assets"))))
	mux.PathPrefix("/static/").Handler(staticHandler)

	mux.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { noteHandler(w, r, cfg) }))

	http.Handle("/", mux)
	listen := cfg.ListenAddress + ":" + strconv.Itoa(cfg.ListenPort)
	log.Println("Open in browser: http://" + listen)
	http.ListenAndServe(listen, nil)
}

func noteHandler(w http.ResponseWriter, r *http.Request, cfg *config) {
	page := r.URL.Path
	p := path.Join(cfg.DataPath, page)
	var postMd string
	if page != "/" {
		postMd = p + ".md"
	} else {
		postMd = p + "/" + cfg.IndexFileName
	}
	post, status, err := notes.getNote(postMd)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := postTemplate.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := errorTemplate.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "" {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
