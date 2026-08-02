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
	defaultRateLimit       = 180

	defaultUserAgent = "github.com/atye/wikitable2json"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		handleErr(fmt.Errorf("PORT env not set"))
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil || cacheSize == 0 {
		log.Printf("CACHE_SIZE env is empty or invalid with error: %v; using %d", err, defaultCacheSize)
		cacheSize = defaultCacheSize
	}

	cacheExpiration, err := time.ParseDuration(os.Getenv("CACHE_EXPIRATION"))
	if err != nil || cacheExpiration == 0 {
		log.Printf("CACHE_EXPIRATION env is empty or invalid with error: %v; using %s", err, defaultCacheExpiration)
		cacheExpiration = defaultCacheExpiration
	}

	rateLimit, err := strconv.Atoi(os.Getenv("WIKIPEDIA_RATE_LIMIT"))
	if err != nil || rateLimit == 0 {
		log.Printf("WIKIPEDIA_RATE_LIMIT env is empty or invalid with error: %v; using %d", err, defaultRateLimit)
		rateLimit = defaultRateLimit
	}

	userAgent := os.Getenv("USER_AGENT")
	if userAgent == "" {
		log.Printf("USER_AGENT env is empty; using %s", defaultUserAgent)
		userAgent = defaultUserAgent
	}

	googleMeasurementId := os.Getenv("GOOGLE_MEASUREMENT_ID")
	googleAPISecret := os.Getenv("GOOGLE_API_SECRET")

	var mp server.MetricsPublisher
	if googleMeasurementId != "" && googleAPISecret != "" {
		mp = metrics.NewGoogleClient(googleMeasurementId, googleAPISecret, &http.Client{Timeout: 5 * time.Second})
	}

	dist, err := fs.Sub(swagger, "static/dist")
	if err != nil {
		handleErr(err)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	app, err := server.NewServer(
		client.NewClient(userAgent, client.WithHTTPClient(httpClient), client.WithRateLimit(rateLimit)),
		server.NewCache(cacheSize, cacheExpiration))
	if err != nil {
		handleErr(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.StripPrefix("/", http.FileServer(http.FS(dist))))
	mux.Handle("GET /api/{page}", server.HeaderMW(server.RequestValidationAndMetricsMW(app, mp)))
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
