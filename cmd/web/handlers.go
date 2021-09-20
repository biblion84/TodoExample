package main

import (
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func (app *application) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		app.db.Find(&todos)
		app.render(w, r, "index.page.gohtml", TemplateData{
			Todos: todos,
		})
	})
}

func (app *application) addPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		text := r.Form.Get("text")
		todo := Todo{
			Checked: false,
			Text:    text,
		}
		app.db.Create(&todo)

		app.index().ServeHTTP(w, r)
	})
}
