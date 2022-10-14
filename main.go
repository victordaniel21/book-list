package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books, Book{ID: 1, Title: "Golang is Fun", Author: "John Doe", Year: "2010"},
		Book{ID: 2, Title: "Golang is Easy", Author: "Gopher", Year: "2011"},
		Book{ID: 3, Title: "Golang is Multipurpose Programming Language", Author: "John Walker", Year: "2011"},
		Book{ID: 4, Title: "Golang is Fast", Author: "Jimmy Doe", Year: "2012"},
		Book{ID: 5, Title: "Golang is Cool", Author: "Nick Name", Year: "2015"})

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))

}
