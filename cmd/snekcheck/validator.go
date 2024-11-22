package main

import (
	"snekcheck/internal/patterns"
	"strings"
)

// Determines if a filename is valid according to snekcheck's opinion.
func IsValid(name string) bool {
	return len(name) > 0 &&
		patterns.IsPosixFileName(name) &&
		(patterns.IsSnakeCase(name) || isAlmostScreamingSnakeCase(name))
}

// Determines if a filename is SCREAMING_SNAKE_CASE with a snake_case file extension.
func isAlmostScreamingSnakeCase(name string) bool {
	lastIndex := strings.LastIndex(name, ".")
	if lastIndex == -1 {
		return patterns.IsScreamingSnakeCase(name)
	}
	return patterns.IsScreamingSnakeCase(name[:lastIndex]) && patterns.IsSnakeCase(name[lastIndex:])
}
