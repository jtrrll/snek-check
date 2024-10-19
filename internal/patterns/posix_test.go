package patterns_test

import (
	"snek-check/internal/patterns"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPosix(t *testing.T) {
	t.Parallel()
	t.Run("IsPosixFileName()", func(t *testing.T) {
		t.Run("identifies valid POSIX filenames", func(t *testing.T) {
			testCases := []string{
				"posix",
				"POSIX_FILE",
				"_POSIX__FILE_.md",
				"012_345",
				"FILE1.txt",
			}
			for _, input := range testCases {
				t.Run(input, func(t *testing.T) {
					assert.True(t, patterns.IsPosixFileName(input))
				})
			}
		})
		t.Run("identifies invalid POSIX filenames", func(t *testing.T) {
			testCases := []string{
				"-TEST",
				"lol@email",
				"invalid%",
				"file(12).pdf",
			}
			for _, input := range testCases {
				t.Run(input, func(t *testing.T) {
					assert.False(t, patterns.IsPosixFileName(input))
				})
			}
		})
	})
	t.Run("ToPosixFileName()", func(t *testing.T) {
		t.Run("does not change valid POSIX filenames", func(t *testing.T) {
			testCases := []string{
				"POSIX",
				".POSIX_123_.md",
				"_DO_NOT_CHANGE_THIS_PLEASE___",
			}
			for _, input := range testCases {
				t.Run(input, func(t *testing.T) {
					require.True(t, patterns.IsPosixFileName(input))
					assert.Equal(t, input, patterns.ToPosixFileName(input))
				})
			}
		})
		t.Run("converts invalid POSIX filenames to valid POSIX filenames", func(t *testing.T) {
			testCases := []struct {
				input  string
				output string
			}{
				{input: "lol#$", output: "lol"},
				{input: "spaced  name", output: "spaced_name"},
				{input: "__012 345.md", output: "__012_345.md"},
			}
			for _, tc := range testCases {
				t.Run(tc.input, func(t *testing.T) {
					require.False(t, patterns.IsPosixFileName(tc.input))
					require.True(t, patterns.IsPosixFileName(tc.output))
					actual := patterns.ToPosixFileName(tc.input)
					assert.Equal(t, tc.output, actual)
					assert.True(t, patterns.IsPosixFileName(actual))
				})
			}
		})
	})
}

func BenchmarkPosix(b *testing.B) {
	b.Run("IsPosixFileName()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			patterns.IsPosixFileName("Bench mark")
		}
	})
	b.Run("ToPosixFileName()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			patterns.ToPosixFileName("Bench mark")
		}
	})
}
