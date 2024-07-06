package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body []byte
}

var templates = template.Must(template.ParseFiles("./templates/view.html", "./templates/edit.html"))

func renderTemplate(w http.ResponseWriter, pathToTemp string, p *Page) {
    err := templates.ExecuteTemplate(w, pathToTemp, p)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (p *Page) save() error {
	filename := p.Title + ".txt"

	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"

	body, err := os.ReadFile(filename)

	if(err != nil) {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		
		return
	}

	renderTemplate(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {	
	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)} 
    
	err := p.save()
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} 

    http.Redirect(w, r, "/view/" + title, http.StatusFound) 
}

var validPath = regexp.MustCompile("^/(edit|save|view|delete)/([a-zA-Z0-9\\-]+)$")

func makeTitleAndPathHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paths := validPath.FindStringSubmatch(r.URL.Path)

		if paths == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, paths[2])
	}
}

func main() {
	
	http.HandleFunc("/view/", makeTitleAndPathHandler(viewHandler))
	http.HandleFunc("/edit/", makeTitleAndPathHandler(editHandler))
	http.HandleFunc("/save/", makeTitleAndPathHandler(saveHandler))

	fmt.Println("Server listen at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}