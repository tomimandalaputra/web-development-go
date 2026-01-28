package main

import (
	"net/http"
)

const (
	loggedInUserKey = "logged_in_user"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "index.html", nil)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		form := NewForm(r.PostForm)
		form.Required("email", "password").
			MaxLength("email", 255).
			MaxLength("password", 20).
			MinLength("email", 3).
			MinLength("password", 6).
			Matches("email", EmailRx)

		if !form.Valid() {
			form.Errors.Add("generic", "The data you submitted was not valid")
			app.render(w, r, "login.html", &templateData{
				Form: form,
			})
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		_, err := app.userRepo.Authenticate(email, password)
		if err != nil {
			form.Errors.Add("generic", err.Error())
			app.render(w, r, "login.html", &templateData{
				Form: form,
			})
			return
		}

		app.session.Put(r, loggedInUserKey, email)
		app.session.Put(r, "flash", "You are logged in")

		http.Redirect(w, r, "/submit", http.StatusSeeOther)
		return
	}

	app.render(w, r, "login.html", &templateData{
		Form: NewForm(r.PostForm),
	})
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, loggedInUserKey)
	app.session.Put(r, "flash", "Your are logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		form := NewForm(r.PostForm)
		form.Required("email", "password", "name").
			MaxLength("email", 255).
			MaxLength("password", 20).
			MinLength("email", 3).
			MinLength("password", 6).
			MinLength("name", 3).
			Matches("email", EmailRx)

		if !form.Valid() {
			form.Errors.Add("generic", "The data you submitted was not valid")
			app.render(w, r, "register.html", &templateData{
				Form: form,
			})
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		name := r.FormValue("name")
		avatar := r.FormValue("avatar")

		_, err := app.userRepo.CreateUser(name, email, password, avatar)
		if err != nil {
			form.Errors.Add("generic", err.Error())
			app.render(w, r, "register.html", &templateData{
				Form: form,
			})
			return
		}

		app.session.Put(r, "flash", "You are registered")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.render(w, r, "register.html", &templateData{
		Form: NewForm(r.PostForm),
	})
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.html", nil)
}

func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "contact.html", nil)
}

func (app *application) submit(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "submit.html", nil)
}
