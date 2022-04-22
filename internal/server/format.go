package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/atye/wikitable-api/internal/status"
)

type verbose map[int]map[int]cell

func format(format string, v verbose, tableIndex int) (interface{}, error) {
	switch format {
	case "keyvalue":
		fallthrough
	case "keyValue":
		kv, err := keyValue(v, tableIndex)
		if err != nil {
			return nil, err
		}
		return kv, nil
	default:
		return matrix(v), nil
	}
}

type Matrix [][]string

func matrix(vf verbose) interface{} {
	matrix := make(Matrix, len(vf))

	var wg sync.WaitGroup
	for i := 0; i < len(vf); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			row := vf[i]
			matrix[i] = make([]string, len(row))
			for j := 0; j < len(row); j++ {
				matrix[i][j] = row[j].value
			}
		}(i)
	}
	wg.Wait()

	return matrix
}

type KeyValue []map[string]string

func keyValue(data verbose, tableIndex int) (interface{}, error) {
	headRows := 0
	for _, col := range data {
		if len(col) > 0 {
			if col[0].header {
				headRows++
			}
		}
	}

	if headRows > 0 {
		if len(data) > 1 {
			var keys []string
			for colNum := 0; colNum < len(data[0]); colNum++ {
				var b strings.Builder
				_, err := b.WriteString(data[0][colNum].value)
				if err != nil {
					return nil, err
				}

				for k := 1; k < headRows; k++ {
					v := data[k][colNum].value
					if v != data[k-1][colNum].value {
						_, err := b.WriteString(fmt.Sprintf(" %s", v))
						if err != nil {
							return nil, err
						}
					}
				}
				keys = append(keys, b.String())
			}
			var kv KeyValue
			for i := headRows; i < len(data); i++ {
				if len(keys) != len(data[i]) {
					msg := "keys length does not match row length"
					return nil, status.NewStatus(msg, http.StatusInternalServerError, status.WithDetails(status.Details{
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
		msg := "table needs at least two rows"
		return nil, status.NewStatus(msg, http.StatusBadRequest, status.WithDetails(status.Details{
			status.TableIndex: tableIndex,
		}))
	}

	/*if len(rows) > 0 {
		var headers []string
		for i := 0; i < len(rows[0]); i++ {
			headers = append(headers, rows[0][i].value)
		}

		if len(rows) > 1 {
			var kv KeyValue
			for i := 1; i <= len(rows)-1; i++ {
				if len(headers) != len(rows[i]) {
					msg := "keys length does not match row length"
					return nil, status.NewStatus(msg, http.StatusInternalServerError, status.WithDetails(status.Details{
						status.TableIndex: tableIndex,
						status.RowNumber:  i,
						status.KeysLength: len(headers),
						status.RowLength:  len(rows[i]),
					}))
				}

				pairs := make(map[string]string)
				for j := 0; j < len(rows[i]); j++ {
					pairs[headers[j]] = rows[i][j].value
				}

				kv = append(kv, pairs)
			}
			return kv, nil
		}
		msg := "table only seems to have one row, need at least two"
		return nil, status.NewStatus(msg, http.StatusInternalServerError, status.WithDetails(status.Details{
			status.TableIndex: tableIndex,
		}))
	}*/
	return KeyValue{}, nil
}
