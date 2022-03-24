package entrypoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/atye/wikitable-api/internal/server/data"
	"github.com/atye/wikitable-api/internal/status"
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
		case "/api/rest_v1/page/html/reference":
			w.Write(getPageBytes(t, "reference"))
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
			want := "github.com/atye/wikitable-api"
			if want != got {
				t.Errorf("want %s, got %s", want, got)
			}
		default:
			t.Fatalf("path %s not supported", r.URL.Path)
		}
	}))
	defer ts.Close()

	go Run(Config{
		Port:    PORT,
		WikiAPI: data.NewWikiClient(ts.URL),
	})

	waitforServer()

	t.Run("Success", func(t *testing.T) {
		t.Run("Matrix", func(t *testing.T) {
			tests := []struct {
				page string
				want interface{}
			}{
				{
					"golden",
					GoldenMatrix,
				},
				{
					"issueOne",
					IssueOneMatrix,
				},
				{
					"dataSortValue",
					DataSortValueMatrix,
				},
				{
					"issue34",
					Issue34Matrix,
				},
			}

			for _, tc := range tests {
				t.Run(tc.page, func(t *testing.T) {
					addr := fmt.Sprintf("http://localhost:%s/api/%s", PORT, tc.page)
					var got [][][]string
					execGetRequest(t, addr, &got)

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				})
			}
		})

		t.Run("KeyValue", func(t *testing.T) {
			tests := []struct {
				page string
				want interface{}
			}{
				{
					"golden",
					GoldenKeyValue,
				},
			}

			for _, tc := range tests {
				t.Run(tc.page, func(t *testing.T) {
					addr := fmt.Sprintf("http://localhost:%s/api/%s?format=keyvalue", PORT, tc.page)
					var got [][]map[string]string
					execGetRequest(t, addr, &got)

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				})
			}
		})

		t.Run("WithParameters", func(t *testing.T) {
			addr := fmt.Sprintf("http://localhost:%s/api/golden?lang=sp&format=matrix&table=0", PORT)
			want := GoldenMatrix

			var got [][][]string
			execGetRequest(t, addr, &got)

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v\n got %v", want, got)
			}
		})

		t.Run("AllTableClasses", func(t *testing.T) {
			addr := fmt.Sprintf("http://localhost:%s/api/allTableClasses", PORT)
			want := [][][]string{
				GoldenMatrix[0],
				GoldenMatrix[0],
				GoldenMatrix[0],
			}

			var got [][][]string
			execGetRequest(t, addr, &got)

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v\n got %v", want, got)
			}
		})

		t.Run("CleanReference", func(t *testing.T) {
			addr := fmt.Sprintf("http://localhost:%s/api/reference?cleanRef=true", PORT)
			want := [][][]string{
				ReferenceMatrix[0],
			}

			var got [][][]string
			execGetRequest(t, addr, &got)

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v\n got %v", want, got)
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
			addr := fmt.Sprintf("http://localhost:%s/api/NoUserAgent", PORT)
			req, err := http.NewRequest(http.MethodGet, addr, nil)
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
		type test struct {
			name string
			url  string
			want status.Status
		}

		tests := []test{
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
						status.TableIndex:   float64(0),
						status.RowNumber:    float64(1),
						status.ColumnNumber: float64(0),
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
						status.TableIndex:   float64(0),
						status.RowNumber:    float64(1),
						status.ColumnNumber: float64(1),
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

			var got status.Status
			err = json.NewDecoder(resp.Body).Decode(&got)
			if err != nil {
				t.Fatal(err)
			}

			want := status.NewStatus(fmt.Sprintf("method %s not allowed", http.MethodPost), http.StatusMethodNotAllowed)
			if !reflect.DeepEqual(want, got) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})

		t.Run("KeyValue", func(t *testing.T) {
			t.Run("MismatchedRow", func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/issueOne?format=keyValue", PORT))
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				var got status.Status
				err = json.NewDecoder(resp.Body).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}

				want := status.NewStatus("keys length does not match row length", http.StatusInternalServerError, status.WithDetails(status.Details{
					status.TableIndex: float64(0),
					status.RowNumber:  float64(1),
					status.KeysLength: float64(3),
					status.RowLength:  float64(1),
				}))
				if !reflect.DeepEqual(want, got) {
					t.Errorf("expected %v, got %v", want, got)
				}
			})

			t.Run("OneRow", func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/dataSortValue?format=keyValue&table=0", PORT))
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				var got status.Status
				err = json.NewDecoder(resp.Body).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}

				want := status.NewStatus("table only seems to have one row, need at least two", http.StatusInternalServerError, status.WithDetails(status.Details{
					status.TableIndex: float64(0),
				}))
				if !reflect.DeepEqual(want, got) {
					t.Errorf("expected %v, got %v", want, got)
				}
			})
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
