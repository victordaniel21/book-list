package drivers

import (
	"book-list/goconf"
	"database/sql"
	"fmt"
	"log"
)

var (
	host     = goconf.Config().GetString("postgres.host")
	port     = goconf.Config().GetInt("postgres.port")
	user     = goconf.Config().GetString("postgres.user")
	password = goconf.Config().GetString("postgres.password")
	dbname   = goconf.Config().GetString("postgres.db")
)

var db *sql.DB

func CreateConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
