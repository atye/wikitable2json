package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atye/wikitable2json/internal/entrypoint"
	"github.com/atye/wikitable2json/internal/server/api"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}

	c := entrypoint.Config{
		Port:    port,
		WikiAPI: api.NewWikiClient(api.BaseURL),
	}

	if err := entrypoint.Run(c); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
