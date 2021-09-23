package json

import (
	"encoding/json"
)

func MarshalMap(key string, val string) ([]byte,error) {
	jsonBytes, err := json.MarshalIndent(map[string]interface{}{
		"apiVersion": "1.0",
		key:  val,
	}, "", "  ")
	if err != nil {
		return nil,err
	}
	return jsonBytes,err
}
