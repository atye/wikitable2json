package service

import (
	"context"
	"net"
	"net/http"

	"github.com/atye/wikitable-api/service/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Config struct {
	HttpGet     func(string) (*http.Response, error)
	HttpSvr     *http.Server
	GrpcSvr     *grpc.Server
	signalReady chan struct{}
}

func Run(ctx context.Context, c Config) error {
	var eg errgroup.Group

	lis, err := net.Listen("tcp", ":2000")
	if err != nil {
		return err
	}

	svc := &Service{
		HttpGet: c.HttpGet,
	}

	pb.RegisterWikiTableJSONAPIServer(c.GrpcSvr, svc)
	eg.Go(func() error {
		return c.GrpcSvr.Serve(lis)
	})

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger/wikitable.swagger.json")
	})
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(http.Dir("swagger/swagger-ui"))))

	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(fromStatusWithDetailsErrorHandler),
	)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	err = pb.RegisterWikiTableJSONAPIHandlerFromEndpoint(ctx, gwMux, "127.0.0.1:2000", opts)
	if err != nil {
		return err
	}

	mux.Handle("/api/", gwMux)

	c.HttpSvr.Handler = mux

	eg.Go(func() error {
		return c.HttpSvr.ListenAndServe()
	})

	if c.signalReady != nil {
		c.signalReady <- struct{}{}
	}

	return eg.Wait()
}
