/*
`snekcheck` recursively lints all provided file paths to ensure all filenames are snake_case.

Usage:

	snekcheck <flag> ... <path> ...

If the `--fix` flag is specified, `snekcheck` will attempt to correct invalid filenames.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"snekcheck/internal/files"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
)

var (
	// A colorful logger.
	// TODO: Use a better view solution.
	logger = configureLogger()
)

// CLI flags.
// TODO: Use a better flag library.
var (
	fix = flag.Bool("fix", false, "Whether snekcheck should attempt to correct invalid filenames")
)

// The snekcheck CLI.
// Will exit with a non-zero exit code upon failure.
func main() {
	// Initialize filesystem.
	rootFs := osfs.New("/")
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic("could not determine present working directory")
	}

	// Parse CLI flags and args.
	flag.Parse()

	paths, pathsErr := absPaths(rootFs, pwd, flag.Args())
	if pathsErr != nil {
		logger.Error(pathsErr)
		exit(1)
	}
	if len(paths) == 0 {
		logger.Error("no valid files or directories specified")
		exit(1)
	}

	// Run sneckcheck.
	if *fix {
		Fix(rootFs, paths)
		exit(0)
	}

	_, invalidPaths := Check(rootFs, paths)
	if len(invalidPaths) != 0 {
		exit(1)
	}
	exit(0)
}

// Converts potentially relative paths to separated, absolute paths.
// Errors if any provided path does not exist.
func absPaths(fs billy.Filesystem, pwd string, paths []string) (absPaths []files.Path, err error) {
	if fs == nil {
		panic("invalid filesystem")
	}

	absPaths = make([]files.Path, len(paths))
	for i, path := range paths {
		absPath := path
		if !strings.HasPrefix(path, "/") {
			absPath = fs.Join(pwd, path)
		}
		_, statErr := fs.Stat(absPath)
		if statErr != nil {
			err = fmt.Errorf("no such file or directory: %s", path)
			return
		}
		absPaths[i] = files.NewPath(absPath)
	}
	return
}

// Configures the CLI logger.
func configureLogger() (logger *log.Logger) {
	logger = log.New(os.Stderr)
	styles := log.DefaultStyles()
	styles.Keys["INVALID"] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f44747"))
	styles.Values["INVALID"] = lipgloss.NewStyle()
	styles.Keys["VALID"] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6a9955"))
	styles.Values["VALID"] = lipgloss.NewStyle()
	styles.Keys["FIXED"] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#dcdcaa"))
	styles.Values["FIXED"] = lipgloss.NewStyle()
	logger.SetStyles(styles)
	return
}

// Terminates the current program with the given status code.
// Panics if the exit code is not in the range [0, 125].
func exit(code uint8) {
	if code > 125 {
		panic(fmt.Errorf("invalid exit code: %d", code))
	}
	os.Exit(int(code))
}

// Parses gitignore patterns in a single directory
func parseGitIgnorePatterns(fs billy.Filesystem, path files.Path) files.GitIgnore {
	patterns, ignoreErr := files.ParseGitIgnore(fs, path)
	if ignoreErr != nil {
		patterns = nil
	}
	return patterns
}

// Parses the list of global gitignore patterns.
// Produces an empty list of patterns upon failure.
func loadGlobalGitIgnore(fs billy.Filesystem) files.GitIgnore {
	globalIgnorePatterns, ignoreErr := files.GlobalGitIgnorePatterns(fs)
	if ignoreErr != nil {
		logger.Warn(ignoreErr)
		globalIgnorePatterns = nil
	}
	return globalIgnorePatterns
}
