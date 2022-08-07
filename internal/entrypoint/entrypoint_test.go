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

	"github.com/atye/wikitable2json/internal/api"
	"github.com/atye/wikitable2json/internal/server/status"
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
		Client: client.NewTableGetter(""),
	})

	waitforServer()

	t.Run("Success", func(t *testing.T) {
		t.Run("Matrix", func(t *testing.T) {
			tests := []struct {
				page  string
				query string
				want  interface{}
			}{
				{
					"golden",
					"",
					GoldenMatrix,
				},
				{
					"issueOne",
					"",
					IssueOneMatrix,
				},
				{
					"dataSortValue",
					"",
					DataSortValueMatrix,
				},
				{
					"issue34",
					"",
					Issue34Matrix,
				},
				{
					"complexKeyValue",
					"?cleanRef=true",
					ComplexMatrix,
				},
			}

			for _, tc := range tests {
				t.Run(tc.page, func(t *testing.T) {
					addr := fmt.Sprintf("http://localhost:%s/api/%s%s", PORT, tc.page, tc.query)
					var got [][][]string
					execGetRequest(t, addr, &got)

					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want %v\n got %v", tc.want, got)
					}
				})
			}
		})

		t.Run("KeyValue", func(t *testing.T) {
			t.Run("Simple", func(t *testing.T) {
				addr := fmt.Sprintf("http://localhost:%s/api/simpleKeyValue?keyRows=1", PORT)

				want := []client.KeyValue{
					{
						{
							"Rank":    "1",
							"Account": "Alpha",
						},
					},
				}

				var got []client.KeyValue
				execGetRequest(t, addr, &got)

				if !reflect.DeepEqual(want, got) {
					t.Errorf("want %v\n got %v", want, got)
				}

				// do it again with table param for coverage
				addr = fmt.Sprintf("http://localhost:%s/api/simpleKeyValue?keyRows=1&table=0", PORT)

				want = []client.KeyValue{
					{
						{
							"Rank":    "1",
							"Account": "Alpha",
						},
					},
				}

				var resp []client.KeyValue
				execGetRequest(t, addr, &resp)

				if !reflect.DeepEqual(want, resp) {
					t.Errorf("want %v\n got %v", want, resp)
				}
			})

			t.Run("Complex", func(t *testing.T) {
				addr := fmt.Sprintf("http://localhost:%s/api/complexKeyValue?keyRows=2&cleanRef=true", PORT)

				want := []client.KeyValue{
					{
						{
							"Date":              "18–24 April 2022",
							"Brand":             "Roy Morgan",
							"Interview mode":    "Telephone/online",
							"Sample size":       "1393",
							"Primary vote L/NP": "35.5%",
							"Primary vote ALP":  "35%",
							"Primary vote GRN":  "12%",
							"Primary vote ONP":  "4.5%",
							"Primary vote UAP":  "1.5%",
							"Primary vote OTH":  "11.5%",
							"UND":               "–",
							"2pp vote L/NP":     "45.5%",
							"2pp vote ALP":      "54.5%",
						},
						{
							"Date":              "20–23 April 2022",
							"Brand":             "Newspoll-YouGov",
							"Interview mode":    "Online",
							"Sample size":       "1538",
							"Primary vote L/NP": "36%",
							"Primary vote ALP":  "37%",
							"Primary vote GRN":  "11%",
							"Primary vote ONP":  "3%",
							"Primary vote UAP":  "4%",
							"Primary vote OTH":  "9%",
							"UND":               "–",
							"2pp vote L/NP":     "47%",
							"2pp vote ALP":      "53%",
						},
						{
							"Date":              "20–23 April 2022",
							"Brand":             "Ipsos",
							"Interview mode":    "Telephone/online",
							"Sample size":       "2302",
							"Primary vote L/NP": "32%",
							"Primary vote ALP":  "34%",
							"Primary vote GRN":  "12%",
							"Primary vote ONP":  "4%",
							"Primary vote UAP":  "3%",
							"Primary vote OTH":  "8%",
							"UND":               "8%",
							"2pp vote L/NP":     "45%",
							"2pp vote ALP":      "55%",
						},
					},
				}

				var got []client.KeyValue
				execGetRequest(t, addr, &got)

				if !reflect.DeepEqual(want, got) {
					t.Errorf("want %v\n got %v", want, got)
				}
			})
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

		t.Run("CleanReferenceCitationNeeded", func(t *testing.T) {
			addr := fmt.Sprintf("http://localhost:%s/api/issue56?cleanRef=true", PORT)
			want := Issue56Matrix

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
			{
				"KeyRows less than one",
				fmt.Sprintf("http://localhost:%s/api/badKeyRows?keyRows=0", PORT),
				status.Status{
					Message: "keyRows must be at least 1",
					Code:    http.StatusBadRequest,
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
				resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/keyValueBadRows?keyRows=1", PORT))
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				var got status.Status
				err = json.NewDecoder(resp.Body).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}

				want := status.NewStatus("number of keys does not equal number of values", http.StatusInternalServerError, status.WithDetails(status.Details{
					status.TableIndex: float64(0),
					status.RowNumber:  float64(2),
					status.KeysLength: float64(2),
					status.RowLength:  float64(3),
				}))
				if !reflect.DeepEqual(want, got) {
					t.Errorf("expected %v, got %v", want, got)
				}
			})

			t.Run("OneRow", func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/keyValueOneRow?keyRows=1", PORT))
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				var got status.Status
				err = json.NewDecoder(resp.Body).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}

				want := status.NewStatus("table needs at least two rows", http.StatusBadRequest, status.WithDetails(status.Details{
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
