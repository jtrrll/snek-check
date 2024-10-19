// TODO: Package comment
package patterns

import "regexp"

// Precompiled regular expressions.
var (
	// Matches several lowercase letters.
	lowers = regexp.MustCompile(`[a-z]+`)

	// Matches several separators.
	separators = regexp.MustCompile(`[_\- ]+`)

	// Matches one or more spaces.
	spaces = regexp.MustCompile(`[ ]+`)

	// Matches several uppercase letters.
	uppers = regexp.MustCompile(`[A-Z]+`)
)
