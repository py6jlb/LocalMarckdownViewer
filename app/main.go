package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

var (
	post_template = template.Must(template.ParseFiles(path.Join("app", "templates", "layout.html"), path.Join("app", "templates", "post.html")))
)

func main() {
	fs := http.FileServer(http.Dir("app/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", dataHandler)
	log.Println("Open in browser: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// fileRead, _ := ioutil.ReadFile("../data/index.md")
	// lines := strings.Split(string(fileRead), "\n")
	// title := string(lines[0])
	// body := strings.Join(lines[1:len(lines)], "\n")
	// body = string(blackfriday.MarckdownCommon([]byte(body)))
	// post := Post{title, template.HTML(body)}
	if err := post_template.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
