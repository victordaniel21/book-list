package main

import (
	"book-list/goconf"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var (
	host     = goconf.Config().GetString("postgres.host")
	port     = goconf.Config().GetInt("postgres.port")
	user     = goconf.Config().GetString("postgres.user")
	password = goconf.Config().GetString("postres.password")
	dbname   = goconf.Config().GetString("postgres.dbname")
)

func createConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {
	//1. catch query parameter
	params := mux.Vars(r)

	//2. convert the data from string to integer
	id, _ := strconv.Atoi(params["id"])

	//3. search the data in slice books, if match show the data
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
		}

	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	//1. Decode body json
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	//2. Adding the data to slice books
	books = append(books, book)

	//3. Encode the data books
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {

	//1. decode body json
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	//2. searching data in slice books. if match, put the data
	for idx, item := range books {
		if item.ID == book.ID {
			books[idx] = book
		}
	}

	//3. Encode the updated data
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	//1. catch query parameter
	params := mux.Vars(r)

	//2. convert data string to int
	id, _ := strconv.Atoi(params["id"])

	//3. searching data in slice books. if match, delete the data
	for idx, book := range books {
		if book.ID == id {
			books = append(books[:idx], books[idx+1:]...)
		}
	}

	//4. Encode the data books
	json.NewEncoder(w).Encode(books)
}

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

	log.Fatal(http.ListenAndServe(":8000", router))

}
