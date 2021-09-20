package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, templateData TemplateData) {
	var ts *template.Template
	if true {
		app.templateCache, _ = newTemplateCache("./ui/html/")
	}

	var ok bool
	ts, ok = app.templateCache[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist", name), http.StatusInternalServerError)
		return
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then
	// return.

	err := ts.Execute(buf, &templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}
