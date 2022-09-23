package book

import (
	"book-list/models"
	"database/sql"
)

type BookRepository struct{}

func (b BookRepository) GetBooks(
	db *sql.DB,
	book models.Book,
	books []models.Book,

) ([]models.Book, error) {
	// 1. query to fetch data in DB
	rows, err := db.Query("SELECT *FROM book.books")
	if err != nil {
		return []models.Book{}, err
	}

	// 2. populate data to struct book
	for rows.Next() {
		errScan := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if errScan != nil {
			return []models.Book{}, err
		}
		books = append(books, book)
	}

	return books, nil
}
