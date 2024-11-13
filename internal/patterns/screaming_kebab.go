package patterns

import (
	"regexp"
	"strings"
)

// Precompiled regular expressions.
var (
	// Matches characters that are not valid in SCREAMING-KEBAB-CASE.
	invalidScreamingKebabCaseCharacters = regexp.MustCompile(`[^A-Z0-9.\-]+`)

	// Matches a valid SCREAMING-KEBAB-CASE string.
	screamingKebabCase = regexp.MustCompile(`^[A-Z0-9.\-]*$`)
)

// Determines if a string is valid SCREAMING-KEBAB-CASE.
func IsScreamingKebabCase(s string) bool {
	return screamingKebabCase.MatchString(s)
}

// Attempts to convert a string to valid SCREAMING-KEBAB-CASE.
func ToScreamingKebabCase(s string) string {
	s = separators.ReplaceAllLiteralString(s, "-")
	s = lowers.ReplaceAllStringFunc(s, func(match string) string {
		return strings.ToUpper(match)
	})
	return invalidScreamingKebabCaseCharacters.ReplaceAllLiteralString(s, "")
}
