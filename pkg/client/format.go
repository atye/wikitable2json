package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/atye/wikitable2json/internal/status"
)

type parsed map[int]map[int]cell

type Verbose struct {
	Text  string   `json:"text,omitempty"`
	Links []string `json:"links,omitempty"`
}

var (
	errNotEnoughRows = errors.New("table needs at least two rows")
)

func formatMatrix(data parsed) [][]string {
	matrix := make([][]string, len(data))

	for i := 0; i < len(data); i++ {
		row := data[i]
		matrix[i] = make([]string, len(row))
		for j := 0; j < len(row); j++ {
			matrix[i][j] = row[j].text
		}
	}

	return matrix
}

func formatMatrixVerbose(data parsed) [][]Verbose {
	matrix := make([][]Verbose, len(data))

	for i := 0; i < len(data); i++ {
		row := data[i]
		matrix[i] = make([]Verbose, len(row))
		for j := 0; j < len(row); j++ {
			matrix[i][j].Text = row[j].text
			matrix[i][j].Links = row[j].links
		}
	}

	return matrix
}

func formatKeyValue(data parsed, keyrows int, tableIndex int) ([]map[string]string, error) {
	if len(data) > 1 && keyrows >= 1 {
		keys, err := generateKeys(data, keyrows)
		if err != nil {
			return nil, err
		}

		var kv []map[string]string
		for i := keyrows; i < len(data); i++ {
			pairs := make(map[string]string)
			for j := 0; j < len(data[i]); j++ {
				key := fmt.Sprintf("null%d", j)
				if j < len(keys) {
					key = keys[j]
				}
				pairs[key] = data[i][j].text
			}
			kv = append(kv, pairs)
		}
		return kv, nil
	}
	return nil, status.NewStatus(errNotEnoughRows.Error(), http.StatusBadRequest, status.WithDetails(status.Details{
		status.TableIndex: tableIndex,
	}))
}

func formatKeyValueVerbose(data parsed, keyrows int, tableIndex int) ([]map[string]Verbose, error) {
	if len(data) > 1 && keyrows >= 1 {
		keys, err := generateKeys(data, keyrows)
		if err != nil {
			return nil, err
		}

		var kv []map[string]Verbose
		for i := keyrows; i < len(data); i++ {
			pairs := make(map[string]Verbose)
			for j := 0; j < len(data[i]); j++ {
				key := fmt.Sprintf("null%d", j)
				if j < len(keys) {
					key = keys[j]
				}
				pairs[key] = Verbose{
					Text:  data[i][j].text,
					Links: data[i][j].links,
				}
			}
			kv = append(kv, pairs)
		}
		return kv, nil
	}
	return nil, status.NewStatus(errNotEnoughRows.Error(), http.StatusBadRequest, status.WithDetails(status.Details{
		status.TableIndex: tableIndex,
	}))
}

func generateKeys(data parsed, keyrows int) ([]string, error) {
	var keys []string
	for colNum := 0; colNum < len(data[0]); colNum++ {
		var b strings.Builder
		_, err := b.WriteString(data[0][colNum].text)
		if err != nil {
			return nil, err
		}

		for k := 1; k < keyrows; k++ {
			v := data[k][colNum].text
			if v != data[k-1][colNum].text && v != "" {
				_, err := b.WriteString(fmt.Sprintf(" %s", v))
				if err != nil {
					return nil, err
				}
			}
		}
		keys = append(keys, b.String())
	}
	return keys, nil
}
