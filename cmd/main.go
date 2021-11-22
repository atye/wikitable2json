package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atye/wikitable-api/internal/entrypoint"
	"github.com/atye/wikitable-api/internal/server/data"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}

	c := entrypoint.Config{
		Port:    port,
		WikiAPI: data.NewWikiClient(data.BaseURL),
	}

	if err := entrypoint.Run(c); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
