//go:generate goversioninfo -icon=assets/icons/icon.ico
package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
)

var (
	postTemplate  = getTemplate("post.html")
	errorTemplate = getTemplate("error.html")
	notes         = newNotesCollection()
)

func getTemplate(templateName string) *template.Template {
	templateBox := rice.MustFindBox("assets")
	layoutTemplateString, err := templateBox.String("templates/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	templateString, err := templateBox.String("templates/" + templateName)
	if err != nil {
		log.Fatal(err)
	}
	template, err := template.New(templateName).Parse(layoutTemplateString + templateString)
	if err != nil {
		log.Fatal(err)
	}
	return template
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}
	assetsBox := rice.MustFindBox("assets").HTTPBox()
	mux := mux.NewRouter()
	staticHandler := http.StripPrefix("/static", noDirListing(http.FileServer(assetsBox)))
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
			errorHandler(w, r, 404)
			return
		}
		h.ServeHTTP(w, r)
	})
}
