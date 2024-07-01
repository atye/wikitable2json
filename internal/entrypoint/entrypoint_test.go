package entrypoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/atye/wikitable2json/internal/api"
	"github.com/atye/wikitable2json/internal/status"
	"github.com/atye/wikitable2json/pkg/client"
)

var (
	PORT = "8080"
)

func TestAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/rest_v1/page/html/golden":
			w.Write(getPageBytes(t, "golden"))
		case "/api/rest_v1/page/html/issueOne":
			w.Write(getPageBytes(t, "issueOne"))
		case "/api/rest_v1/page/html/dataSortValue":
			w.Write(getPageBytes(t, "dataSortValue"))
		case "/api/rest_v1/page/html/allTableClasses":
			w.Write(getPageBytes(t, "allTableClasses"))
		case "/api/rest_v1/page/html/badRowSpan":
			w.Write(getPageBytes(t, "badRowSpan"))
		case "/api/rest_v1/page/html/badColSpan":
			w.Write(getPageBytes(t, "badColSpan"))
		case "/api/rest_v1/page/html/issue34":
			w.Write(getPageBytes(t, "issue34"))
		case "/api/rest_v1/page/html/issue56":
			w.Write(getPageBytes(t, "issue56"))
		case "/api/rest_v1/page/html/issue77":
			w.Write(getPageBytes(t, "issue77"))
		case "/api/rest_v1/page/html/issue85":
			w.Write(getPageBytes(t, "issue85"))
		case "/api/rest_v1/page/html/issue93":
			w.Write(getPageBytes(t, "issue93"))
		case "/api/rest_v1/page/html/reference":
			w.Write(getPageBytes(t, "reference"))
		case "/api/rest_v1/page/html/simpleKeyValue":
			w.Write(getPageBytes(t, "simpleKeyValue"))
		case "/api/rest_v1/page/html/complexKeyValue":
			w.Write(getPageBytes(t, "complexKeyValue"))
		case "/api/rest_v1/page/html/keyValueBadRows":
			w.Write(getPageBytes(t, "keyValueBadRows"))
		case "/api/rest_v1/page/html/keyValueOneRow":
			w.Write(getPageBytes(t, "keyValueOneRow"))
		case "/api/rest_v1/page/html/noTables":
			w.Write(getPageBytes(t, "noTables"))
		case "/api/rest_v1/page/html/StatusRequestEntityTooLarge":
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("StatusRequestEntityTooLarge"))
		case "/api/rest_v1/page/html/UserAgent":
			got := r.Header.Get("User-Agent")
			want := "test@mail.com"
			if want != got {
				t.Errorf("want %s, got %s", want, got)
			}
		case "/api/rest_v1/page/html/NoUserAgent":
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

	originalBaseURL := api.BaseURL
	api.BaseURL = ts.URL
	defer func() {
		api.BaseURL = originalBaseURL
	}()

	go Run(Config{
		Port:   PORT,
		Client: client.NewTableGetter("", client.WithCache(3, 500*time.Millisecond, 500*time.Millisecond), client.WithHTTPClient(http.DefaultClient)),
	})

	waitforServer()

	t.Run("Success", func(t *testing.T) {
		t.Run("Matrix", func(t *testing.T) {
			tests := []struct {
				name string
				url  string
				want [][][]string
			}{
				{
					"Golden",
					fmt.Sprintf("http://localhost:%s/api/golden", PORT),
					GoldenMatrix,
				},
				{
					"GoldenWithParameters",
					fmt.Sprintf("http://localhost:%s/api/golden?lang=sp&format=matrix&table=0", PORT),
					GoldenMatrix,
				},
				{
					"AllTableClasses",
					fmt.Sprintf("http://localhost:%s/api/allTableClasses", PORT),
					AllTableClasses,
				},
				{
					"CleanReference",
					fmt.Sprintf("http://localhost:%s/api/reference?cleanRef=true", PORT),
					ReferenceMatrix,
				},
				{
					"IssueOne",
					fmt.Sprintf("http://localhost:%s/api/issueOne", PORT),
					IssueOneMatrix,
				},
				{
					"DataSortValue",
					fmt.Sprintf("http://localhost:%s/api/dataSortValue", PORT),
					DataSortValueMatrix,
				},
				{
					"Issue34",
					fmt.Sprintf("http://localhost:%s/api/issue34", PORT),
					Issue34Matrix,
				},
				{
					"Issue56",
					fmt.Sprintf("http://localhost:%s/api/issue56?cleanRef=true", PORT),
					Issue56Matrix,
				},
				{
					"Issue77",
					fmt.Sprintf("http://localhost:%s/api/issue77", PORT),
					Issue77Matrix,
				},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					var got [][][]string
					execGetRequest(t, tc.url, &got)

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				})
			}
		})

		t.Run("MatrixVerbose", func(t *testing.T) {
			tests := []struct {
				name string
				url  string
				want [][][]client.Verbose
			}{
				{
					"Issue77",
					fmt.Sprintf("http://localhost:%s/api/issue77?verbose=true", PORT),
					Issue77MatrixVerbose,
				},
				{
					"Issue93",
					fmt.Sprintf("http://localhost:%s/api/issue93?verbose=true", PORT),
					Issue93MatrixVerbose,
				},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					var got [][][]client.Verbose
					execGetRequest(t, tc.url, &got)

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				})
			}
		})

		t.Run("KeyValue", func(t *testing.T) {
			tests := []struct {
				name string
				url  string
				want [][]map[string]string
			}{
				{
					"Simple",
					fmt.Sprintf("http://localhost:%s/api/simpleKeyValue?keyRows=1", PORT),
					SimpleKeyValue,
				},
				{
					"SimpleWithTableParameter",
					fmt.Sprintf("http://localhost:%s/api/simpleKeyValue?keyRows=1&table=0", PORT),
					SimpleKeyValue,
				},
				{
					"Complex",
					fmt.Sprintf("http://localhost:%s/api/complexKeyValue?keyRows=2&cleanRef=true", PORT),
					ComplexKeyValue,
				},
				{
					"Issue85",
					fmt.Sprintf("http://localhost:%s/api/issue85?keyRows=1", PORT),
					Issue85KeyValue,
				},
				{
					"MismatchRows",
					fmt.Sprintf("http://localhost:%s/api/keyValueBadRows?keyRows=1", PORT),
					MismatchRowsKeyValue,
				},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					var got [][]map[string]string
					execGetRequest(t, tc.url, &got)

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				})
			}
		})

		t.Run("NoTables", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/noTables", PORT), nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			want := "[]\n"
			if string(b) != want {
				t.Errorf("want %v, got %b", want, b)
			}
		})

		t.Run("UserAgent", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/UserAgent", PORT), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("User-Agent", "test@mail.com")

			_, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("NoUserAgent", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/NoUserAgent", PORT), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("User-Agent", "")

			_, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("ResponseContentType", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/golden", PORT), nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			want := "application/json"
			got := resp.Header.Get("Content-Type")

			if want != got {
				t.Errorf("want %s, got %s", want, got)
			}

		})
	})

	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			name string
			url  string
			want status.Status
		}{
			{
				"GettingData",
				fmt.Sprintf("http://localhost:%s/api/StatusRequestEntityTooLarge", PORT),
				status.Status{
					Message: "StatusRequestEntityTooLarge",
					Code:    http.StatusRequestEntityTooLarge,
					Details: status.Details{
						status.Page: "StatusRequestEntityTooLarge",
					},
				},
			},
			{
				"BadRowSpan",
				fmt.Sprintf("http://localhost:%s/api/badRowSpan?table=0", PORT),
				status.NewStatus("no integer value in span attribute: [x]", http.StatusInternalServerError, status.WithDetails(
					status.Details{
						status.TableIndex:  float64(0),
						status.RowIndex:    float64(1),
						status.ColumnIndex: float64(0),
					},
				)),
			},
			{
				"BadColSpan",
				fmt.Sprintf("http://localhost:%s/api/badColSpan", PORT),
				status.Status{
					Message: "no integer value in span attribute: [x]",
					Code:    http.StatusInternalServerError,
					Details: status.Details{
						status.TableIndex:  float64(0),
						status.RowIndex:    float64(1),
						status.ColumnIndex: float64(1),
					},
				},
			},
			{
				"BadTableParameter",
				fmt.Sprintf("http://localhost:%s/api/badTableParameter?table=x", PORT),
				status.Status{
					Message: `strconv.Atoi: parsing "x": invalid syntax`,
					Code:    http.StatusBadRequest,
				},
			},
			{
				"KeyRows less than one",
				fmt.Sprintf("http://localhost:%s/api/badKeyRows?keyRows=0", PORT),
				status.Status{
					Message: "keyRows must be at least 1",
					Code:    http.StatusBadRequest,
				},
			},
			{
				"Not enough rows",
				fmt.Sprintf("http://localhost:%s/api/keyValueOneRow?keyRows=1", PORT),
				status.Status{
					Message: "table needs at least two rows",
					Code:    http.StatusBadRequest,
					Details: status.Details{
						status.TableIndex: float64(0),
					},
				},
			},
			{
				"KeyRows not a number",
				fmt.Sprintf("http://localhost:%s/api/badKeyRows?keyRows=x", PORT),
				status.Status{
					Message: `strconv.Atoi: parsing "x": invalid syntax`,
					Code:    http.StatusBadRequest,
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				var got status.Status
				execGetRequest(t, tc.url, &got)

				if tc.want.Message != got.Message {
					t.Errorf("expected %v, got %v", tc.want.Message, got.Message)
				}
				if tc.want.Code != got.Code {
					t.Errorf("expected %v, got %v", tc.want.Code, got.Code)
				}
				if !reflect.DeepEqual(tc.want.Details, got.Details) {
					t.Errorf("expected %v, got %v", tc.want.Details, got.Details)
				}
			})
		}

		t.Run("NonGetMethod", func(t *testing.T) {
			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/api/golden", PORT), "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
			}

			want := "Method Not Allowed\n"
			if string(b) != want {
				t.Errorf("expected %v, got %v", want, string(b))
			}
		})
	})
}

func execGetRequest(t *testing.T, url string, v interface{}) {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		t.Fatal(err)
	}
}

func getPageBytes(t *testing.T, page string) []byte {
	t.Helper()

	f, err := os.ReadFile(fmt.Sprintf("testdata/%s.html", page))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func waitforServer() {
	for {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s", PORT))
		if err != nil || resp.StatusCode != http.StatusOK {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return
	}
}
