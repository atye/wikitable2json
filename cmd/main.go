package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/atye/wikitable2json/internal/server"
	"github.com/atye/wikitable2json/pkg/client"
)

//go:embed static/dist/*
var swagger embed.FS

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		handleErr(fmt.Errorf("PORT env not set"))
	}

	cacheSize, ok := os.LookupEnv("CACHE_SIZE")
	if !ok {
		handleErr(fmt.Errorf("CACHE_SIZE env not set"))
	}

	cacheExpiration, ok := os.LookupEnv("CACHE_EXPIRATION")
	if !ok {
		handleErr(fmt.Errorf("CACHE_EXPIRATION env not set"))
	}

	size, err := strconv.Atoi(cacheSize)
	if err != nil {
		handleErr(fmt.Errorf("parsing CACHE_SIZE %s: %v", cacheSize, err))
	}

	dur, err := time.ParseDuration(cacheExpiration)
	if err != nil {
		handleErr(fmt.Errorf("parsing CACHE_EXPIRATION %s: %v", cacheExpiration, err))
	}

	dist, err := fs.Sub(swagger, "static/dist")
	if err != nil {
		handleErr(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.StripPrefix("/", http.FileServer(http.FS(dist))))
	mux.Handle("GET /api/{page}", server.HeaderMW(server.NewServer(client.NewClient("", client.WithHTTPClient(&http.Client{Timeout: 10 * time.Second})), server.NewCache(size, dur))))
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	httpErrors := make(chan error, 1)
	go func() {
		httpErrors <- svr.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-httpErrors:
		handleErr(err)
	case <-shutdown:
		log.Println("main: handling shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := svr.Shutdown(ctx)
		if err != nil {
			log.Printf("main: shutting down server: %v", err)
			_ = svr.Close()
		}
	}
}

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
