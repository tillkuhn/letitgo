package json

import (
	"encoding/json"
	"strings"
)

// MarshalMapNicely serializes a Map into nice (i.e. with multiple lines + ident) JSON
func MarshalMapNicely(key string, val string) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(map[string]interface{}{
		"apiVersion": "1.0",
		key:          val,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, err
}

// GenericUnmarshal play around with generics
func GenericUnmarshal[J any](jsonStr string) J {
	var input2 J
	_ = json.NewDecoder(strings.NewReader(jsonStr)).Decode(&input2)
	return input2
}
