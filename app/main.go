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
	post_template  = template.Must(template.ParseFiles(path.Join("app", "templates", "layout.html"), path.Join("app", "templates", "post.html")))
	error_template = template.Must(template.ParseFiles(path.Join("app", "templates", "layout.html"), path.Join("app", "templates", "error.html")))
	notes          = newNotesCollection()
)

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	mux := mux.NewRouter()
	s := http.StripPrefix("/static/", noDirListing(http.FileServer(http.Dir("app/assets"))))
	mux.PathPrefix("/static/").Handler(s)
	mux.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { noteHandler(w, r, cfg) }))

	http.Handle("/", mux)
	listen := cfg.ListenAddress + ":" + strconv.Itoa(cfg.ListenPort)
	log.Println("Open in browser: http://" + listen)
	http.ListenAndServe(listen, nil)
}

func noteHandler(w http.ResponseWriter, r *http.Request, cfg *config) {
	page := r.URL.Path
	p := path.Join(cfg.DataPath, page)
	var post_md string
	if page != "/" {
		post_md = p + ".md"
	} else {
		post_md = p + "/" + cfg.IndexFileName
	}
	post, status, err := notes.getNote(post_md)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := error_template.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
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
