package main

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available", "system_info": map[string]string{
			"environment": app.Config.Env,
			"version":     version},
	}

	app.respondWithJSON(w, http.StatusOK, env)

}