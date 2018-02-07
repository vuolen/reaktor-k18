package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vuolen/reaktor-k18/db"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cdb, err := db.CachedDatabase()
	if err != nil {
		log.Panicf("Error while opening database: %s\n", err)
	}
	defer cdb.Close()

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json, err := cdb.GetJson()
			log.Print(cdb.Cache)
			if err != nil {
				log.Panicf("Error getting database: %s\n", err)
			}
			w.Write(json)
		} else if r.Method == "POST" {

		}
	})
	http.Handle("/", http.FileServer(http.Dir("./public-html")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
