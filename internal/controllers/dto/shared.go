package dto

import "strings"

func ParseIncludes(includesString *string) map[string]bool {
	includes := make(map[string]bool)

	if includesString == nil {
		return includes
	}

	includeStr := *includesString

	if includeStr == "" {
		return includes
	}

	for _, include := range strings.Split(includeStr, ",") {
		normalizedField := strings.TrimSpace(strings.ToLower(include))
		includes[normalizedField] = true
	}

	return includes
}
