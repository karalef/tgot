package tg

import "encoding/json"

func mergeJSON(v ...interface{}) ([]byte, error) {
	var data []byte
	for i := range v {
		d, err := json.Marshal(v[i])
		if err != nil {
			return nil, err
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
