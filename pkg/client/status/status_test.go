package status

import (
	"net/http"
	"strings"
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

	if !strings.Contains(got, "message") {
		t.Errorf("want %s to contain %s", want, "message")
	}

	if !strings.Contains(got, "TableIndex: 0") {
		t.Errorf("want %s to contain %s", want, "TableIndex: 0")
	}

	if !strings.Contains(got, "RowIndex: 0") {
		t.Errorf("want %s to contain %s", want, "RowIndex: 0")
	}

	if !strings.Contains(got, "ColumnIndex: 0") {
		t.Errorf("want %s to contain %s", want, "ColumnIndex: 0")
	}
}
