package server

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vuolen/reaktor-k18/db"
)

type server struct {
	handler Handler
}

func New() (server, error) {
	tdb, err := db.OpenDatabase("./database.db")
	if err != nil {
		return server{}, err
	}
	handler := Handler{tdb}
	return server{handler}, nil
}

func (s server) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/locations", s.handler.GetLocations)
	r.HandleFunc("/logs", s.handler.GetLogs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public-html")))
	//r.HandleFunc("/logs/add", s.handler.AddLog)
	return http.ListenAndServe(":8080", r)
}

func (s server) Close() error {
	return s.handler.tdb.Close()
}
