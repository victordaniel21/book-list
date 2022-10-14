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

	// 2. Populate data to struct Book
	for rows.Next() {
		errScan := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if errScan != nil {
			return []models.Book{}, err
		}
	}
	return books, nil
}

func (b BookRepository) GetBook(
	db *sql.DB,
	book models.Book,
	id string,

) (models.Book, error) {
	// 1. query to fetch data in DB
	row := db.QueryRow("SELECT * FROM book.books WHERE id=$1", id)

	// 2. populate data to struct book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		return book, err
	}

	return book, nil
}

func (b BookRepository) AddBook(db *sql.DB, book models.Book) (int, error) {
	var bookID int

	//  command query to insert data
	err := db.QueryRow("INSERT INTO boook.books(title, author, year) VALUES($1, $2, $3) RETURNING id", book.Title, book.Author, book.Year).Scan(&bookID)

	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func (b BookRepository) DeleteBook(db *sql.DB, id string) (int64, error) {
	//  query to delete data in db
	result, err := db.Exec("DELETE FROM book.books WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	// checking how many rows was updated
	rowDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowDeleted, nil
}

func (b BookRepository) UpdateBook(db *sql.DB, book models.Book) (int64, error) {
	// query to fetch data in DB
	result, err := db.Exec("UPDATE book.books SET title=$1, author=$2, year=$3 WHERE id=$4 id", &book.Title, &book.Author, &book.Year, &book.ID)
	if err != nil {
		return 0, err
	}

	// checking how many rows was updated => result.RowsAffected()
	rowUpdated, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	// return data
	return rowUpdated, nil
}
