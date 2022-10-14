package controller

import (
	"book-list/models"
	bookRep "book-list/repository/book"
	"book-list/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type controller struct{}

var books []models.Book

func (c controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user   models.User
			errMsg models.Error
		)

		// decode request body to JSON to User struct
		json.NewDecoder(r.Body).Decode(&user)

		// validation field email should not be empty
		if user.Email == "" {
			errMsg.Message = "Email should not be empty"
			utils.SendError(w, http.StatusBadRequest, errMsg)
			return
		}

		// validation field password should not be empty
		if user.Password == "" {
			errMsg.Message = "Password should not be empty"
			utils.SendError(w, http.StatusBadRequest, errMsg)
			return
		}

		// Hashing password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Fatal(err)
		}

		// populated data hash to user.password
		user.Password = string(hash)

		// calling repository
		userRep := userRepo.UserRepository{}
		user, errUser := userRep.Signup(db, user)
		if errUser != nil {
			errMsg.Message = "Soemthing gone wrogn with sql" + errUser.Error()
			utils.SendError(w, http.StatusInternalServerError, errMsg)
			return
		}

		utils.SendSuccess(w, user)
	}
}

func (c controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. initialize; creating var
		var (
			book   models.Book
			errMsg models.Error
		)

		// 2. Manage DB
		books := []models.Book{}
		bookRepo := bookRep.BookRepository{}
		books, err := bookRepo.GetBooks(db, book, books)
		if err != nil {
			errMsg.Message = "Failed catching data books"
			utils.SendError(w, http.StatusInternalServerError, errMsg)
			return
		}

		// Response DB
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, books)

	}
}

func (c controller) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Initiallize struct book
		var (
			book   models.Book
			errMsg models.Error
		)
		// 2. Get query parameter
		params := mux.Vars(r)

		// 3. Manage database

		bookRepo := bookRep.BookRepository{}
		books, err := bookRepo.GetBook(db, book, params["id"])
		if err != nil {
			errMsg.Message = "failed fetching data books"
			utils.SendError(w, http.StatusInternalServerError, errMsg)

		}

		// 3. send response data

		utils.SendSuccess(w, books)
	}
}

func (c controller) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			book   models.Book
			bookID int
			errMsg models.Error
		)

		// 2. Get Data Request, then decode to book
		json.NewDecoder(r.Body).Decode(&book)

		// 3. Called repository
		bookRepo := bookRep.BookRepository{}
		bookID, err := bookRepo.AddBook(db, book)
		if err != nil {
			errMsg.Message = "Failed Adding Book"
			utils.SendError(w, http.StatusInternalServerError, errMsg)
		}

		// 4. Return data (encode)
		utils.SendSuccess(w, bookID)
	}
}

func (c controller) DeleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errMsg models.Error

		// get query parameter  {id}
		qparams := mux.Vars(r)

		// calling book repository
		bookRepo := bookRep.BookRepository{}
		rowsDeleted, err := bookRepo.DeleteBook(db, qparams["id"])
		if err != nil {
			errMsg.Message = "Failed deleting book"
			utils.SendError(w, http.StatusInternalServerError, errMsg)
		}

		// response data
		utils.SendSuccess(w, rowsDeleted)
	}
}

func (c controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// initiate struct Book
		var (
			book   models.Book
			errMsg models.Error
		)

		// Get data request than decode to Book
		json.NewDecoder(r.Body).Decode(&book)

		// calling bookRepository
		bookRepo := bookRep.BookRepository{}
		rowUpdated, err := bookRepo.UpdateBook(db, book)
		if err != nil {
			errMsg.Message = "Failed updating book"
			utils.SendError(w, http.StatusInternalServerError, errMsg)
		}

		// response data
		utils.SendSuccess(w, rowUpdated)
	}
}
