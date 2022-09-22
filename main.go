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
	// 1. Initiate Book and ID
	var (
		book   Book
		bookID int
	)

	// 2. Get data request, then decode to Book
	json.NewDecoder(r.Body).Decode(&book)

	// 3. Get connection
	dbConn := createConnection()
	defer dbConn.Close()

	// 4. Command query to insert data
	err := dbConn.QueryRow(`
		INSERT INTO book.books(title, author, year)
		VALUES($1, $2, $3)
		RETURNING id`, book.Title, book.Author, book.Year).Scan(&bookID)

	if err != nil {
		log.Fatal(err)
	}

	// 5. Return data(encode)
	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	// 0. Initiate struct Book
	var book Book

	// 1. Get data request, then decode to Book
	json.NewDecoder(r.Body).Decode(&book)

	// 2. Get connection
	dbConn := createConnection()
	defer dbConn.Close()

	// 3. Query to fetch data in DB
	result, err := dbConn.Exec("UPDATE book.books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Checking how many rows was updated => err.RowsAffected()
	rowUpdated, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	// 5. Response data
	json.NewEncoder(w).Encode(rowUpdated)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	//1. catch query parameter
	qParams := mux.Vars(r)

	//2. get connection
	dbConn := createConnection()
	defer dbConn.Close()

	//3. Query to delete data in DB
	result, err := dbConn.Exec("DELETE FROM boo.books WHERE id=$1", qParams["id"])
	if err != nil {
		log.Fatal(err)
	}
	//4. Checking how many rows was updated => result.RowsAffected()
	rowDeleted, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	//5. Response data
	json.NewEncoder(w).Encode(rowDeleted)
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
