package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = writeJSONError(w, http.StatusNotFound, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("duplicate keys found", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = writeJSONError(w, http.StatusConflict, "not found")
}
