package main

import (

	"fmt"

	"net/http"
	"pos-rs/pkg/pos/model"
	"strings"
	"golang.org/x/time/rate"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

		if err := recover(); err != nil {
			w.Header().Set("Connection", "close")
			app.respondWithError(w, http.StatusInternalServerError, "Recover Panic") 
		}
		}()
			next.ServeHTTP(w, r)
	})
}


func (app *Application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
	app.rateLimitExceededResponse(w, r)
	return
	}
	next.ServeHTTP(w, r)
	})
	}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	w.Header().Add("Vary", "Authorization")

	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		r = app.contextSetUser(r, model.AnonymousEmployee)
		next.ServeHTTP(w, r)
	return
	}

	headerParts := strings.Split(authorizationHeader, " ")
	fmt.Println(headerParts)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		app.invalidAuthenticationTokenResponse(w, r)
	return
	}

	token := headerParts[1]


	user, err := app.Models.Employee.GetForToken(model.ScopeAuthentication, token)
	if err != nil {

			app.invalidAuthenticationTokenResponse(w, r)
	return
	}

	r = app.contextSetUser(r, user)

	next.ServeHTTP(w, r)
	})
}


func (app *Application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
		return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

func (app *Application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		user := app.contextGetUser(r)

		permissions, err := app.Models.Permissions.GetAllForUser(user.Id)
		if err != nil {
			app.respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(code)
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		

		next.ServeHTTP(w, r)
		}

	return app.requireActivatedUser(fn)
}
