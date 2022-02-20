package main

import (
	"context"
	"net/http"
)

func (app *application) loginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie == nil {
			app.signup().ServeHTTP(w, r)
			return
		}
		session, err := app.SessionFind(r.Context(), cookie.Value)
		if err != nil || session.ID == 0 {
			app.signup().ServeHTTP(w, r)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))

		next.ServeHTTP(w, r)
	})
}
