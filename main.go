package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

const HTTP_PORT = 8000

const ROUTE_SHOW = "/show/"
const ROUTE_EDIT = "/edit/"
const ROUTE_SAVE = "/save/"

type Page struct {
	Title string
	Body  []byte
}

func getTplPath(name string) string {
	return fmt.Sprintf("templates/%s.html", name)
}

func getPagePath(name string) string {
	return fmt.Sprintf("pages/%s.txt", name)
}

func (p *Page) save() error {
	f := getPagePath(p.Title)
	return ioutil.WriteFile(f, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	f := getPagePath(title)
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func viewPage(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len(ROUTE_SHOW):]
	fmt.Printf("Showing `%s` page...\n", title)
	page, err := loadPage(title)

	if err != nil {
		fmt.Fprintf(res, "Page `%s` NOT FOUND :(", title)
		return
	}

	tpl, _ := template.ParseFiles(getTplPath("show"))
	tpl.Execute(res, page)
}

func editPage(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len(ROUTE_EDIT):]
	fmt.Printf("Showing `%s` edit page...\n", title)
	page, err := loadPage(title)

	if err != nil {
		fmt.Fprintf(res, "Page `%s` NOT FOUND :(", title)
		return
	}

	tpl, _ := template.ParseFiles(getTplPath("edit"))
	tpl.Execute(res, page)
}

func savePage(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len(ROUTE_SAVE):]
	body := req.FormValue("body")
	fmt.Printf("Saving `%s` page...\n", title)
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	redirectTo := fmt.Sprintf("%s%s", ROUTE_SHOW, title)
	http.Redirect(res, req, redirectTo, http.StatusFound)
}

func main() {
	http.HandleFunc(ROUTE_SHOW, viewPage)
	http.HandleFunc(ROUTE_EDIT, editPage)
	http.HandleFunc(ROUTE_SAVE, savePage)
	fmt.Println(fmt.Sprintf("Listening on http://localhost:%d", HTTP_PORT))
	http.ListenAndServe(fmt.Sprintf(":%d", HTTP_PORT), nil)
}
