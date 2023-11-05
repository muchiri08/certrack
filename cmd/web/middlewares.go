package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/muchiri08/certrack/internal/models"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.sessionManager.Exists(r.Context(), "username")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		username := app.sessionManager.Get(r.Context(), "username").(string)
		user, err := app.models.Users.GetUserByUsername(username)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrNoRecordFound):
				app.sessionManager.Remove(r.Context(), "username")
				return
			default:
				app.serverError(w, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requiredAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
