package main

import (
	"net/http"
)

func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.respondWithError(w, http.StatusTooManyRequests, message)
}

func (app *Application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.respondWithError(w, http.StatusUnauthorized, message)
}

func (app *Application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.respondWithError(w, http.StatusUnauthorized, message)
	}

func (app *Application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.respondWithError(w, http.StatusForbidden, message)
}

func (app *Application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.respondWithError(w, http.StatusForbidden, message)
}
