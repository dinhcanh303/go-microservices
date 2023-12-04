package utils

import (
	"encoding/json"
)

func JsonRawMessageToArrayString(rawJson json.RawMessage) ([]string, error) {
	var strings []string
	err := json.Unmarshal([]byte(rawJson), &strings)
	if err != nil {
		return nil, err
	}
	return strings, nil
}
