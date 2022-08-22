package main

import (
	"context"
	"fmt"
	"log"

	"github.com/atye/wikitable2json/pkg/client"
)

func main() {
	tg := client.NewTableGetter("github.com/atye/wikitable2json/examples")

	// You can also create a client with a cache (capacity, page expiration, interval to check expiration of each page)
	// tg = client.NewTableGetter("github.com/atye/wikitable2json/examples", client.WithCache(5, 5*time.Second, 5*time.Second))

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
