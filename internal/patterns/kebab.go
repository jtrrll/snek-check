package patterns

import (
	"regexp"
	"strings"
)

// Precompiled regular expressions.
var (
	// Matches characters that are not valid in `kebab-case`.
	invalidKebabCaseCharacters = regexp.MustCompile(`[^a-z0-9.\-]+`)

	// Matches a valid `kebab-case` string.
	kebabCase = regexp.MustCompile(`^[a-z0-9.\-]+$`)
)

// Determines if a string is valid `kebab-case`.
func IsKebabCase(s string) bool {
	return kebabCase.MatchString(s)
}

// Attempts to convert a string to valid `kebab-case`.
func ToKebabCase(s string) string {
	s = separators.ReplaceAllLiteralString(s, "-")
	s = uppers.ReplaceAllStringFunc(s, func(match string) string {
		return strings.ToLower(match)
	})
	return invalidKebabCaseCharacters.ReplaceAllLiteralString(s, "")
}
