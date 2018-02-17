package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type note struct {
	Title   string
	Body    template.HTML
	ModTime int64
}

type notesCollection struct {
	Notes map[string]note
	sync.RWMutex
}

func newNotesCollection() *notesCollection {
	n := notesCollection{}
	n.Notes = make(map[string]note)
	return &n
}

func (n *notesCollection) getNote(md string) (note, int, error) {
	info, err := os.Stat(md)
	if err != nil {
		if os.IsNotExist(err) {
			return note{}, http.StatusNotFound, err
		}
	}
	if info.IsDir() {
		return note{}, http.StatusNotFound, fmt.Errorf("dir")
	}
	val, ok := n.Notes[md]
	if !ok || (ok && val.ModTime != info.ModTime().UnixNano()) {
		n.RLock()
		defer n.RUnlock()
		fileread, _ := ioutil.ReadFile(md)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		body := strings.Join(lines[1:len(lines)], "\n")
		body = string(blackfriday.Run([]byte(body)))
		post := note{title, template.HTML(body), info.ModTime().UnixNano()}
		n.Notes[md] = post
	}
	result := n.Notes[md]
	return result, 200, nil
}
