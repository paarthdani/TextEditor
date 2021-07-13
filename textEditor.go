package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Editor struct {
	Title string
	Text  []byte
}

func loadPage(title string) (*Editor, error) {
	filename := title + ".txt"
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Editor{Title: title, Text: text}, nil
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, template string, editor *Editor) {
	err := templates.ExecuteTemplate(w, template+".html", editor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a View Page.")
}


func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	editor, err := loadPage(title)
	if err != nil {
		editor = &Editor{Title: title}
	}
	renderTemplate(w, "edit", editor)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a Save Page.")
}

func makeHandler(function func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		splitString := strings.Split(r.URL.Path, "/")
		if splitString == nil {
			http.NotFound(w, r)
			return
		}
		function(w, r, splitString[2])
	}
}

func main() {
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save", saveHandler)
	http.ListenAndServe(":8040", nil)
}
