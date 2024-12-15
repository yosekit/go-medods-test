package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"goauth/config"

	_ "github.com/lib/pq"
)

var (
	Client *sql.DB

	db = config.DB()
	err error
)

func init() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		"postgres", db.Port, db.User, db.Password, db.DBName)

	log.Println(connStr)

	Client, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
}