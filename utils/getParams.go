package utils

import (
	"fmt"
	"strings"
)

// Example: For "/user/123/profile", GetPathParam("/user/123/profile", 1) returns "123".
func GetPathParam(path string, position int) (string, error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")

	if position < 0 || position >= len(segments) {
		return "", fmt.Errorf("no path segment at position %d", position)
	}

	return segments[position], nil
}
