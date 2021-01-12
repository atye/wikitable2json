package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/atye/wikitable-api/internal/service/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	HTTPGet     func(string) (*http.Response, error)
	HTTPSvr     *http.Server
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
		HTTPGet: c.HTTPGet,
	}
	pb.RegisterWikipediaTableJSONAPIServer(c.GrpcSvr, svc)
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
		runtime.WithErrorHandler(fromStatusErrorHandler),
	)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	err = pb.RegisterWikipediaTableJSONAPIHandlerFromEndpoint(ctx, gwMux, "127.0.0.1:2000", opts)
	if err != nil {
		return err
	}
	mux.Handle("/api/", gwMux)
	c.HTTPSvr.Handler = mux
	eg.Go(func() error {
		return c.HTTPSvr.ListenAndServe()
	})
	if c.signalReady != nil {
		c.signalReady <- struct{}{}
	}
	return eg.Wait()
}

func fromStatusErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	s, _ := status.FromError(err)
	if s.Details() == nil || len(s.Details()) == 0 {
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	} else {
		switch d := s.Details()[0].(type) {
		case *errdetails.ErrorInfo:
			if respCodeStr, ok := d.Metadata["ResponseStatusCode"]; ok {
				respCode, err := strconv.Atoi(respCodeStr)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					s = status.New(codes.Internal, fmt.Sprintf("failed to process error response: %v", err))
				}
				w.WriteHeader(respCode)
			} else {
				w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
			}
		default:
			w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
		}
	}
	data, err := marshaler.Marshal(s.Proto())
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshaling error response: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", data)
}
