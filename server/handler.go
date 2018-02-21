package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vuolen/reaktor-k18/db"
)

type Handler struct {
	tdb db.TemperatureDatabase
}

func (handler Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logs, err := handler.tdb.GetLogs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%+v", err)
		return
	}
	json.NewEncoder(w).Encode(logs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%+v", errors.WithStack(err))
		return
	}
}

func (handler Handler) GetLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	locations, err := handler.tdb.GetLocations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%+v", err)
		return
	}
	json.NewEncoder(w).Encode(locations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%+v", errors.WithStack(err))
		return
	}
}
