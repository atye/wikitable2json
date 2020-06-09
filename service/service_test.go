package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/atye/wikitable-api/service/pb"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	startMocks()
	defer httpmock.DeactivateAndReset()
	os.Exit(m.Run())
}

func Test_GetTables(t *testing.T) {
	tests := []struct {
		Page                   string
		N                      []string
		ExepctedTablesResponse *pb.GetTablesResponse
	}{
		{
			"table1",
			[]string{},
			&pb.GetTablesResponse{
				Tables: []*pb.Table{
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
		},
		{
			"table1",
			[]string{"0"},
			&pb.GetTablesResponse{
				Tables: []*pb.Table{
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
		},
	}

	svc := &Service{}

	for _, test := range tests {
		gtReq := &pb.GetTablesRequest{
			Page: test.Page,
			N:    test.N,
		}

		tables, err := svc.GetTables(context.Background(), gtReq)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(tables, test.ExepctedTablesResponse) {
			t.Log("expected:")
			print(test.ExepctedTablesResponse.Tables[0])
			t.Log("got:")
			print(tables.Tables[0])
			t.Errorf("expected %v, got %v", test.ExepctedTablesResponse.Tables[0], tables.Tables[0])
		}
	}
}

func Test_GetTables_Error(t *testing.T) {
	tests := []struct {
		Page string
		N    []string
	}{
		{
			"rowspanError",
			[]string{},
		},
		{
			"colspanError",
			[]string{},
		},
		{
			"rowspanError",
			[]string{"0"},
		},
		{
			"table1",
			[]string{"x"},
		},
	}

	svc := &Service{}

	for _, test := range tests {
		gtReq := &pb.GetTablesRequest{
			Page: test.Page,
			N:    test.N,
		}

		_, err := svc.GetTables(context.Background(), gtReq)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}

func startMocks() {
	httpmock.Activate()

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", baseURL, "table1"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: getRespBody("table1.html"),
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", baseURL, "colspanError"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: getRespBody("colspanError.html"),
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", baseURL, "rowspanError"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: getRespBody("rowspanError.html"),
			}, nil
		})
}

func getRespBody(file string) io.ReadCloser {
	tables, err := os.Open(fmt.Sprintf("%s/%s", "testdata", file))
	if err != nil {
		panic(err)
	}

	return tables
}

func print(table *pb.Table) {
	for i := 0; i < len(table.Rows); i++ {
		fmt.Printf("row %d:\n", i)
		for j := 0; j < len(table.Rows[int64(i)].Columns); j++ {
			fmt.Printf("column %d: %#v\n", j, table.Rows[int64(i)].Columns[int64(j)])
		}
	}
}
