package entrypoint

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/atye/wikitable2json/internal/server"
)

//go:embed static/dist/*
var swagger embed.FS

type Config struct {
	Port   string
	Client server.TableGetter
}

func Run(c Config) error {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.StripPrefix("/", http.FileServer(getSwagger())))
	mux.Handle("GET /api/{page}", headerMW(server.NewServer(c.Client)))
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
