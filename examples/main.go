package main

import (
	"context"
	"fmt"
	"log"

	"github.com/atye/wikitable2json/pkg/client"
)

func main() {
	tg := client.NewTableGetter("github.com/atye/wikitable2json/examples")

	// You can create a client with a cache (size, item expiration, interval to check expiration of each item)
	// tg := client.NewTableGetter("github.com/atye/wikitable2json/examples", client.WithCache(cache.New(5, 5*time.Second, 5*time.Second)))

	matrix, err := tg.GetTablesMatrix(context.Background(), "Arhaan_Khan", "en", false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matrix)

	matrix, err = tg.GetTablesMatrix(context.Background(), "Arhaan_Khan", "en", false, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matrix)

	keyValue, err := tg.GetTablesKeyValue(context.Background(), "Arhaan_Khan", "en", false, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(keyValue)
}
