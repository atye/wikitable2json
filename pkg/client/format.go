package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/atye/wikitable2json/internal/status"
)

type verbose map[int]map[int]cell

var (
	errNotEnoughRows = errors.New("table needs at least two rows")
)

func formatMatrix(data verbose) [][]string {
	matrix := make([][]string, len(data))

	for i := 0; i < len(data); i++ {
		row := data[i]
		matrix[i] = make([]string, len(row))
		for j := 0; j < len(row); j++ {
			matrix[i][j] = row[j].value
		}
	}

	return matrix
}

func formatKeyValue(data verbose, keyrows int, tableIndex int) ([]map[string]string, error) {
	if len(data) > 1 && keyrows >= 1 {
		var keys []string
		for colNum := 0; colNum < len(data[0]); colNum++ {
			var b strings.Builder
			_, err := b.WriteString(data[0][colNum].value)
			if err != nil {
				return nil, err
			}

			for k := 1; k < keyrows; k++ {
				v := data[k][colNum].value
				if v != data[k-1][colNum].value && v != "" {
					_, err := b.WriteString(fmt.Sprintf(" %s", v))
					if err != nil {
						return nil, err
					}
				}
			}
			keys = append(keys, b.String())
		}
		var kv []map[string]string
		for i := keyrows; i < len(data); i++ {
			pairs := make(map[string]string)
			for j := 0; j < len(data[i]); j++ {
				key := fmt.Sprintf("null%d", j)
				if j < len(keys) {
					key = keys[j]
				}
				pairs[key] = data[i][j].value
			}
			kv = append(kv, pairs)
		}
		return kv, nil
	}
	return nil, status.NewStatus(errNotEnoughRows.Error(), http.StatusBadRequest, status.WithDetails(status.Details{
		status.TableIndex: tableIndex,
	}))
}
