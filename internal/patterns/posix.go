package patterns

import (
	"regexp"
)

// Precompiled regular expressions.
var (
	// Matches characters that are not valid in POSIX filenames according to
	// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_276.
	invalidPosixFileNameCharacters = regexp.MustCompile(`[^a-zA-Z0-9._\-]*`)

	// Matches a valid POSIX filename according to
	// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_276.
	// The length limitation is enforced by file systems rather than POSIX,
	// but many file systems have a similar limit.
	validPosixFileName = regexp.MustCompile(`^[a-zA-Z0-9._][a-zA-Z0-9._\-]{0,254}$`)
)

// Determines if a string is a valid POSIX filename.
func IsPosixFileName(s string) bool {
	return validPosixFileName.MatchString(s)
}

// Attempts to convert a string to a valid POSIX filename.
func ToPosixFileName(s string) string {
	s = spaces.ReplaceAllLiteralString(s, "_")
	s = invalidPosixFileNameCharacters.ReplaceAllLiteralString(s, "")
	return s[:min(len(s), 255)]
}
