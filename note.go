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
	Content template.HTML
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
		markdown := strings.Replace(string(fileread), "\r\n", "\n", -1)
		body := string(blackfriday.Run([]byte(markdown)))
		post := note{template.HTML(body), info.ModTime().UnixNano()}
		n.Notes[md] = post
	}
	result := n.Notes[md]
	return result, 200, nil
}
