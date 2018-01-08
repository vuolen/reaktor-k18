package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vuolen/reaktor-k18/db"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbHandle, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panicf("Error while opening database: %s\n", err)
	}
	defer dbHandle.Close()
	err = db.CreateDatabase(dbHandle)
	if err != nil {
		log.Fatalf("SQL error: %+v\n", err)
	}
}
