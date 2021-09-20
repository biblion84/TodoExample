package main

import (
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func (app *application) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "index.page.gohtml", TemplateData{
			Todos: []Todo{{
				Checked: false,
				Text:    "Buy milk",
			}, {
				Checked: true,
				Text:    "Buy eggs",
			}},
		})
	})
}
