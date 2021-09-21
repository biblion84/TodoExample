package main

import (
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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

func (app *application) checkTodoPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.Form.Get("id")
		state := r.Form.Get("state")
		var todo Todo
		app.db.Find(&todo, id)
		if todo.ID != 0 {
			todo.Checked = state == "true"
			app.db.Save(&todo)
		}
		app.index().ServeHTTP(w, r)
	})
}

func (app *application) signup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "signup.page.gohtml", TemplateData{})
	})
}

func (app *application) signupPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		var user User
		app.db.First(&user, User{Email: email})
		if user.ID != 0 {
			app.render(w, r, "signup.page.gohtml", TemplateData{
				Flash: "The email already exist, try login in instead",
			})
			return
		}
		user = User{}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			app.render(w, r, "signup.page.gohtml", TemplateData{
				Flash: "Internal error, try again",
			})
			return
		}

		user.PasswordHash = string(passwordHash)
		user.Email = email

		app.db.Create(&user)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
func (app *application) signin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "signin.page.gohtml", TemplateData{})
	})
}

func (app *application) signinPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		var user User
		app.db.First(&user, User{Email: email})
		if user.ID == 0 {
			app.render(w, r, "signin.page.gohtml", TemplateData{
				Flash: "Wrong password or email",
			})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			app.render(w, r, "signin.page.gohtml", TemplateData{
				Flash: "Wrong password or email",
			})
			return
		}

		session := Session{}
		session.Cookie = generateCookie()
		session.Email = user.Email
		app.db.Create(&session)
		cookie := http.Cookie{
			Name:  "session",
			Value: session.Cookie,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
