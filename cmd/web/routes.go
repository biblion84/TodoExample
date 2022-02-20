package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	router := http.NewServeMux()
	router.Handle("/", app.loginRequired(app.index()))
	router.Handle("/checkTodo", app.loginRequired(app.checkTodoPost()))
	router.Handle("/signup", app.signup())
	router.Handle("/signin", app.signin())
	router.Handle("/test", app.test())

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handle("/static/", (http.StripPrefix("/static/", fileServer)))

	//router.NotFound = loginRequired.Then(app.notFound())

	return router
}
