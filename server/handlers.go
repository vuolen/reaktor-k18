package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/vuolen/reaktor-k18/db"
)

type Context struct {
	Tdb db.TemperatureDatabase
}

type HandlerWithContext struct {
	Ctx     *Context
	Handler func(ctx *Context, w http.ResponseWriter, r *http.Request) error
}

func writeError(w http.ResponseWriter, code int, message string) {
	if message == "" {
		message = http.StatusText(code)
	}
	w.WriteHeader(code)
	var data struct {
		Error string `json:error`
	}
	data.Error = message
	json.NewEncoder(w).Encode(data)
}

func (hwc HandlerWithContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	err := hwc.Handler(hwc.Ctx, w, r)
	log.Printf(
		"%s\t%s\t%s\t%fms",
		r.Method,
		r.RequestURI,
		mux.CurrentRoute(r).GetName(),
		time.Since(start).Seconds()*1000,
	)
	if err != nil {
		log.Printf("%+v", err)
	}
}

func GetLocations(ctx *Context, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	locations := make([]db.Location, 0)
	err := ctx.Tdb.Select(&locations, "select * from locations")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return errors.WithStack(err)
	}
	json.NewEncoder(w).Encode(locations)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return errors.WithStack(err)
	}
	return nil
}

func GetLogs(ctx *Context, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logs := make([]db.TemperatureLog, 0)
	err := ctx.Tdb.Select(&logs, "select * from logs")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return errors.WithStack(err)
	}
	json.NewEncoder(w).Encode(logs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return errors.WithStack(err)
	}
	return nil
}

func AddLog(ctx *Context, w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return errors.WithStack(err)
	}
	var tlog db.TemperatureLog
	err = json.Unmarshal(b, &tlog)
	if err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			writeError(w, http.StatusBadRequest, "JSON syntax error")
		case *json.UnmarshalTypeError:
			writeError(w, http.StatusBadRequest, "Invalid JSON object type")
		default:
			writeError(w, http.StatusBadRequest, "")
		}
		return errors.WithStack(err)
	}
	if tlog.Time < 0 || tlog.Time > time.Now().UTC().Unix() {
		writeError(w, http.StatusBadRequest, "Time out of bounds")
		return errors.New("Time out of bounds")
	}
	// 373.15 kelvin equals 100 degrees Celcius
	if tlog.Temperature < 0 || tlog.Temperature > 373.15 {
		writeError(w, http.StatusBadRequest, "Temperature out of bounds")
		return errors.New("Temperature out of bounds")
	}
	locations := make([]db.Location, 0)
	err = ctx.Tdb.Select(&locations, "select * from locations")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return errors.WithStack(err)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return err
	}
	validLocation := false
	for _, loc := range locations {
		if loc.Id == tlog.LocationId {
			validLocation = true
		}
	}
	if !validLocation {
		writeError(w, http.StatusBadRequest, "Invalid location")
		return errors.New("Invalid location")
	}
	_, err = ctx.Tdb.Exec("insert into logs(locationId, time, temperature) values (?, ?, ?)", tlog.LocationId, tlog.Time, tlog.Temperature)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "")
		return err
	}
	return nil
}
