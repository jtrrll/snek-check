package main_test

import (
	main "snekcheck/cmd/snekcheck"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Parallel()
	t.Run("identifies valid file names", func(t *testing.T) {
		testCases := []string{
			"main.go",
			"flake.nix",
			"LICENSE",
			"README.md",
			"snake_test.go",
		}
		for _, input := range testCases {
			t.Run(input, func(t *testing.T) {
				assert.True(t, main.IsValid(input))
			})
		}
	})
	t.Run("identifies invalid file names", func(t *testing.T) {
		testCases := []string{
			"Snake",
			"snake case 123",
			"snake-case",
			"Readme.md",
			"snake.PNG",
		}
		for _, input := range testCases {
			t.Run(input, func(t *testing.T) {
				assert.False(t, main.IsValid(input))
			})
		}
	})
}
