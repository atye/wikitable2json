package service

import (
	"fmt"
	"log"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type wikiApiError struct {
	statusCode int
	message    string
}

func (e *wikiApiError) Error() string {
	return ""
}

func wikiAPIStatusErr(apiErr *wikiApiError) error {
	st := status.New(codes.Aborted, apiErr.message)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikipedia.org/api/rest_v1/#",
		Reason: fmt.Sprintf("expected response code 200/OK from the wikipedia API, got %d", apiErr.statusCode),
	})
	if err != nil {
		log.Printf("failed to apply error details: %v", err)
		return status.Error(codes.Internal, "not sure what happened. Open an issue at https://github.com/atye/wikitable-api if you'd like.")
	}
	return st.Err()
}

func tableParseStatusErr(err error) error {
	st := status.New(codes.Internal, err.Error())
	st, err = st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikitable2json.com",
		Reason: "something unexpected was encountered while parsing tables",
	})
	if err != nil {
		log.Printf("failed to apply error details: %v", err)
		return status.Error(codes.Internal, "not sure what happened. Open an issue at https://github.com/atye/wikitable-api if you'd like.")
	}
	return st.Err()
}

func getDocumentStatusErr(err error) error {
	st := status.New(codes.Aborted, err.Error())
	st, err = st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikitable2json.com",
		Reason: "something unexpected was encountered while retrieving the wikipedia API response document",
	})
	if err != nil {
		log.Printf("failed to apply error details: %v", err)
		return status.Error(codes.Internal, "not sure what happened. Open an issue at https://github.com/atye/wikitable-api if you'd like.")
	}
	return st.Err()
}
