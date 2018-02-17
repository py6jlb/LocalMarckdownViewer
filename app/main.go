package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bmizerany/pat"
	"gopkg.in/russross/blackfriday.v2"
)

var (
	post_template  = template.Must(template.ParseFiles(path.Join("app", "templates", "layout.html"), path.Join("app", "templates", "post.html")))
	error_template = template.Must(template.ParseFiles(path.Join("app", "templates", "layout.html"), path.Join("app", "templates", "error.html")))
)

type Post struct {
	Title string
	Body  template.HTML
}

func main() {
	fs := http.FileServer(http.Dir("app/assets"))

	mux := pat.New()
	mux.Get("/:page", http.HandlerFunc(dataHandler))
	mux.Get("/:page/", http.HandlerFunc(dataHandler))
	mux.Get("/", http.HandlerFunc(dataHandler))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", mux)
	log.Println("Open in browser: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	page := params.Get(":page")
	p := path.Join("data", page)
	var post_md string
	if page != "" {
		post_md = p + ".md"
	} else {
		post_md = p + "/index.md"
	}
	post, status, err := load_post(post_md)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func load_post(md string) (Post, int, error) {
	info, err := os.Stat(md)
	if err != nil {
		if os.IsNotExist(err) {
			return Post{}, http.StatusNotFound, err
		}
	}
	if info.IsDir() {
		return Post{}, http.StatusNotFound, fmt.Errorf("dir")
	}
	fileread, _ := ioutil.ReadFile(md)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	body := strings.Join(lines[1:len(lines)], "\n")
	body = string(blackfriday.Run([]byte(body)))
	post := Post{title, template.HTML(body)}
	return post, 200, nil
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := error_template.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
