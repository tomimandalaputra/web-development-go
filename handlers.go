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
	homeContent := fmt.Sprintf(htmlContent, "Home", "<h1>Welcome to my web app</h1>")
	_, _ = w.Write([]byte(homeContent))
}
