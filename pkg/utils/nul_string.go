package utils

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func HandleNullString(value interface{}) string {
	if value != "" {
		return value.(string)
	}
	return "" // or any default value you prefer
}

func HandleInterfaceToArrayString(value interface{}) []string {
	switch v := value.(type) {
	case []string:
		return v
	case pq.StringArray:
		return []string(v)
	case []uint8:
		fmt.Printf("Type: []uint8\n")
		result := convertStringToArray(string(v))
		// Assuming it's binary data, you might want to convert it to a string
		return result
	case nil:
		return nil
	default:
		return []string{}
	}
}
func convertStringToArray(str string) []string {
	if str == "{}" {
		return []string{}
	}
	str = strings.Trim(str, "{}")
	parts := strings.Split(str, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}
