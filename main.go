package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vuolen/reaktor-k18/db"
	"github.com/vuolen/reaktor-k18/server"
)

func main() {
	tdb := db.TemperatureDatabase{sqlx.MustConnect("sqlite3", "./database.db")}
	defer tdb.Close()
	err := tdb.CreateTables()
	if err != nil {
		log.Printf("%+v", err)
	}
	err = tdb.PopulateWithDefaults()
	if err != nil {
		log.Printf("%+v", err)
	}

	ctx := &server.Context{tdb}

	r := mux.NewRouter()
	r.Methods("GET").Path("/locations").Name("GetLocations").Handler(server.HandlerWithContext{ctx, server.GetLocations})
	r.Methods("GET").Path("/logs").Name("GetLogs").Handler(server.HandlerWithContext{ctx, server.GetLogs})
	r.Methods("GET").Path("/logs/{locationId}").Name("GetLogsByLocationId").Handler(server.HandlerWithContext{ctx, server.GetLogsByLocationId})
	r.Methods("POST").Path("/logs/add").Name("AddLog").Handler(server.HandlerWithContext{ctx, server.AddLog})
	r.Methods("GET").PathPrefix("/").Name("FileServer").Handler(http.FileServer(http.Dir("./public-html")))
	log.Printf("%+v", http.ListenAndServe(":8080", r))
}
