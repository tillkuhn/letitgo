package json

import (
	"encoding/json"
	"strings"
)

func MarshalMap(key string, val string) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(map[string]interface{}{
		"apiVersion": "1.0",
		key:          val,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, err
}

func GenericUnmarshal[J any](jsonStr string) J {
	var input2 J
	_ = json.NewDecoder(strings.NewReader(jsonStr)).Decode(&input2)
	return input2
}
