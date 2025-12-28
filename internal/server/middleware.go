package server

import (
	"context"
	"log"
	"net/http"

	"github.com/atye/wikitable2json/pkg/client/status"
)

type contextKey string

var (
	pageKey  contextKey = "page"
	queryKey contextKey = "query"
)

func HeaderMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type MetricsPublisher interface {
	Publish(code int, ip string, page string, lang string) error
}

func RequestValidationAndMetricsMW(main http.Handler, mp MetricsPublisher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.PathValue("page")
		if page == "" {
			writeError(w, status.NewStatus("page value must be supplied in /api/{page}", http.StatusBadRequest))
			return
		}

		qv, err := parseParameters(r)
		if err != nil {
			writeError(w, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, pageKey, page)
		ctx = context.WithValue(ctx, queryKey, qv)
		r = r.WithContext(ctx)

		rec := &statusRecorder{ResponseWriter: w, Status: http.StatusOK}
		main.ServeHTTP(rec, r)

		if mp != nil {
			go func() {
				err := mp.Publish(rec.Status, r.RemoteAddr, page, qv.lang)
				if err != nil {
					log.Printf("publishing metric: %v\n", err)
				}
			}()
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.Status = code
	rec.ResponseWriter.WriteHeader(code)
}
