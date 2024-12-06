package main

import (
	"context"
	"fmt"
	"log"

	"github.com/atye/wikitable2json/pkg/client"
)

func main() {
	tg := client.NewClient("github.com/atye/wikitable2json/examples")

	// You can also create a client with a cache (capacity, page expiration, interval to check expiration of each page)
	// tg = client.NewTableGetter("github.com/atye/wikitable2json/examples", client.WithCache(5, 5*time.Second, 5*time.Second))

	matrix, err := tg.GetMatrix(context.Background(), "Arhaan_Khan", "en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matrix)

	matrix, err = tg.GetMatrix(context.Background(), "Arhaan_Khan", "en", client.WithTables(0))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matrix)

	matrixVerbose, err := tg.GetMatrixVerbose(context.Background(), "Arhaan_Khan", "en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matrixVerbose)

	keyValue, err := tg.GetKeyValue(context.Background(), "Arhaan_Khan", "en", 1, client.WithTables(1))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(keyValue)

	keyValueVerbose, err := tg.GetKeyValueVerbose(context.Background(), "Arhaan_Khan", "en", 1, client.WithTables(1))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(keyValueVerbose)
}
