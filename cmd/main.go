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
	"github.com/atye/wikitable2json/internal/server/metrics"
	"github.com/atye/wikitable2json/pkg/client"
)

//go:embed static/dist/*
var swagger embed.FS

var (
	defaultCacheSize       = 20
	defaultCacheExpiration = 60 * time.Second
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		handleErr(fmt.Errorf("PORT env not set"))
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil || cacheSize == 0 {
		log.Printf("CACHE_SIZE is empty or invalid with error: %v; using %d", err, defaultCacheSize)
		cacheSize = defaultCacheSize
	}

	cacheExpiration, err := time.ParseDuration(os.Getenv("CACHE_EXPIRATION"))
	if err != nil || cacheExpiration == 0 {
		log.Printf("CACHE_EXPIRATION is empty or invalid with error: %v; using %s", err, defaultCacheExpiration)
		cacheExpiration = defaultCacheExpiration
	}

	googleMeasurementId := os.Getenv("GOOGLE_MEASUREMENT_ID")
	googleAPISecret := os.Getenv("GOOGLE_API_SECRET")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	var mp server.MetricsPublisher
	if googleMeasurementId != "" && googleAPISecret != "" {
		mp = metrics.NewGoogleClient(googleMeasurementId, googleAPISecret, httpClient)
	}

	dist, err := fs.Sub(swagger, "static/dist")
	if err != nil {
		handleErr(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.StripPrefix("/", http.FileServer(http.FS(dist))))
	mux.Handle("GET /api/{page}",
		server.HeaderMW(server.RequestValidationAndMetricsMW(
			server.NewServer(client.NewClient("", client.WithHTTPClient(httpClient)), server.NewCache(cacheSize, cacheExpiration)), mp)))
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
			log.Printf("main: shutting down server: %v\n", err)
			_ = svr.Close()
		}
	}
}

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
