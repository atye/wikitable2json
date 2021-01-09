package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fromStatusWithDetailsErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	if s, ok := status.FromError(err); ok {
		if s.Details() != nil && len(s.Details()) > 0 {
			switch d := s.Details()[0].(type) {
			case *errdetails.ErrorInfo:
				respCodeStr, ok := d.Metadata["ResponseStatusCode"]
				if !ok {
					panic("response status code wasn't included in the error details")
				}
				respCode, err := strconv.Atoi(respCodeStr)
				if err != nil {
					http.Error(w, fmt.Sprintf("error processing response code : %v", err.Error()), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(respCode)
			}
		} else {
			http.Error(w, fmt.Sprintf("encountered an error but no other details were provided: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		data, err := marshaler.Marshal(s.Proto())
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshaling error response: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", data)
	} else {
		fmt.Fprintf(w, "encountered an error but no other details were provided: %v", err.Error())
	}
}

type wikiApiError struct {
	statusCode int
	message    string
}

func (e *wikiApiError) Error() string {
	return ""
}

func wikiAPIStatusErr(apiErr *wikiApiError) error {
	st := status.New(codes.Internal, apiErr.message)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikipedia.org/api/rest_v1/#",
		Reason: fmt.Sprintf("expected response status 200/OK from the wikipedia API, got something else"),
		Metadata: map[string]string{
			"ResponseStatusCode": strconv.Itoa(apiErr.statusCode),
			"ResponseStatusText": http.StatusText(apiErr.statusCode),
		},
	})
	if err != nil {
		log.Printf("failed to apply error details: %v", err)
		return status.Error(codes.Internal, "not sure what happened. Open an issue at https://github.com/atye/wikitable-api if you'd like.")
	}
	return st.Err()
}

func tableParseStatusErr(ptErr *parseTableError) error {
	st := status.New(codes.Internal, ptErr.err.Error())
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikitable2json.com",
		Reason: "something unexpected was encountered while parsing tables",
		Metadata: map[string]string{
			"ResponseStatusCode": strconv.Itoa(http.StatusInternalServerError),
			"ResponseStatusText": http.StatusText(http.StatusInternalServerError),
			"TableIndex":         strconv.Itoa(ptErr.tableIndex),
			"RowNumber":          strconv.Itoa(ptErr.rowNum),
			"CellNumber":         strconv.Itoa(ptErr.cellNum),
		},
	})
	if err != nil {
		log.Printf("failed to apply error details: %v", err)
		return status.Error(codes.Internal, "not sure what happened. Open an issue at https://github.com/atye/wikitable-api if you'd like.")
	}
	return st.Err()
}

func getDocumentStatusErr(err error) error {
	st := status.New(codes.Internal, err.Error())
	st, err = st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikitable2json.com",
		Reason: "something unexpected was encountered while retrieving the wikipedia API response document",
		Metadata: map[string]string{
			"ResponseStatusCode": strconv.Itoa(http.StatusInternalServerError),
			"ResponseStatusText": http.StatusText(http.StatusInternalServerError),
		},
	})
	if err != nil {
		log.Printf("failed to apply error details: %v", err)
		return status.Error(codes.Internal, "not sure what happened. Open an issue at https://github.com/atye/wikitable-api if you'd like.")
	}
	return st.Err()
}

type parseTableError struct {
	err        error
	tableIndex int
	rowNum     int
	cellNum    int
}

func (e *parseTableError) Error() string {
	return ""
}
