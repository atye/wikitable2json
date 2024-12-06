package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static/dist/*
var swagger embed.FS

type Config struct {
	Port   string
	Client tableGetter
	Cache  *cache
}

func Run(c Config) error {
	dist, err := fs.Sub(swagger, "static/dist")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.StripPrefix("/", http.FileServer(http.FS(dist))))
	mux.Handle("GET /api/{page}", headerMW(newServer(c.Client, c.Cache)))
	return http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux)
}

func headerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
