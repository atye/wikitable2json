package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/atye/wikitable-api/service/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Config struct {
	HttpGet func(string) (*http.Response, error)
	HttpSvr *http.Server
	GrpcSvr *grpc.Server
}

func Run(ctx context.Context, c Config) error {
	var eg errgroup.Group

	lis, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	svc := &Service{
		HttpGet: c.HttpGet,
	}

	pb.RegisterWikiTableServer(c.GrpcSvr, svc)
	eg.Go(func() error {
		return c.GrpcSvr.Serve(lis)
	})

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger/apidocs.swagger.json")
	})
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(http.Dir("swagger/swagger-ui"))))

	gwMux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	err = pb.RegisterWikiTableHandlerFromEndpoint(ctx, gwMux, "127.0.0.1:2000", opts)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/api/", gwMux)

	runtime.HTTPError = customErrorHandler

	c.HttpSvr.Handler = mux

	eg.Go(func() error {
		return c.HttpSvr.ListenAndServe()
	})

	return eg.Wait()
}

func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	w.Header().Set("Content-type", marshaler.ContentType())

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		http.Error(w, "error getting server metadata", http.StatusInternalServerError)
		return
	}

	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting server status code: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(code)
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	} else {
		w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	}

	w.Write([]byte(grpc.ErrorDesc(err)))
}
