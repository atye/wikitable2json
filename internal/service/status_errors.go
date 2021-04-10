package service

import (
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func wikiAPIStatusErr(apiErr *wikiApiError) error {
	st := status.New(codes.Internal, apiErr.message)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikipedia.org/api/rest_v1/#",
		Reason: fmt.Sprintf("expected response status 200/OK from the wikipedia API, got %d/%s", apiErr.statusCode, http.StatusText(apiErr.statusCode)),
		Metadata: map[string]string{
			"ResponseStatusCode": strconv.Itoa(apiErr.statusCode),
			"ResponseStatusText": http.StatusText(apiErr.statusCode),
			"Page":               apiErr.page,
		},
	})
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to process error response: %v", err))
	}
	return st.Err()
}

func tableParseStatusErr(ptErr *parseTableError) error {
	st := status.New(codes.Internal, ptErr.err.Error())
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: "wikitable2json.com",
		Reason: "something unexpected was encountered while parsing tables",
		Metadata: map[string]string{
			"TableIndex": strconv.Itoa(ptErr.tableIndex),
			"RowIndex":   strconv.Itoa(ptErr.rowNum),
			"CellIndex":  strconv.Itoa(ptErr.cellNum),
		},
	})
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to process error response: %v", err))
	}
	return st.Err()
}
