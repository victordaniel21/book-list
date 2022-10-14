package main

import (
	"book-list/controller"
	"book-list/drivers"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var dbPsql *sql.DB

func main() {
	dbPsql = drivers.CreateConnection()
	c := controller.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/books", c.GetBooks(dbPsql)).Methods("GET")
	router.HandleFunc("/books/{id}", c.GetBook(dbPsql)).Methods("GET")
	router.HandleFunc("/books", c.AddBook(dbPsql)).Methods("POST")
	router.HandleFunc("/books", c.UpdateBook(dbPsql)).Methods("PUT")
	router.HandleFunc("/books/{id}", c.DeleteBook(dbPsql)).Methods("DELETE")

	router.HandleFunc("/signup", c.Signup(dbPsql)).Methods("POST")
	router.HandleFunc("/signin", c.Signin(dbPsql)).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))

}
