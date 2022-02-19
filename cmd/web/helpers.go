package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

// TODO put user value in context.Context
func (app *application) getUser(r *http.Request) User {
	cookie, err := r.Cookie("session")
	if err != nil || cookie == nil {
		return User{}
	}

	session, err := app.SessionFind(cookie.Value)

	if err != nil || session.ID == 0 {
		return User{}
	}

	user, err := app.UserFindByEmail(session.Email)

	if err != nil || user.ID == 0 {
		return user
	}

	return user
}

func generateCookie() string {
	b := make([]byte, 60)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

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
