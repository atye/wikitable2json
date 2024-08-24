package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/atye/wikitable2json/internal/entrypoint"
	"github.com/atye/wikitable2json/pkg/client"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}

	c := entrypoint.Config{
		Port:   port,
		Client: client.NewTableGetter("", client.WithCache(10, 8*time.Second, 8*time.Second)),
	}

	if err := entrypoint.Run(c); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
