package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Editor struct {
	Title string
	Text  []byte
}

type Editors struct {
	FileName []string
}

func (editor *Editor) save() error {
	dir, _ := os.Getwd()
	os.Chdir(dir + "/files")
	filename := editor.Title + ".txt"
	err := ioutil.WriteFile(filename, editor.Text, 0600)
	os.Chdir(dir)
	return err
}

func loadPage(title string) (*Editor, error) {
	filename := title + ".txt"
	dir, _ := os.Getwd()
	os.Chdir(dir + "/files")
	text, err := ioutil.ReadFile(filename)
	os.Chdir(dir)
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	editor, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", editor)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	editor, err := loadPage(title)
	if err != nil {
		editor = &Editor{Title: title}
	}
	renderTemplate(w, "edit", editor)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	editor := &Editor{Title: title, Text: []byte(body)}
	err := editor.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("files/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	filename := make([]string, 0)
	for _, file := range files {
		filename = append(filename, strings.Split(file.Name(), ".")[0])
	}
	fmt.Println(filename)
	editors := &Editors{FileName: filename}

	fmt.Println(editors.FileName)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
