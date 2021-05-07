package main

import (
	"context"
	"log"
	"os"

	v2 "github.com/atye/wikitable-api/internal/service/v2"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}
	/*log.Fatal(service.Run(context.Background(), service.Config{
		HTTPGet: http.Get,
		HTTPSvr: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
		},
		GrpcSvr: grpc.NewServer(),
	}))*/

	log.Fatal(v2.Run(context.Background(), v2.Config{
		Port: port,
	}))
}
