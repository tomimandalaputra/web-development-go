package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	userRepo UserRepository
}

func main() {

	db, err := connectToDB("users_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: log.New(os.Stdout, "ERROR\t", log.Ltime|log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ltime|log.LstdFlags),
		userRepo: NewSQLUserRepository(db),
	}

	log.Print("Listening on :8080")
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

func connectToDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
