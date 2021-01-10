package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func wikiAPIRespNotOKStatusErr(apiErr *wikiApiError) error {
	st := status.New(codes.Internal, apiErr.message)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikipedia.org/api/rest_v1/#",
		Reason: fmt.Sprintf("expected response status 200/OK from the wikipedia API, got something else"),
		Metadata: map[string]string{
			"ResponseStatusCode": strconv.Itoa(apiErr.statusCode),
			"ResponseStatusText": http.StatusText(apiErr.statusCode),
			"Page":               apiErr.page,
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

func getGeneralStatusErr(err error, reason string) error {
	st := status.New(codes.Internal, err.Error())
	st, err = st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikitable2json.com",
		Reason: reason,
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
