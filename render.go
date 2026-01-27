package main

import (
	"html/template"
	"net/http"
	"path"
)

func (app *application) render(w http.ResponseWriter, filename string, data interface{}) {
	fullPath := path.Join(app.templateDir, filename)
	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
