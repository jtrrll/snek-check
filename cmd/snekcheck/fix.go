package main

import (
	"fmt"
	"snekcheck/internal/files"
	"snekcheck/internal/patterns"

	"github.com/go-git/go-billy/v5"
)

// Renames a collection of filenames to satisfy snekcheck's opinionated validator.
// Recursively descends into directories.
func Fix(fs billy.Filesystem, paths []files.Path) (validPaths []files.Path, renamedPaths []renamedPath) {
	if fs == nil {
		panic("invalid filesystem")
	}

	gitIgnore := loadGlobalGitIgnore(fs)
	match := func(path files.Path, isDir bool) bool {
		return !gitIgnore.Match(path, isDir)
	}

	validPaths = make([]files.Path, 0, len(paths))
	renamedPaths = make([]renamedPath, 0, len(paths))

	for _, path := range paths {
		for path, fileInfo := range files.IterTree(fs, match, path) {
			if fileInfo.IsDir() {
				gitIgnore = append(gitIgnore, parseGitIgnorePatterns(fs, path)...)
			}

			if IsValid(path.Base()) {
				validPaths = append(validPaths, path)
				continue
			}

			var newPath files.Path
			newPath = append(append(newPath, path.Parent()...), patterns.ToSnakeCase(path.Base()))
			if fs.Rename(path.String(), newPath.String()) != nil {
				panic(fmt.Errorf("unable to rename %s to %s", path.String(), newPath.String()))
			}
			renamedPaths = append(renamedPaths, renamedPath{old: path, new: newPath})
		}
	}
	return
}

// A path that has been renamed.
type renamedPath struct {
	old files.Path
	new files.Path
}
