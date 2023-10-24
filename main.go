package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = make(map[string]*template.Template)

func main() {
	templates["fileviewer"] = loadTemplate("fileviewer")
	mux := http.NewServeMux()
	fileServer := FileServer2(Dir("./public"))
	mux.Handle("/contents/", http.StripPrefix("/contents/", fileServer))
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()

}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles("template/" + name + ".html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}
