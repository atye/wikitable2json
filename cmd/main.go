package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/atye/wikitable-api/service"
	"github.com/atye/wikitable-api/service/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	errCh := make(chan error)

	lis, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	svc := &service.Service{}
	svr := grpc.NewServer()

	pb.RegisterWikiTableServer(svr, svc)
	go func() {
		errCh <- svr.Serve(lis)
	}()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	mux := runtime.NewServeMux()

	err = pb.RegisterWikiTableHandlerFromEndpoint(context.Background(), mux, "127.0.0.1:2000", opts)
	if err != nil {
		log.Fatal(err)
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env not set")
	}

	go func() {
		errCh <- http.ListenAndServe(":"+port, mux)
	}()

	log.Fatal(<-errCh)
}
