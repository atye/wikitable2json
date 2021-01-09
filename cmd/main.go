package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atye/wikitable-api/internal/service"
	"google.golang.org/grpc"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}
	log.Fatal(service.Run(context.Background(), service.Config{
		HTTPGet: http.Get,
		HTTPSvr: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
		},
		GrpcSvr: grpc.NewServer(),
	}))
}
