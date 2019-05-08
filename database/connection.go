package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DSN = ""

func Initialize() {
	log.Print("Initializing Database...")
	var DB_USER = os.Getenv("SETSISAW_DB_USER")
	var DB_PASS = os.Getenv("SETSISAW_DB_PASS")
	var DB_HOST = os.Getenv("SETSISAW_DB_HOST")
	var DB_NAME = os.Getenv("SETSISAW_DB_NAME")

	error := false
	if DB_USER == "" {
		log.Printf("SETSISAW_DB_USER environment variable not set")
		error = true
	}
	if DB_PASS == "" {
		log.Printf("SETSISAW_DB_PASS environment variable not set")
		error = true
	}
	if DB_NAME == "" {
		log.Printf("SETSISAW_DB_NAME environment variable not set")
		error = true
	}
	if DB_HOST == "" {
		log.Printf("SETSISAW_DB_HOST environment variable not set")
		error = true
	}

	if error {
		log.Fatal("Missing Database login information. Not Starting.")
	}

	DSN = DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/" + DB_NAME

	connection, err := GetConnection()
	if err != nil {
		log.Fatal("Failed to initiate database connection. Not Starting.")
	}

	_ = connection.Close()
	log.Print("done")
}

func GetConnection() (*sql.DB, error) {
	// Open DB connection
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		return nil, err
	}

	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
