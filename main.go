package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates = make(map[string]*template.Template)

func main() {
	templates["fileviewer"] = loadTemplate("fileviewer")
	mux := http.NewServeMux()
	fileServer := FileServer2(Dir("./public"))
	mux.Handle("/contents/", http.StripPrefix("/contents/", fileServer))
	uploadHandler := &UploadHandler{}
	mux.Handle("/upload/", http.StripPrefix("/upload/", uploadHandler))
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

type UploadHandler struct {
}

func (u *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, fmt.Sprintf("FormFile Error: %s", err), 500)
		return
	}
	defer file.Close()
	basename := filepath.Base(header.Filename)
	path := filepath.Join("public", r.URL.Path, basename)
	if f, err := os.Stat(path); !os.IsNotExist(err) && f.IsDir() {
		http.Error(w, fmt.Sprintf("Directory Exist: %s", err), 500)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, fmt.Sprintf("FileOpenError %s", err), 500)
	}
	defer file.Close()
	io.Copy(f, file)

	//w.Header().Set("Location", "/contents/"+r.URL.Path)
	//w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Fprintf(w, "Upload Success")
}
