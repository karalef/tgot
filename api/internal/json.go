package internal

import "encoding/json"

// MergeJSON merges JSON objects.
func MergeJSON(v ...interface{}) ([]byte, error) {
	var data []byte
	for i := range v {
		d, err := json.Marshal(v[i])
		if err != nil {
			return nil, err
		}
		if string(d) == "{}" {
			continue
		}
		if len(data) == 0 {
			data = append(data, d...)
		} else {
			data[len(data)-1] = ','
			data = append(data, d[1:]...)
		}
	}
	return data, nil
}
