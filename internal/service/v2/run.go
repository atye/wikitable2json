package v2

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	Port string
}

func Run(ctx context.Context, conf Config) error {
	mux := http.NewServeMux()
	mux.Handle("/api/", &Server{})
	return http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), mux)
}
