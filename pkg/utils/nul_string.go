package utils

import (
	"log/slog"
	"strings"
)

func HandleNullString(value interface{}) string {
	if value != "" {
		return value.(string)
	}
	return "" // or any default value you prefer
}

func HandleNullStringSlice(value interface{}) []string {
	slog.Info("Value::", value)
	if value != "{}" {
		return []string{} //extractStrings(value)
	}
	return []string{} // or any default value you prefer
}
func extractStrings(input interface{}) []string {
	// Remove the surrounding curly braces
	input = strings.Trim(input.(string), "{}")
	// Split the string using commas
	parts := strings.Split(input.(string), ",")
	// Trim any leading or trailing spaces from each part
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}
