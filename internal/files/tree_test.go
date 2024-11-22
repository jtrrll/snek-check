package files_test

import (
	"fmt"
	"os"
	"snekcheck/internal/files"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterTree(t *testing.T) {
	var matchAll files.Matcher = func(_ files.Path, _ bool) bool { return true }
	var matchNone files.Matcher = func(_ files.Path, _ bool) bool { return false }

	// Creates a new file system of directories, filled with dummy files
	initFs := func(dirs map[string]uint) billy.Filesystem {
		fs := memfs.New()
		for dir, count := range dirs {
			require.Nil(t, fs.MkdirAll(dir, os.ModeDir))
			for i := range count {
				_, createErr := fs.Create(fs.Join(dir, fmt.Sprintf("%d", i)))
				require.Nil(t, createErr)
			}
		}
		return fs
	}

	t.Parallel()
	t.Run("stops iteration upon break", func(t *testing.T) {
		fs := initFs(map[string]uint{
			"parent": 20,
		})

		var yieldedDirs uint = 0
		var yieldedFiles uint = 0
		for _, fileInfo := range files.IterTree(fs, matchAll, files.NewPath("parent")) {
			if fileInfo.IsDir() {
				yieldedDirs += 1
			} else {
				yieldedFiles += 1
			}
			if yieldedFiles == 10 {
				break
			}
		}
		assert.EqualValues(t, 1, yieldedDirs)
		assert.EqualValues(t, 10, yieldedFiles)
	})
	t.Run("yields only matching paths", func(t *testing.T) {
		fs := initFs(map[string]uint{
			"grandparent":         20,
			"grandparent/parent1": 10,
			"grandparent/parent2": 12,
		})
		var match files.Matcher = func(path files.Path, _ bool) bool {
			return path.Base() == "grandparent" ||
				(path.Parent().Base() == "grandparent" && path.Base() == "parent2") ||
				path.Parent().Base() == "parent2"
		}

		var yieldedDirs uint = 0
		var yieldedFiles uint = 0
		for _, fileInfo := range files.IterTree(fs, match, files.NewPath("grandparent")) {
			if fileInfo.IsDir() {
				yieldedDirs += 1
			} else {
				yieldedFiles += 1
			}
		}
		assert.EqualValues(t, 2, yieldedDirs)
		assert.EqualValues(t, 12, yieldedFiles)
	})
	t.Run("yields everything if all paths match", func(t *testing.T) {
		fs := initFs(map[string]uint{
			"grandparent":         20,
			"grandparent/parent1": 10,
			"grandparent/parent2": 12,
		})

		var yieldedDirs uint = 0
		var yieldedFiles uint = 0
		for _, fileInfo := range files.IterTree(fs, matchAll, files.NewPath("grandparent")) {
			if fileInfo.IsDir() {
				yieldedDirs += 1
			} else {
				yieldedFiles += 1
			}
		}
		assert.EqualValues(t, 3, yieldedDirs)
		assert.EqualValues(t, 42, yieldedFiles)
	})
	t.Run("yields nothing if no paths match", func(t *testing.T) {
		fs := initFs(map[string]uint{
			"grandparent":         20,
			"grandparent/parent1": 10,
			"grandparent/parent2": 12,
		})

		var yieldedDirs uint = 0
		var yieldedFiles uint = 0
		for _, fileInfo := range files.IterTree(fs, matchNone, files.NewPath("grandparent")) {
			if fileInfo.IsDir() {
				yieldedDirs += 1
			} else {
				yieldedFiles += 1
			}
		}
		assert.EqualValues(t, 0, yieldedDirs)
		assert.EqualValues(t, 0, yieldedFiles)
	})
	t.Run("yields nothing if the starting path is invalid", func(t *testing.T) {
		fs := initFs(map[string]uint{
			"grandparent": 10,
		})

		var yieldedDirs uint = 0
		var yieldedFiles uint = 0
		for _, fileInfo := range files.IterTree(fs, matchNone, files.NewPath("invalid")) {
			if fileInfo.IsDir() {
				yieldedDirs += 1
			} else {
				yieldedFiles += 1
			}
		}
		assert.EqualValues(t, 0, yieldedDirs)
		assert.EqualValues(t, 0, yieldedFiles)
	})
	t.Run("yields parent directories before their children", func(t *testing.T) {
		fs := initFs(map[string]uint{
			"grandparent":         20,
			"grandparent/parent1": 10,
			"grandparent/parent2": 12,
		})

		grandparentYielded := false
		parentYielded := false
		parent2Yielded := false
		for path := range files.IterTree(fs, matchAll, files.NewPath("grandparent")) {
			if !grandparentYielded && path.Base() == "grandparent" {
				grandparentYielded = true
				continue
			}
			if !parentYielded && path.Base() == "parent" {
				parentYielded = true
				continue
			}
			if !parent2Yielded && path.Base() == "parent2" {
				parent2Yielded = true
				continue
			}

			if path.Parent().Base() == "grandparent" {
				assert.True(t, grandparentYielded)
			} else if path.Parent().Base() == "parent" {
				assert.True(t, grandparentYielded)
				assert.True(t, parentYielded)
			} else if path.Parent().Base() == "parent2" {
				assert.True(t, grandparentYielded)
				assert.True(t, parent2Yielded)
			}
		}
	})
}
