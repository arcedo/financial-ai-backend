package utils

import "strings"

func SanitizeString(str string) string {
	// Remove leading and trailing whitespace
	str = strings.TrimSpace(str)

	// Convert to lowercase
	str = strings.ToLower(str)

	// Replace spaces with underscores
	str = strings.ReplaceAll(str, " ", "_")

	return str
}
