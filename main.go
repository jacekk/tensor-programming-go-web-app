package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

// HTTPPort ...
const HTTPPort = 8000

// RouteShow ...
const RouteShow = "/show/"

// RouteEdit ...
const RouteEdit = "/edit/"

// RouteSave ...
const RouteSave = "/save/"

// Page ...
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
	title := req.URL.Path[len(RouteShow):]
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
	title := req.URL.Path[len(RouteEdit):]
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
	title := req.URL.Path[len(RouteSave):]
	body := req.FormValue("body")
	fmt.Printf("Saving `%s` page...\n", title)
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	redirectTo := fmt.Sprintf("%s%s", RouteShow, title)
	http.Redirect(res, req, redirectTo, http.StatusFound)
}

func main() {
	http.HandleFunc(RouteShow, viewPage)
	http.HandleFunc(RouteEdit, editPage)
	http.HandleFunc(RouteSave, savePage)
	fmt.Println(fmt.Sprintf("Listening on http://localhost:%d", HTTPPort))
	http.ListenAndServe(fmt.Sprintf(":%d", HTTPPort), nil)
}
