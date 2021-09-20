package main

import (
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func (app *application) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "index.page.gohtml", TemplateData{
			Todos: app.todos,
		})
	})
}

func (app *application) addPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		text := r.Form.Get("text")
		app.todos = append(app.todos, Todo{
			Checked: false,
			Text:    text,
		})

		app.index().ServeHTTP(w, r)
	})
}
