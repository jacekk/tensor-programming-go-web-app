package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

const HTTP_PORT = 8000
const TPLS_DIR = "templates/"
const ROUTE_TEST = "/test/"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	f := fmt.Sprintf("%s%s.txt", TPLS_DIR, p.Title)
	return ioutil.WriteFile(f, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	f := fmt.Sprintf("%s%s.txt", TPLS_DIR, title)
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func view(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len(ROUTE_TEST):]
	fmt.Printf("Serving `%s` page...\n", title)
	page, err := loadPage(title)

	if err != nil {
		fmt.Fprintf(res, "Page `%s` not found :(", title)
		return
	}

	tpl, _ := template.ParseFiles(fmt.Sprintf("%stpl.html", TPLS_DIR))
	tpl.Execute(res, page)
}

func main() {
	// p := &Page{Title: "edit", Body: []byte("Welcome to the Edit Page :)")}
	// p.save()
	http.HandleFunc(ROUTE_TEST, view)
	fmt.Println(fmt.Sprintf("Listening on http://localhost:%d", HTTP_PORT))
	http.ListenAndServe(fmt.Sprintf(":%d", HTTP_PORT), nil)
}
