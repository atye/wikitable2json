package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/atye/wikitable-api/internal/service/pb"
	"google.golang.org/grpc"
)

func TestRunSuccess(t *testing.T) {
	tests := []struct {
		Name           string
		Page           string
		Tables         []int32
		Lang           string
		Config         Config
		ExpectedTables []*pb.Table
	}{
		{
			"table1",
			"table1",
			[]int32{},
			"",
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("table1.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			[]*pb.Table{
				{
					Caption: "test",
					Rows: map[int64]*pb.Row{
						0: {
							Columns: map[int64]string{
								0: "Column 1",
								1: "Column 2",
								2: "Column 3",
							},
						},
						1: {
							Columns: map[int64]string{
								0: "A",
								1: "B",
								2: "B",
							},
						},
						2: {
							Columns: map[int64]string{
								0: "A",
								1: "C",
								2: "D",
							},
						},
						3: {
							Columns: map[int64]string{
								0: "E",
								1: "F",
								2: "F",
							},
						},
						4: {
							Columns: map[int64]string{
								0: "G",
								1: "F",
								2: "F",
							},
						},
						5: {
							Columns: map[int64]string{
								0: "H",
								1: "H",
								2: "H",
							},
						},
					},
				},
			},
		},
		{
			"table1_n=0",
			"table1",
			[]int32{0},
			"",
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("table1.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			[]*pb.Table{
				{
					Caption: "test",
					Rows: map[int64]*pb.Row{
						0: {
							Columns: map[int64]string{
								0: "Column 1",
								1: "Column 2",
								2: "Column 3",
							},
						},
						1: {
							Columns: map[int64]string{
								0: "A",
								1: "B",
								2: "B",
							},
						},
						2: {
							Columns: map[int64]string{
								0: "A",
								1: "C",
								2: "D",
							},
						},
						3: {
							Columns: map[int64]string{
								0: "E",
								1: "F",
								2: "F",
							},
						},
						4: {
							Columns: map[int64]string{
								0: "G",
								1: "F",
								2: "F",
							},
						},
						5: {
							Columns: map[int64]string{
								0: "H",
								1: "H",
								2: "H",
							},
						},
					},
				},
			},
		},
		{
			"table1_lang=cs",
			"table1",
			[]int32{},
			"cs",
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("table1.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			[]*pb.Table{
				{
					Caption: "test",
					Rows: map[int64]*pb.Row{
						0: {
							Columns: map[int64]string{
								0: "Column 1",
								1: "Column 2",
								2: "Column 3",
							},
						},
						1: {
							Columns: map[int64]string{
								0: "A",
								1: "B",
								2: "B",
							},
						},
						2: {
							Columns: map[int64]string{
								0: "A",
								1: "C",
								2: "D",
							},
						},
						3: {
							Columns: map[int64]string{
								0: "E",
								1: "F",
								2: "F",
							},
						},
						4: {
							Columns: map[int64]string{
								0: "G",
								1: "F",
								2: "F",
							},
						},
						5: {
							Columns: map[int64]string{
								0: "H",
								1: "H",
								2: "H",
							},
						},
					},
				},
			},
		},
		{
			"issue1",
			"issue1",
			[]int32{},
			"",
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("issue_1.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			[]*pb.Table{
				{
					Rows: map[int64]*pb.Row{
						0: {
							Columns: map[int64]string{
								0: "Jeju",
								1: "South Korea",
								2: "official, in Jeju Island",
							},
						},
						1: {
							Columns: map[int64]string{
								0: "Jeju",
							},
						},
					},
				},
			},
		},
		{
			"data-sort-value_in_span",
			"data-sort-value",
			[]int32{},
			"",
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("data-sort-value.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			[]*pb.Table{
				{
					Rows: map[int64]*pb.Row{
						0: {
							Columns: map[int64]string{
								0: "Abu Dhabi, United Arab Emirates",
								1: "N/A",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			go func() {
				Run(context.Background(), tc.Config)
			}()

			defer func() {
				tc.Config.GrpcSvr.GracefulStop()
				tc.Config.HTTPSvr.Shutdown(context.Background())
			}()

			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/api/v2/%s", tc.Page), nil)
			if err != nil {
				t.Fatal(err)
			}

			queryParams := req.URL.Query()
			for _, table := range tc.Tables {
				queryParams.Add("table", strconv.Itoa(int(table)))
			}

			if tc.Lang != "" {
				queryParams.Add("lang", tc.Lang)
			}

			req.URL.RawQuery = queryParams.Encode()

			<-tc.Config.signalReady
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			var tablesResp pb.TablesResponse

			err = json.NewDecoder(resp.Body).Decode(&tablesResp)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tablesResp.Tables, tc.ExpectedTables) {
				t.Errorf("expected %v, got %v", tc.ExpectedTables, tablesResp.Tables)
			}
		})
	}
}

func TestRunError(t *testing.T) {
	tests := []struct {
		Name               string
		Page               string
		Tables             []int32
		Config             Config
		ExpectedStatusCode int
	}{
		{
			"RowSpanError",
			"rowSpanError",
			[]int32{},
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("rowspanError.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			http.StatusInternalServerError,
		},
		{
			"ColSpanError",
			"colSpanError",
			[]int32{},
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("colspanError.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			http.StatusInternalServerError,
		},
		{
			"SpanError_n=0",
			"SpanError_n=0",
			[]int32{0},
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       getRespBody("colspanError.html"),
						StatusCode: http.StatusOK,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			http.StatusInternalServerError,
		},
		{
			"WikiAPINotOkError",
			"",
			[]int32{0},
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(&bytes.Buffer{}),
						StatusCode: http.StatusRequestEntityTooLarge,
					}, nil
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			http.StatusRequestEntityTooLarge,
		},
		{
			"HTTPGetError",
			"test",
			[]int32{0},
			Config{
				HTTPGet: func(string) (*http.Response, error) {
					return nil, errors.New("error")
				},
				HTTPSvr: &http.Server{
					Addr: fmt.Sprintf(":%s", "8080"),
				},
				GrpcSvr:     grpc.NewServer(),
				signalReady: make(chan struct{}),
			},
			http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			go func() {
				Run(context.Background(), tc.Config)
			}()

			defer func() {
				tc.Config.GrpcSvr.GracefulStop()
				tc.Config.HTTPSvr.Shutdown(context.Background())
			}()

			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/api/v2/%s", tc.Page), nil)
			if err != nil {
				t.Fatal(err)
			}

			queryParams := req.URL.Query()
			for _, table := range tc.Tables {
				queryParams.Add("table", strconv.Itoa(int(table)))
			}

			req.URL.RawQuery = queryParams.Encode()

			<-tc.Config.signalReady
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.ExpectedStatusCode {
				t.Errorf("expected status %d, got %d", tc.ExpectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func getRespBody(file string) io.ReadCloser {
	tables, err := os.Open(fmt.Sprintf("%s/%s", "testdata", file))
	if err != nil {
		panic(err)
	}

	return tables
}
