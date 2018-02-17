package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

func main() {
	http.HandleFunc("/", dataHandler)
	log.Println("Open in browser: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	fileRead, _ := ioutil.ReadFile("../data/index.md")
	lines := strings.Split(string(fileRead), "\n")
	title := string(lines[0])
	body := strings.Join(lines[1:len(lines)], "\n")
	body = string(blackfriday.MarckdownCommon([]byte(body)))
	post := Post{title, template.HTML(body)}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
