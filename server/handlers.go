package server

import (
	"encoding/json"
	"fmt"
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
	Handler func(ctx *Context, w ApiResponseWriter, r *http.Request) error
}

func (hwc HandlerWithContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	err := hwc.Handler(hwc.Ctx, ApiResponseWriter{w}, r)
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

func GetLocations(ctx *Context, w ApiResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	locations := make([]db.Location, 0)
	err := ctx.Tdb.Select(&locations, "select * from locations")
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return errors.WithStack(err)
	}
	json.NewEncoder(w).Encode(locations)
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return errors.WithStack(err)
	}
	return nil
}

func GetLogs(ctx *Context, w ApiResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logs := make([]db.TemperatureLog, 0)
	err := ctx.Tdb.Select(&logs, "select * from logs")
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return errors.WithStack(err)
	}
	json.NewEncoder(w).Encode(logs)
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return errors.WithStack(err)
	}
	return nil
}

func AddLog(ctx *Context, w ApiResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return errors.WithStack(err)
	}

	var tlog db.TemperatureLog
	err = json.Unmarshal(b, &tlog)
	if err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			serr := err.(*json.SyntaxError)
			w.WriteErrorWithMessage(http.StatusBadRequest, fmt.Sprintf("JSON syntax error at offset %d", serr.Offset))
		case *json.UnmarshalTypeError:
			w.WriteErrorWithMessage(http.StatusBadRequest, "Invalid JSON object type")
		default:
			w.WriteDefaultError(http.StatusBadRequest)
		}
		return errors.WithStack(err)
	}

	if tlog.Time < 0 || tlog.Time > time.Now().UTC().Unix() {
		w.WriteErrorWithMessage(http.StatusBadRequest, "Invalid time")
		return errors.New("Invalid time")
	}

	// 373.15 kelvin equals 100 degrees Celcius
	if tlog.Temperature < 0 || tlog.Temperature > 373.15 {
		w.WriteErrorWithMessage(http.StatusBadRequest, "Invalid temperature")
		return errors.New("Invalid temperature")
	}

	isValidLocation, err := ctx.Tdb.IsValidLocationId(tlog.LocationId)
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return err
	} else if !isValidLocation {
		w.WriteErrorWithMessage(http.StatusBadRequest, "Invalid location")
		return errors.New("Invalid location")
	}

	err = ctx.Tdb.AddLog(tlog)
	if err != nil {
		w.WriteDefaultError(http.StatusInternalServerError)
		return err
	}
	w.WriteMessage("Log added")
	return nil
}
