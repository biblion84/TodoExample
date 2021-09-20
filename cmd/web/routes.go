package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/lpar/gzipped/v2"
	"net/http"
)

func (app *application) routes() http.Handler {

	//standardMiddleware := alice.New(gziphandler.GzipHandler)
	standardMiddleware := alice.New()

	router := pat.New()
	router.Get("/", app.index())
	router.Post("/add", app.addPost())

	fileServer := gzipped.FileServer(gzipped.Dir("./ui/static"))
	router.Get("/static/", (http.StripPrefix("/static/", fileServer)))

	//router.NotFound = loginRequired.Then(app.notFound())

	return standardMiddleware.Then(router)
}
