package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"user-mgt-system/pkg/handlers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var tmpl *template.Template
var Store = sessions.NewCookieStore([]byte("usermanagementsecret"))

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")

	//set up sessions
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		HttpOnly: true,
	}
}

func initDB() {
	var err error

	db, err = sql.Open("mysql", "root:password@(127.0.0.1:3306)/usermanagement?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// check db connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	gRouter := mux.NewRouter()

	initDB()
	defer db.Close()

	gRouter.HandleFunc("/", handlers.Homepage(db, tmpl, Store)).
		Methods("GET")

	gRouter.HandleFunc("/register", handlers.RegisterPage(db, tmpl)).
		Methods("GET")

	gRouter.HandleFunc("/register", handlers.RegisterHandler(db, tmpl)).
		Methods("POST")

	gRouter.HandleFunc("/login", handlers.LoginPage(db, tmpl)).
		Methods("GET")

	gRouter.HandleFunc("/login", handlers.LoginHandler(db, tmpl, Store)).
		Methods("POST")

	http.ListenAndServe(":4000", gRouter)
}
