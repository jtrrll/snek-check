/*
`snekcheck` recursively lints all provided file paths to ensure all filenames are snake_case.

Usage:

	snekcheck <flag> ... <path> ...

If the `--fix` flag is specified, `snekcheck` will attempt to correct invalid filenames.
*/
package main

import (
	"flag"
	"os"
	"path/filepath"
	"snekcheck/internal/files"
	"snekcheck/internal/patterns"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
)

var (
	// A colorful logger.
	logger = configureLogger()
)

// CLI flags.
var (
	fix = flag.Bool("fix", false, "Whether snekcheck should attempt to correct invalid filenames")
)

// A runtime configuration for snek.
type config struct {
	fix        bool
	fileSystem billy.Filesystem
	gitIgnore  files.GitIgnore
	paths      []files.Path
}

// The snekcheck CLI.
// Will exit with a non-zero exit code upon failure.
func main() {
	rootFs := osfs.New("/")

	flag.Parse()
	relativePaths := flag.Args()
	if len(relativePaths) == 0 {
		logger.Error("no files or directories specified")
		os.Exit(1)
	}

	absolutePaths := make([]files.Path, len(relativePaths))
	for i, path := range relativePaths {
		absolutePath, absoluteErr := filepath.Abs(path)
		if absoluteErr != nil {
			logger.Error(absoluteErr)
			os.Exit(1)
		}
		_, statErr := rootFs.Stat(absolutePath)
		if statErr != nil {
			logger.Error(absoluteErr)
			os.Exit(1)
		}
		absolutePaths[i] = files.NewPath(absolutePath)
	}

	config := config{
		fix:        *fix,
		fileSystem: rootFs,
		gitIgnore:  loadGlobalGitIgnore(rootFs),
		paths:      absolutePaths,
	}
	os.Exit(SnekCheck(config))
}

// Recursively lints all provided file paths to ensure all filenames are snake_case.
func SnekCheck(config config) (exitCode int) {
	match := func(path files.Path, isDir bool) bool {
		return !config.gitIgnore.Match(path, isDir)
	}

	for _, path := range config.paths {
		for path, fileInfo := range files.IterTree(config.fileSystem, match, path) {
			if fileInfo.IsDir() {
				moreIgnores, ignoreErr := files.ParseGitIgnore(config.fileSystem, path)
				if ignoreErr != nil {
					moreIgnores = nil
				}
				config.gitIgnore = append(config.gitIgnore, moreIgnores...)
			}

			if isValidFileName(path.Base()) {
				logger.Print("", "VALID", path)
				continue
			}

			if !config.fix {
				exitCode = 1
				logger.Print("", "INVALID", path)
				continue
			}

			var newPath files.Path
			newPath = append(append(newPath, path.Parent()...), patterns.ToSnakeCase(path.Base()))
			renameErr := config.fileSystem.Rename(path.String(), newPath.String())
			if renameErr != nil {
				exitCode = 1
				logger.Error(renameErr)
				return
			}
			logger.Print("", "FIXED", path)
		}
	}
	return
}

// Determines if a filename is valid according to snekcheck's opinion.
func isValidFileName(name string) bool {
	if !patterns.IsPosixFileName(name) {
		return false
	}

	if patterns.IsSnakeCase(name) {
		return true
	}

	name, extension := splitExtension(name)
	return patterns.IsScreamingSnakeCase(name) && patterns.IsSnakeCase(extension)
}

// Configures the CLI logger.
func configureLogger() *log.Logger {
	l := log.New(os.Stderr)
	styles := log.DefaultStyles()
	styles.Keys["INVALID"] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f44747"))
	styles.Values["INVALID"] = lipgloss.NewStyle()
	styles.Keys["VALID"] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6a9955"))
	styles.Values["VALID"] = lipgloss.NewStyle()
	styles.Keys["FIXED"] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#dcdcaa"))
	styles.Values["FIXED"] = lipgloss.NewStyle()
	l.SetStyles(styles)
	return l
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

// Splits the extension from a filename.
func splitExtension(name string) (string, string) {
	lastIndex := strings.LastIndex(name, ".")
	if lastIndex == -1 {
		return name, ""
	}
	return name[:lastIndex], name[lastIndex:]
}
