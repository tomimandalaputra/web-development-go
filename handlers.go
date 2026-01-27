package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// app.infoLog.Printf("Session data: %s", app.session.GetString(r, "userID"))
	app.render(w, r, "index.html", nil)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	// app.session.Put(r, "userID", "tomi")

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
			app.errorLog.Printf("Invalid form: %+v", form.Errors)
			app.render(w, r, "login.html", nil)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		app.infoLog.Printf("Ligged in with email %s; %s\n", email, password)
	}

	app.render(w, r, "login.html", nil)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.html", nil)
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
