package user

import (
	"book-list/models"
	"database/sql"
)

type UserRepository struct{}

func (u UserRepository) Signup(
	db *sql.DB,
	user models.User,
) (models.User, error) {
	statement := "insert into book.user (email, password) values($1, $2) returning id"

	err := db.QueryRow(statement, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return user, err
	}

	user.Password = ""
	return user, nil

}

func (u UserRepository) Signin(
	db *sql.DB,
	user models.User,
) (models.User, error) {
	row := db.QueryRow("select * from book.user where email=$1", user.Email)
	err := row.Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
