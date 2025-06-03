package utils

import (
	"encoding/json"
)

// ToJSON converts a value to its JSON string representation
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
} 