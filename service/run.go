package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/atye/wikitable-api/service/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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

	gwMux := runtime.NewServeMux()

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

func httpHeaderModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("error generating response status code: %v", err), http.StatusInternalServerError)
			return nil
		}
		w.WriteHeader(code)
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	}
	return nil
}

func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return
	}
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("error generating response status code: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(code)
	} else {
		w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	}
	delete(md.HeaderMD, "x-http-code")
	delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	w.Write([]byte(grpc.ErrorDesc(err)))
}
