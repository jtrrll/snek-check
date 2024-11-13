// TODO: Package comment
package patterns

import "regexp"

// Precompiled regular expressions.
var (
	// Matches several lowercase letters.
	lowers = regexp.MustCompile(`[a-z]+`)

	// Matches a single separator.
	separators = regexp.MustCompile(`[_\- ]`)

	// Matches a single space.
	spaces = regexp.MustCompile(`[ ]`)

	// Matches several uppercase letters.
	uppers = regexp.MustCompile(`[A-Z]+`)
)
