package server

import (
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"text/template"
)

//go:embed templates/*
var staticContent embed.FS

func registerPaths(mux *http.ServeMux) error {
	if mux == nil {
		return errors.New("mux is nil")
	}
	tmpl, err := template.ParseFS(staticContent, "templates/*.tmpl")
	Server.Log.Debug("setup templates", slog.Any("templates", tmpl))
	if err != nil {
		Server.Log.Error("error parsing embedded templates", slog.Any("error", err))
		return err
	}

	mux.HandleFunc("/", logConnection(
		func(w http.ResponseWriter, r *http.Request) {
			data := struct{ Title, Content string }{Title: "Example", Content: "Hello, world!"}
			if err := tmpl.ExecuteTemplate(w, "index.tmpl", data); err != nil {
				Server.Log.Error("error executing template", slog.Any("error", err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}))
	return nil
}

func logConnection[T func(w http.ResponseWriter, r *http.Request)](f T) T {
	return func(w http.ResponseWriter, r *http.Request) {
		Server.Log.Debug("handling connection", slog.Any("request", r))
		f(w, r)
	}
}
