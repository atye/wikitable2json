package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/atye/wikitable2json/pkg/client/status"
)

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/golden":
			w.Write(getPageBytes(t, "golden"))
		case "/goldenDouble":
			w.Write(getPageBytes(t, "goldenDouble"))
		case "/issueOne":
			w.Write(getPageBytes(t, "issueOne"))
		case "/dataSortValue":
			w.Write(getPageBytes(t, "dataSortValue"))
		case "/allTableClasses":
			w.Write(getPageBytes(t, "allTableClasses"))
		case "/badRowSpan":
			w.Write(getPageBytes(t, "badRowSpan"))
		case "/badColSpan":
			w.Write(getPageBytes(t, "badColSpan"))
		case "/issue34":
			w.Write(getPageBytes(t, "issue34"))
		case "/issue56":
			w.Write(getPageBytes(t, "issue56"))
		case "/issue77":
			w.Write(getPageBytes(t, "issue77"))
		case "/issue85":
			w.Write(getPageBytes(t, "issue85"))
		case "/issue93":
			w.Write(getPageBytes(t, "issue93"))
		case "/issue105":
			w.Write(getPageBytes(t, "issue105"))
		case "/reference":
			w.Write(getPageBytes(t, "reference"))
		case "/simpleKeyValue":
			w.Write(getPageBytes(t, "simpleKeyValue"))
		case "/complexKeyValue":
			w.Write(getPageBytes(t, "complexKeyValue"))
		case "/keyValueBadRows":
			w.Write(getPageBytes(t, "keyValueBadRows"))
		case "/keyValueOneRow":
			w.Write(getPageBytes(t, "keyValueOneRow"))
		case "/noTables":
			w.Write(getPageBytes(t, "noTables"))
		case "/StatusRequestEntityTooLarge":
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("StatusRequestEntityTooLarge"))
		case "/UserAgent":
			got := r.Header.Get("User-Agent")
			want := "test@mail.com"
			if want != got {
				t.Errorf("want %s, got %s", want, got)
			}
		case "/NoUserAgent":
			got := r.Header.Get("User-Agent")
			want := "github.com/atye/wikitable2json"
			if want != got {
				t.Errorf("want %s, got %s", want, got)
			}
		default:
			t.Fatalf("path %s not supported", r.URL.Path)
		}
	}))
	defer ts.Close()

	originalgetApiURLFn := getApiURLFn
	getApiURLFn = func(_, page string) string { return fmt.Sprintf("%s/%s", ts.URL, page) }
	defer func() {
		getApiURLFn = originalgetApiURLFn
	}()

	sut := NewClient("", WithHTTPClient(&http.Client{Timeout: 1 * time.Second}))

	t.Run("Matrix", func(t *testing.T) {
		tests := []struct {
			page         string
			options      []TableOption
			want         [][][]string
			wantErr      bool
			wantErrValue status.Status
		}{
			{
				"golden",
				nil,
				GoldenMatrix,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				nil,
				GoldenMatrixDouble,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithTables(0)},
				GoldenMatrix,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithTables(1)},
				GoldenMatrix,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithTables(0, 1)},
				GoldenMatrixDouble,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithSections("First")},
				GoldenMatrix,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithSections("Second_Table")},
				GoldenMatrix,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithSections("First", "Second_Table")},
				GoldenMatrixDouble,
				false,
				status.Status{},
			},
			{
				"goldenDouble",
				[]TableOption{WithTables(0), WithSections("Second_Table")},
				GoldenMatrixDouble,
				false,
				status.Status{},
			},
			{
				"allTableClasses",
				nil,
				AllTableClasses,
				false,
				status.Status{},
			},
			{
				"reference",
				[]TableOption{WithCleanReferences()},
				ReferenceMatrix,
				false,
				status.Status{},
			},
			{
				"issueOne",
				nil,
				IssueOneMatrix,
				false,
				status.Status{},
			},
			{
				"dataSortValue",
				nil,
				DataSortValueMatrix,
				false,
				status.Status{},
			},
			{
				"issue34",
				nil,
				Issue34Matrix,
				false,
				status.Status{},
			},
			{
				"issue56",
				[]TableOption{WithCleanReferences()},
				Issue56Matrix,
				false,
				status.Status{},
			},
			{
				"issue77",
				nil,
				Issue77Matrix,
				false,
				status.Status{},
			},
			{
				"issue105",
				[]TableOption{WithBRNewLine()},
				Issue105Matrix,
				false,
				status.Status{},
			},
			{
				"noTables",
				nil,
				NoTablesMatrix,
				false,
				status.Status{},
			},
			{
				"StatusRequestEntityTooLarge",
				nil,
				nil,
				true,
				status.Status{
					Message: "StatusRequestEntityTooLarge",
					Code:    http.StatusRequestEntityTooLarge,
					Details: status.Details{
						status.Page: "StatusRequestEntityTooLarge",
					},
				},
			},
			{
				"badRowSpan",
				nil,
				nil,
				true,
				status.NewStatus("no integer value in span attribute: [x]", http.StatusInternalServerError, status.WithDetails(
					status.Details{
						status.TableIndex:  0,
						status.RowIndex:    1,
						status.ColumnIndex: 0,
					},
				)),
			},
			{
				"badColSpan",
				nil,
				nil,
				true,
				status.Status{
					Message: "no integer value in span attribute: [x]",
					Code:    http.StatusInternalServerError,
					Details: status.Details{
						status.TableIndex:  0,
						status.RowIndex:    1,
						status.ColumnIndex: 1,
					},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.page, func(t *testing.T) {
				got, err := sut.GetMatrix(context.Background(), tc.page, "en", tc.options...)
				if tc.wantErr {
					if err == nil {
						t.Fatal("expected error, got nil")
					}

					if !reflect.DeepEqual(tc.wantErrValue, err.(status.Status)) {
						t.Errorf("want %v\n got %v", tc.wantErrValue, err)
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				}
			})
		}
	})

	t.Run("MatrixVerbose", func(t *testing.T) {
		tests := []struct {
			page         string
			options      []TableOption
			want         [][][]Verbose
			wantErr      bool
			wantErrValue status.Status
		}{
			{
				"issue77",
				nil,
				Issue77MatrixVerbose,
				false,
				status.Status{},
			},
			{
				"issue93",
				nil,
				Issue93MatrixVerbose,
				false,
				status.Status{},
			},
			{
				"issue105",
				[]TableOption{WithBRNewLine()},
				Issue105MatrixVerbose,
				false,
				status.Status{},
			},
		}

		for _, tc := range tests {
			t.Run(tc.page, func(t *testing.T) {
				got, err := sut.GetMatrixVerbose(context.Background(), tc.page, "en", tc.options...)
				if tc.wantErr {
					if err == nil {
						t.Fatal("expected error, got nil")
					}

					if !reflect.DeepEqual(tc.wantErrValue, err.(status.Status)) {
						t.Errorf("want %v\n got %v", tc.wantErrValue, err)
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				}
			})
		}
	})

	t.Run("KeyValue", func(t *testing.T) {
		tests := []struct {
			page         string
			options      []TableOption
			keyRows      int
			want         [][]map[string]string
			wantErr      bool
			wantErrValue status.Status
		}{
			{
				"simpleKeyValue",
				nil,
				1,
				SimpleKeyValue,
				false,
				status.Status{},
			},
			{
				"complexKeyValue",
				[]TableOption{WithCleanReferences()},
				2,
				ComplexKeyValue,
				false,
				status.Status{},
			},
			{
				"issue85",
				nil,
				1,
				Issue85KeyValue,
				false,
				status.Status{},
			},
			{
				"keyValueOneRow",
				nil,
				1,
				nil,
				true,
				status.Status{
					Message: "table needs at least two rows",
					Code:    http.StatusBadRequest,
					Details: status.Details{
						status.TableIndex: 0,
					},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.page, func(t *testing.T) {
				got, err := sut.GetKeyValue(context.Background(), tc.page, "en", tc.keyRows, tc.options...)
				if tc.wantErr {
					if err == nil {
						t.Fatal("expected error, got nil")
					}

					if !reflect.DeepEqual(tc.wantErrValue, err.(status.Status)) {
						t.Errorf("want %v\n got %v", tc.wantErrValue, err)
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				}
			})
		}
	})

	t.Run("KeyValueVerbose", func(t *testing.T) {
		tests := []struct {
			page         string
			options      []TableOption
			keyRows      int
			want         [][]map[string]Verbose
			wantErr      bool
			wantErrValue status.Status
		}{
			{
				"issue77",
				nil,
				1,
				Issue77KeyValueVerbose,
				false,
				status.Status{},
			},
			{
				"issue93",
				nil,
				1,
				Issue93KeyValueVerbose,
				false,
				status.Status{},
			},
		}

		for _, tc := range tests {
			t.Run(tc.page, func(t *testing.T) {
				got, err := sut.GetKeyValueVerbose(context.Background(), tc.page, "en", tc.keyRows, tc.options...)
				if tc.wantErr {
					if err == nil {
						t.Fatal("expected error, got nil")
					}

					if !reflect.DeepEqual(tc.wantErrValue, err.(status.Status)) {
						t.Errorf("want %v\n got %v", tc.wantErrValue, err)
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				}
			})
		}
	})

	t.Run("UserAgent", func(t *testing.T) {
		sut.SetUserAgent("test@mail.com")
		defer sut.SetUserAgent("")

		_, err := sut.GetMatrix(context.Background(), "UserAgent", "en")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NoUserAgent", func(t *testing.T) {
		_, err := sut.GetMatrix(context.Background(), "NoUserAgent", "en")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func getPageBytes(t *testing.T, page string) []byte {
	t.Helper()

	f, err := os.ReadFile(fmt.Sprintf("testdata/%s.html", page))
	if err != nil {
		t.Fatal(err)
	}
	return f
}
