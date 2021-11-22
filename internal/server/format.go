package server

import (
	"net/http"
	"sync"

	"github.com/atye/wikitable-api/internal/status"
)

type verbose map[int]map[int]string

func toFormat(format string, v verbose, tableIndex int) (interface{}, error) {
	switch format {
	case "keyvalue":
		fallthrough
	case "keyValue":
		kv, err := toKeyValue(v, tableIndex)
		if err != nil {
			return nil, err
		}
		return kv, nil
	default:
		return toMatrix(v), nil
	}
}

type Matrix [][]string

func toMatrix(vf verbose) interface{} {
	matrix := make(Matrix, len(vf))

	var wg sync.WaitGroup
	for i := 0; i < len(vf); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			row := vf[i]
			matrix[i] = make([]string, len(row))
			for j := 0; j < len(row); j++ {
				matrix[i][j] = row[j]
			}
		}(i)
	}
	wg.Wait()

	return matrix
}

type KeyValue []map[string]string

func toKeyValue(rows verbose, tableIndex int) (interface{}, error) {
	if len(rows) > 0 {
		var headers []string
		for i := 0; i < len(rows[0]); i++ {
			headers = append(headers, rows[0][i])
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
					pairs[headers[j]] = rows[i][j]
				}

				kv = append(kv, pairs)
			}
			return kv, nil
		}
		msg := "table only seems to have one row, need at least two"
		return nil, status.NewStatus(msg, http.StatusInternalServerError, status.WithDetails(status.Details{
			status.TableIndex: tableIndex,
		}))
	}
	return KeyValue{}, nil
}
