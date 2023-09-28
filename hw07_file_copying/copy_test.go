package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	fileInput := "./testdata/input.txt"
	fileOutput := "./testdata/test.txt"

	t.Run("case error ErrOpenFile", func(t *testing.T) {
		err := Copy("random_name.txt", fileOutput, 50, 0)
		require.Truef(t, errors.Is(err, ErrOpenFile), "actual error %q", err)
	})

	t.Run("case error ErrOffsetExceedsFileSize", func(t *testing.T) {
		err := Copy("./testdata/1.txt", fileOutput, 50, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("case custom copy success", func(t *testing.T) {
		res := Copy("./testdata/1.txt", fileOutput, 3, 1)

		file, _ := os.Open(fileOutput)
		content, _ := io.ReadAll(file)

		defer func() {
			file.Close()
			os.Remove(fileOutput)
		}()

		require.Nil(t, res)
		require.Equal(t, "s", string(content))
	})

	t.Run("case offset=0 limit=0", func(t *testing.T) {
		res := Copy(fileInput, fileOutput, 0, 0)

		fileOut, _ := os.Open(fileOutput)
		contentOut, _ := io.ReadAll(fileOut)

		fileMatch, _ := os.Open("./testdata/out_offset0_limit0.txt")
		contentMatch, _ := io.ReadAll(fileMatch)

		defer func() {
			fileOut.Close()
			fileMatch.Close()
			os.Remove(fileOutput)
		}()

		require.Nil(t, res)
		require.Equal(t, string(contentMatch), string(contentOut))
	})

	t.Run("case offset=0 limit=10", func(t *testing.T) {
		res := Copy(fileInput, fileOutput, 0, 10)

		fileOut, _ := os.Open(fileOutput)
		contentOut, _ := io.ReadAll(fileOut)

		fileMatch, _ := os.Open("./testdata/out_offset0_limit10.txt")
		contentMatch, _ := io.ReadAll(fileMatch)

		defer func() {
			fileOut.Close()
			fileMatch.Close()
			os.Remove(fileOutput)
		}()

		require.Nil(t, res)
		require.Equal(t, string(contentMatch), string(contentOut))
	})
}
