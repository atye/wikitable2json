package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atye/wikitable-api/service"
	"google.golang.org/grpc"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}

	conf := service.Config{
		HttpGet: http.Get,
		HttpSvr: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
		},
		GrpcSvr: grpc.NewServer(),
	}

	log.Fatal(service.Run(context.Background(), conf))
}
