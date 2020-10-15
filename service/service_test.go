package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
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
		Name                   string
		TablesRequest          *pb.GetTablesRequest
		ExepctedTablesResponse *pb.GetTablesResponse
	}{
		{
			"table1",
			&pb.GetTablesRequest{
				Page: "table1",
				N:    []string{},
			},
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
			"table1_n=0",
			&pb.GetTablesRequest{
				Page: "table1",
				N:    []string{"0"},
			},
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
			"table1_lang=cs",
			&pb.GetTablesRequest{
				Page: "table1",
				N:    []string{},
				Lang: "cs",
			},
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
			"issue1",
			&pb.GetTablesRequest{
				Page: "issue_1",
				N:    []string{},
			},
			&pb.GetTablesResponse{
				Tables: []*pb.Table{
					{
						Rows: map[int64]*pb.Row{
							0: {
								Columns: map[int64]string{
									0: "Language",
									1: "Country",
									2: "Status",
								},
							},
							1: {
								Columns: map[int64]string{
									0: "Korean",
									1: "Australia",
									2: "minority",
								},
							},
							2: {
								Columns: map[int64]string{
									0: "Korean",
									1: "Brazil",
									2: "minority",
								},
							},
							3: {
								Columns: map[int64]string{
									0: "Korean",
									1: "Canada",
									2: "minority",
								},
							},
							4: {
								Columns: map[int64]string{
									0: "Korean",
									1: "China",
									2: "minority, co-official with Chinese in Yanbian Korean Autonomous Prefecture in Jilin Province",
								},
							},
							5: {
								Columns: map[int64]string{
									0: "Korean",
									1: "Japan",
									2: "minority",
								},
							},
							6: {
								Columns: map[int64]string{
									0: "Korean",
									1: "North Korea",
									2: "official",
								},
							},
							7: {
								Columns: map[int64]string{
									0: "Korean",
									1: "Philippines",
									2: "minority",
								},
							},
							8: {
								Columns: map[int64]string{
									0: "Korean",
									1: "South Korea",
									2: "official",
								},
							},
							9: {
								Columns: map[int64]string{
									0: "Korean",
									1: "Taiwan",
									2: "minority",
								},
							},
							10: {
								Columns: map[int64]string{
									0: "Korean",
									1: "United States",
									2: "minority",
								},
							},
							11: {
								Columns: map[int64]string{
									0: "Jeju",
									1: "South Korea",
									2: "official, in Jeju Island",
								},
							},
							12: {
								Columns: map[int64]string{
									0: "Jeju",
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
		t.Run(test.Name, func(t *testing.T) {
			tables, err := svc.GetTables(context.Background(), test.TablesRequest)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tables, test.ExepctedTablesResponse) {
				//t.Log("expected:")
				//print(test.ExepctedTablesResponse.Tables[0])
				//t.Logf("%v", test.ExepctedTablesResponse.Tables)
				//t.Log("got:")
				//print(tables.Tables[0])
				//t.Logf("%v", tables)
				t.Errorf("expected %v, got %v", test.ExepctedTablesResponse.Tables[0], tables.Tables[0])
			}
		})
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
		{
			"httpError",
			[]string{},
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://en.%s/%s", baseURL, "table1"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       getRespBody("table1.html"),
				StatusCode: http.StatusOK,
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://en.%s/%s", baseURL, "issue_1"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       getRespBody("issue_1.html"),
				StatusCode: http.StatusOK,
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cs.%s/%s", baseURL, "table1"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       getRespBody("table1.html"),
				StatusCode: http.StatusOK,
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://en.%s/%s", baseURL, "colspanError"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       getRespBody("colspanError.html"),
				StatusCode: http.StatusOK,
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://en.%s/%s", baseURL, "rowspanError"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       getRespBody("rowspanError.html"),
				StatusCode: http.StatusOK,
			}, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://en.%s/%s", baseURL, "httpError"),
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
				StatusCode: http.StatusRequestEntityTooLarge,
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
