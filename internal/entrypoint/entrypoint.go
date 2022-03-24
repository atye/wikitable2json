package entrypoint

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/atye/wikitable-api/internal/server"
)

//go:embed static/dist/*
var swagger embed.FS

type Config struct {
	Port    string
	WikiAPI server.WikiAPI
}

func Run(c Config) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(getSwagger())))
	mux.Handle("/api/", headerMW(server.NewServer(c.WikiAPI)))
	return http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux)
}

func getSwagger() http.FileSystem {
	dist, err := fs.Sub(swagger, "static/dist")
	if err != nil {
		panic(err)
	}
	return http.FS(dist)
}

func headerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
