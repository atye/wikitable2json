package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atye/wikitable-api/internal/service"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/api/", service.NewServer(service.BaseURL))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
