package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Result  string
	Message string
	Status  string
}

var status = map[string]string{
	"pass":         "pass",
	"fail":         "fail",
	"ok":           "pass",
	"ko":           "fail",
	"pending":      "pending",
	"wait":         "pending",
	"inconclusive": "inconclusive",
}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9_.-]*)$")

func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("static/%s", r.URL.Path[1:]))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	renderStaticTemplate(w, r.URL.Path[1:])
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	title := m[1]
	if title == "" {
		renderStaticTemplate(w, "home.html")
		return
	}
	msg := r.URL.Query().Get("msg")
	p := &Page{Result: strings.ToUpper(title), Message: msg, Status: status[title]}
	renderDynamicTemplate(w, p)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	referer := strings.Split(r.Header.Get("Referer"), "?")
	referer = strings.Split(referer[0], "/")
	icon := strings.ToLower(referer[len(referer)-1])
	if validPath.MatchString("/" + icon) {
		icon := status[icon]
		if icon == "" {
			icon = "default"
		}
		http.ServeFile(w, r, fmt.Sprintf("favicon/%s.ico", icon))
		return
	}
	http.Error(w, "Invalid Referer", http.StatusPreconditionFailed)
}

func renderDynamicTemplate(w http.ResponseWriter, p *Page) {
	templates := template.Must(template.ParseFiles("templates/base.html", "templates/dynamic.html"))
	err := templates.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderStaticTemplate(w http.ResponseWriter, title string) {
	templates := template.Must(template.ParseFiles("templates/base.html", "templates/"+title))
	err := templates.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", dynamicHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)

	http.HandleFunc("/legal.html", staticHandler)
	http.HandleFunc("/home.html", staticHandler)
	http.HandleFunc("/style.css", fileHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
