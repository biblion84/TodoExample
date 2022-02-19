package main

import "net/http"

func (app *application) loginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie == nil {
			app.signup().ServeHTTP(w, r)
			return
		}
		session, err := app.SessionFind(cookie.Value)
		if err != nil || session.ID == 0 {
			app.signup().ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
