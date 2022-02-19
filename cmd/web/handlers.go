package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func (app *application) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			app.addPost().ServeHTTP(w, r)
			return
		}

		user := app.getUser(r)
		todos, err := app.TodoGetByUserId(user.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		app.render(w, r, "index.page.gohtml", TemplateData{
			Todos: todos,
		})
	})
}

func (app *application) addPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		text := r.Form.Get("text")
		user := app.getUser(r)
		todo := Todo{
			Checked: false,
			Text:    text,
			UserID:  user.ID,
		}
		err := app.TodoCreate(&todo)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (app *application) checkTodoPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.index().ServeHTTP(w, r)
			return
		}

		r.ParseForm()

		id, err := strconv.Atoi(r.Form.Get("id"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = app.TodoSetCheck(r.Form.Get("checked") == "true", id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})
}

func (app *application) signup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			app.signupPost().ServeHTTP(w, r)
			return
		}

		app.render(w, r, "signup.page.gohtml", TemplateData{})
	})
}

func (app *application) signupPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.index().ServeHTTP(w, r)
			return
		}

		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		userExist, err := app.UserExist(email)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if userExist {
			app.render(w, r, "signup.page.gohtml", TemplateData{
				Flash: "The email already exist, try login in instead",
			})
			return
		}
		user := User{}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			app.render(w, r, "signup.page.gohtml", TemplateData{
				Flash: "Internal error, try again",
			})
			return
		}

		user.PasswordHash = string(passwordHash)
		user.Email = email

		err = app.UserCreate(&user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	})
}

func (app *application) signin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			app.signinPost().ServeHTTP(w, r)
			return
		}

		app.render(w, r, "signin.page.gohtml", TemplateData{})
	})
}

func (app *application) signinPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.index().ServeHTTP(w, r)
			return
		}

		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user, err := app.UserFindByEmail(email)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.render(w, r, "signin.page.gohtml", TemplateData{
					Flash: "Wrong password or email",
				})
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

		if err != nil {
			app.render(w, r, "signin.page.gohtml", TemplateData{
				Flash: "Wrong password or email",
			})
			return
		}

		session := Session{}
		session.Cookie = generateCookie()
		session.Email = user.Email

		err = app.SessionCreate(&session)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		cookie := http.Cookie{
			Name:  "session",
			Value: session.Cookie,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
