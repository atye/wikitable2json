package service

import (
	"net/http"
)

type ServerError struct {
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (e *ServerError) Error() string {
	return e.Message
}

func wikiAPIErr(e *wikiApiError) *ServerError {
	return &ServerError{
		Message: e.Error(),
		Metadata: map[string]interface{}{
			"ResponseStatusCode": e.statusCode,
			"ResponseStatusText": http.StatusText(e.statusCode),
			"Page":               e.page,
		},
	}
}

func tableParseErr(e *parseTableError) *ServerError {
	return &ServerError{
		Message: e.err.Error(),
		Metadata: map[string]interface{}{
			"ResponseStatusCode": http.StatusInternalServerError,
			"ResponseStatusText": http.StatusText(http.StatusInternalServerError),
			"TableIndex":         e.tableIndex,
			"RowNumber":          e.rowNum,
			"CellNumber":         e.cellNum,
		},
	}
}

func generalErr(e error, code int) *ServerError {
	return &ServerError{
		Message: e.Error(),
		Metadata: map[string]interface{}{
			"ResponseStatusCode": code,
			"ResponseStatusText": http.StatusText(code),
		},
	}
}
