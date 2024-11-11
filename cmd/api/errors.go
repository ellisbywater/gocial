package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error %s path: %s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered an issue")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error %s path: %s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictingResourceError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflicting resource error %s path: %s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error %s path: %s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found")
}
