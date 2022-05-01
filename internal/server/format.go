package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/atye/wikitable2json/internal/status"
)

type verbose map[int]map[int]cell

var (
	ErrNotEnoughRows         = errors.New("table needs at least two rows")
	ErrNumKeysValuesMismatch = errors.New("number of keys does not equal number of values")
)

type Matrix [][]string

func formatMatrix(data verbose) Matrix {
	matrix := make(Matrix, len(data))

	for i := 0; i < len(data); i++ {
		row := data[i]
		matrix[i] = make([]string, len(row))
		for j := 0; j < len(row); j++ {
			matrix[i][j] = row[j].value
		}
	}

	return matrix
}

type KeyValue []map[string]string

func formatKeyValue(data verbose, keyrows int, tableIndex int) (KeyValue, error) {
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
		var kv KeyValue
		for i := keyrows; i < len(data); i++ {
			if len(keys) != len(data[i]) {
				return nil, status.NewStatus(ErrNumKeysValuesMismatch.Error(), http.StatusInternalServerError, status.WithDetails(status.Details{
					status.TableIndex: tableIndex,
					status.RowNumber:  i,
					status.KeysLength: len(keys),
					status.RowLength:  len(data[i]),
				}))
			}

			pairs := make(map[string]string)
			for j := 0; j < len(data[i]); j++ {
				pairs[keys[j]] = data[i][j].value
			}
			kv = append(kv, pairs)
		}
		return kv, nil
	}
	return nil, status.NewStatus(ErrNotEnoughRows.Error(), http.StatusBadRequest, status.WithDetails(status.Details{
		status.TableIndex: tableIndex,
	}))
}
