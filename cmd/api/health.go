package main

import (
	"log"
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		Checks the health of the server
//	@Description	Checks the health of the server
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	""
//	@Failure		404	{object}	error	"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/health/ [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Println("Failed writing to JSON error")
		}
	}

}
