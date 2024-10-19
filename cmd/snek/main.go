/*
`snek` recursively lints all provided file paths to ensure all filenames are snake_case.

Usage:

	snek <flag> ... <path> ...

If the `--fix` flag is specified, `snek` will attempt to correct invalid filenames.
*/
package main

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"snek-check/internal/files"
	"snek-check/internal/patterns"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-billy/v5/osfs"
)

var (
	// The program exit code.
	exitCode = 0
	// A colorful logger.
	logger = configureLogger()
	// The root file system.
	rootFs = osfs.New("/")
)

// CLI flags.
var (
	fix = flag.Bool("fix", false, "Whether snek should attempt to correct invalid filenames")
)

// A runtime configuration for snek.
type config struct {
	fix   bool
	paths []files.Path
}

// The snek CLI.
// Will exit with a non-zero exit code upon failure.
func main() {
	config := loadConfig()
	gitIgnore := loadGlobalGitIgnore()

	match := func(path files.Path, isDir bool) bool {
		return !gitIgnore.Match(path, isDir)
	}

	fn := func(path files.Path, fileInfo fs.FileInfo) {
		if fileInfo.IsDir() {
			moreIgnores, ignoreErr := files.ParseGitIgnore(rootFs, path)
			if ignoreErr != nil {
				moreIgnores = nil
			}
			gitIgnore = append(gitIgnore, moreIgnores...)
		}

		if !patterns.IsSnakeCase(path.Base()) {
			if config.fix {
				var newPath files.Path
				newPath = append(append(newPath, path.Parent()...), patterns.ToSnakeCase(path.Base()))
				renameErr := rootFs.Rename(path.String(), newPath.String())
				if renameErr != nil {
					exitCode = 1
					logger.Error(renameErr)
					os.Exit(exitCode)
				}
				logger.Print("", "FIXED", path)
			} else {
				exitCode = 1
				logger.Print("", "INVALID", path)
			}
		} else {
			logger.Print("", "VALID", path)
		}
	}

	for _, path := range config.paths {
		for path, fileInfo := range files.IterTree(rootFs, match, path) {
			fn(path, fileInfo)
		}
	}

	os.Exit(exitCode)
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

// Parses CLI inputs into a runtime configuration.
// Will exit with a non-zero exit code upon failure.
func loadConfig() config {
	flag.Parse()

	relativePaths := flag.Args()
	if len(relativePaths) == 0 {
		exitCode = 1
		logger.Error("no files or directories specified")
		os.Exit(exitCode)
	}

	absolutePaths := make([]files.Path, len(relativePaths))
	for i, path := range relativePaths {
		absolutePath, absoluteErr := filepath.Abs(path)
		if absoluteErr != nil {
			exitCode = 1
			logger.Error(absoluteErr)
			os.Exit(exitCode)
		}
		_, statErr := rootFs.Stat(absolutePath)
		if statErr != nil {
			exitCode = 1
			logger.Error(absoluteErr)
			os.Exit(exitCode)
		}
		absolutePaths[i] = files.NewPath(absolutePath)
	}

	return config{
		fix:   *fix,
		paths: absolutePaths,
	}
}

// Parses the list of global gitignore patterns.
// Produces an empty list of patterns upon failure.
func loadGlobalGitIgnore() files.GitIgnore {
	globalIgnorePatterns, ignoreErr := files.GlobalGitIgnorePatterns(rootFs)
	if ignoreErr != nil {
		logger.Warn(ignoreErr)
		globalIgnorePatterns = nil
	}
	return globalIgnorePatterns
}
