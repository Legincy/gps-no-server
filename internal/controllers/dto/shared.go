package dto

import "strings"

func ParseIncludes(includesString string) map[string]bool {
	includes := make(map[string]bool)

	if includesString == "" {
		return includes
	}

	for _, include := range strings.Split(includesString, ",") {
		normalizedField := strings.TrimSpace(strings.ToLower(include))
		includes[normalizedField] = true
	}

	return includes
}
