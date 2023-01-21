package status

import (
	"net/http"
	"testing"
)

func TestStatus(t *testing.T) {
	s := NewStatus("message", http.StatusInternalServerError, WithDetails(Details{
		TableIndex:  0,
		RowIndex:    0,
		ColumnIndex: 0,
	}))

	want := "message, TableIndex: 0, RowIndex: 0, ColumnIndex: 0"
	got := s.Error()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}
