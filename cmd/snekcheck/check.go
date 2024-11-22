package main

import (
	"snekcheck/internal/files"

	"github.com/go-git/go-billy/v5"
)

// Determines if a collection of filenames are valid according to snekcheck's opinionated validator.
// Recursively descends into directories.
func Check(fs billy.Filesystem, paths []files.Path) (validPaths []files.Path, invalidPaths []files.Path) {
	if fs == nil {
		panic("invalid filesystem")
	}

	gitIgnore := loadGlobalGitIgnore(fs)
	match := func(path files.Path, isDir bool) bool {
		return !gitIgnore.Match(path, isDir)
	}

	validPaths = make([]files.Path, 0, len(paths))
	invalidPaths = make([]files.Path, 0, len(paths))

	for _, path := range paths {
		for path, fileInfo := range files.IterTree(fs, match, path) {
			if fileInfo.IsDir() {
				gitIgnore = append(gitIgnore, parseGitIgnorePatterns(fs, path)...)
			}

			if IsValid(path.Base()) {
				logger.Print("", "VALID", path)
				validPaths = append(validPaths, path)
			} else {
				logger.Print("", "INVALID", path)
				invalidPaths = append(invalidPaths, path)
			}
		}
	}
	return
}
