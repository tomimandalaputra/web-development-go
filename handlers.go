package main

import (
	"fmt"
	"net/http"
)

var htmlContent = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>%s</title>
		</head>
		<body>
			%s
		</body>
		</html>
		`

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Printf("%s %s", r.Method, r.URL.Path)
	app.render(w, "index.html", nil)
}

func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Printf("%s %s", r.Method, r.URL.Path)
	homeContent := fmt.Sprintf(htmlContent, "Contact", "<h1>Contact page</h1>")
	_, _ = w.Write([]byte(homeContent))
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Printf("%s %s", r.Method, r.URL.Path)
	homeContent := fmt.Sprintf(htmlContent, "About", "<h1>About page</h1>")
	_, _ = w.Write([]byte(homeContent))
}
