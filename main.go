package main

import (
	"html/template"
	"log"
	"net/http"
	"google.golang.org/appengine"
)

//Type template from Package Template
var tpl *template.Template

//pageData type with underlying type struct
//Title and title are different. title would be unexported and could not be used in a template
//Title is exported due to capalization of the first letter.

type pageData struct {
	Title string
	FirstName string
}

func init() {
	//needs relative reference
	tpl = template.Must(template.ParseGlob("gohtml-templates/*.gohtml"))
}

func main() {
	appengine.Main()
	http.HandleFunc("/", idx)
	//http.HandleFunc("/index", idx)
	http.HandleFunc("/about", about)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/apply", apply)
	http.HandleFunc("/redirect", redirect)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	//http.ListenAndServe(":8080", nil)
}

func idx(w http.ResponseWriter, req *http.Request) {

	pd := pageData{
		Title: "Index Page",
	}

	//Denies any other requests except GET
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	err := tpl.ExecuteTemplate(w, "index.gohtml", pd)
	if err != nil {
		//Println is Printline
		log.Println(err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
}

func about(w http.ResponseWriter, req *http.Request) {

	pd := pageData{
		Title: "About Page",
	}

	err := tpl.ExecuteTemplate(w, "about.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
}

func contact(w http.ResponseWriter, req *http.Request) {

	pd := pageData{
		Title: "Contact Page",
	}

	err := tpl.ExecuteTemplate(w, "contact.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
}

func apply(w http.ResponseWriter, req *http.Request) {
	var first string

	pd := pageData{
		Title: "Apply Page",
	}

	// Wenn die HTTP-Methode POST anstatt GET ist dann führe folgendes aus
	//req.Method ist der http-Request der eine Konstante hat
	if req.Method == http.MethodPost {
		//single equal sign because we are not initializing the variable
		first = req.FormValue("fname")
		pd.FirstName = first
	}

	err := tpl.ExecuteTemplate(w, "apply.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/contact", http.StatusSeeOther)
}