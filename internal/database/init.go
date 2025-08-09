package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

type Data struct {
	Hash        string
	Name        string
	Category    string
	Description string
	Time        string
	Size        string
	Uploader    string
}

func InitDatabase() (*sql.DB, error) {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		log.Fatalln("[init.go]", err.Error())
	}

	return db, nil
}
