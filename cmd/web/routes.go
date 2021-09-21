package main

import (
	"github.com/bmizerany/pat"
	"github.com/lpar/gzipped/v2"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := pat.New()
	router.Get("/", app.loginRequired(app.index()))
	router.Post("/", app.loginRequired(app.addPost()))
	router.Post("/checkTodo", app.loginRequired(app.checkTodoPost()))
	router.Get("/signup", app.signup())
	router.Post("/signup", app.signupPost())
	router.Get("/signin", app.signin())
	router.Post("/signin", app.signinPost())

	fileServer := gzipped.FileServer(gzipped.Dir("./ui/static"))
	router.Get("/static/", (http.StripPrefix("/static/", fileServer)))

	//router.NotFound = loginRequired.Then(app.notFound())

	return router
}
