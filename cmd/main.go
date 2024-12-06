package main

import (
	"fmt"
	"os"
	"time"

	"github.com/atye/wikitable2json/internal/server"
	"github.com/atye/wikitable2json/pkg/client"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		fmt.Fprintf(os.Stderr, "%s", "PORT env not set")
		os.Exit(1)
	}

	c := server.Config{
		Port:   port,
		Client: client.NewClient(""),
		Cache:  server.NewCache(10, 8*time.Second, 8*time.Second),
	}

	if err := server.Run(c); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
