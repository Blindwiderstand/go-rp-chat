package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"google.golang.org/appengine"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"fmt"
	"os"
	"bufio"
	"google.golang.org/appengine/file"
)

//Type template from Package Template
var tpl *template.Template
var db *sql.DB

//pageData type with underlying type struct
//Title and title are different. title would be unexported and could not be used in a template
//Title is exported due to capalization of the first letter.

type pageData struct {
	Title string
	FirstName string
	CharacterName string
	UserID int
}

func init() {
	var err error
	var password string

	if file, err := os.Open("mysql-google.config"); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			password = scanner.Text()
			log.Println(password)
		}

	} else {
		log.Fatal(err)
	}

	//Opens Connection to database. Needs a database driver for the right database.
	db, err = sql.Open("mysql",
		"go-admin:" + password + "@vivid-cargo-180511:europe-west1:character-db/character_db?tls=true")
	if err != nil {
		panic(err)
	}
	// Make sure it's connected
	if err = db.Ping(); err != nil {
		panic(err)
	}

	//needs relative reference
	tpl = template.Must(template.ParseGlob("gohtml-templates/*.gohtml"))
	http.HandleFunc("/", idx)
	//http.HandleFunc("/index", idx)
	http.HandleFunc("/about", about)
	http.HandleFunc("/apply", apply)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/create", create)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/favicon.ico", favicon)

	// Handle css and js
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./go-pub"))))
}

func main() {
	appengine.Main()
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

func create(w http.ResponseWriter, req *http.Request) {
	var character_name string

	pd := pageData{
		Title: "Create Page",
		UserID: 1,
	}

	if req.Method == http.MethodPost {
		character_name = req.FormValue("character_name")
		pd.CharacterName = character_name
	}

	db.Exec("INSERT INTO character_db(user_id, character_name) VALUES" +
		" ($1, $2)", pd.UserID, pd.CharacterName)

	tmplErr := tpl.ExecuteTemplate(w, "apply.gohtml", pd)
	if tmplErr != nil {
		log.Println(tmplErr)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}


}

func favicon (w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./favicon.ico")
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

func redirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/contact", http.StatusSeeOther)
}