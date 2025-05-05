package main
import (
	"context"
	"net/http"
	"pos-rs/pkg/pos/model"
)

type contextKey string

const userContextKey = contextKey("employee")

func (app *Application) contextSetUser(r *http.Request, user *model.Employee) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) *model.Employee {
	user, ok := r.Context().Value(userContextKey).(*model.Employee)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}