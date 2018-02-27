package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type ApiResponseWriter struct {
	http.ResponseWriter
}

type Message struct {
  Message string `json:"message"`
}

type Error struct {
  Error string `json:"error"`
}

func (w ApiResponseWriter) WriteJson(obj interface{}) {
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
	}
}

func (w ApiResponseWriter) WriteDefaultError(code int) {
	w.WriteHeader(code)
	w.WriteJson(Error{http.StatusText(code)})
}

func (w ApiResponseWriter) WriteErrorWithMessage(code int, message string) {
	w.WriteHeader(code)
	w.WriteJson(Error{message})
}

func (w ApiResponseWriter) WriteMessage(message string) {
	w.WriteHeader(200)
	w.WriteJson(struct {
		Message string `json:"message"`
	}{message})
}
