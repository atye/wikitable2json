package main

import (
	"context"
	"fmt"
	"log"

	"github.com/atye/wikitable2json/pkg/client"
)

func main() {
	tg := client.NewClient("user@email.com")

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

	matrix, err = tg.GetMatrix(context.Background(), "Arhaan_Khan", "en", client.WithSections("Film"))
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
