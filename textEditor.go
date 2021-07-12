package main

import (
	"fmt"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a View Page.")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is an Edit Page.")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a Save Page.")
}

func main() {
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/save", saveHandler)
	http.ListenAndServe(":8090", nil)
}
