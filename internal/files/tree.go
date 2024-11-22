package files

import (
	"io/fs"
	"iter"

	"github.com/go-git/go-billy/v5"
)

// A Matcher determines if a file path matches implementation-specific constraints or not.
type Matcher func(path Path, isDir bool) bool

// Iterates over a file tree, only producing paths that match the given matcher.
// The iterator is guaranteed to yield parent directories before their children.
func IterTree(fileSystem billy.Filesystem, match Matcher, p Path) iter.Seq2[Path, fs.FileInfo] {
	return func(yield func(path Path, fileInfo fs.FileInfo) bool) {
		// Process this path
		fileInfo, statErr := fileSystem.Stat(p.String())
		if statErr != nil {
			return
		}

		if !match(p, fileInfo.IsDir()) {
			return
		}

		if !yield(p, fileInfo) || !fileInfo.IsDir() {
			return
		}

		// Attempt to read directory entries
		entries, readErr := fileSystem.ReadDir(p.String())
		if readErr != nil {
			entries = nil
		}

		// Process entries if it is a directory
		for _, entry := range entries {
			for entryPath, fileInfo := range IterTree(fileSystem, match, append(p, entry.Name())) {
				if !yield(entryPath, fileInfo) {
					return
				}
			}
		}
	}
}
