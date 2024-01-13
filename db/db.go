package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

func DB() (*sql.DB, error) {
	port, _ := strconv.Atoi(os.Getenv("DBPORT"))
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"), port, os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
