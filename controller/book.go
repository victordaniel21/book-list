package controller

import (
	"book-list/models"
	bookRep "book-list/repository/book"
	"book-list/utils"
	"database/sql"
	"net/http"
)

type controller struct{}

var books []models.Book

func (c controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Initiallize
		var (
			book   models.Book
			errMsg models.Error
		)

		// 2. Manage database
		books = []models.Book{}
		booksRepo := bookRep.BookRepository{}
		books, err := booksRepo.GetBooks(db, book, books)
		if err != nil {
			errMsg.Message = "failed fetching data books"
			utils.SendError(w, http.StatusInternalServerError, errMsg)
			return
		}

		// 3. send response data
		w.Header().Set("content-type", "application/json")
		utils.SendSuccess(w, books)
	}
}
