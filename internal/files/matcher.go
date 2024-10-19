package files

// A Matcher determines if a file path matches implementation-specific constraints or not.
type Matcher func(path Path, isDir bool) bool
