package main

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyErrors(t *testing.T) {
	fileOutput := "./testdata/test.txt"

	t.Run("case error ErrOpenFile", func(t *testing.T) {
		err := Copy("random_name.txt", fileOutput, 50, 0)
		require.Truef(t, errors.Is(err, ErrOpenFile), "actual error %q", err)
	})

	t.Run("case error ErrOffsetExceedsFileSize", func(t *testing.T) {
		err := Copy("./testdata/1.txt", fileOutput, 50, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("case error ErrUnsupportedFile", func(t *testing.T) {
		err := Copy("./testdata/image.png", fileOutput, 0, 10)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})
}

func TestCopy(t *testing.T) {
	fileInput := "./testdata/input.txt"
	fileOutput := "./testdata/test.txt"

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

	t.Run("case offset=0 limit=1000", func(t *testing.T) {
		res := Copy(fileInput, fileOutput, 0, 1000)

		fileOut, _ := os.Open(fileOutput)
		contentOut, _ := io.ReadAll(fileOut)

		fileMatch, _ := os.Open("./testdata/out_offset0_limit1000.txt")
		contentMatch, _ := io.ReadAll(fileMatch)

		defer func() {
			fileOut.Close()
			fileMatch.Close()
			os.Remove(fileOutput)
		}()

		require.Nil(t, res)
		require.Equal(t, string(contentMatch), string(contentOut))
	})

	t.Run("case offset=0 limit=10000", func(t *testing.T) {
		res := Copy(fileInput, fileOutput, 0, 10000)

		fileOut, _ := os.Open(fileOutput)
		contentOut, _ := io.ReadAll(fileOut)

		fileMatch, _ := os.Open("./testdata/out_offset0_limit10000.txt")
		contentMatch, _ := io.ReadAll(fileMatch)

		defer func() {
			fileOut.Close()
			fileMatch.Close()
			os.Remove(fileOutput)
		}()

		require.Nil(t, res)
		require.Equal(t, string(contentMatch), string(contentOut))
	})

	t.Run("case offset=100 limit=1000", func(t *testing.T) {
		res := Copy(fileInput, fileOutput, 100, 1000)

		fileOut, _ := os.Open(fileOutput)
		contentOut, _ := io.ReadAll(fileOut)

		fileMatch, _ := os.Open("./testdata/out_offset100_limit1000.txt")
		contentMatch, _ := io.ReadAll(fileMatch)

		defer func() {
			fileOut.Close()
			fileMatch.Close()
			os.Remove(fileOutput)
		}()

		require.Nil(t, res)
		require.Equal(t, string(contentMatch), string(contentOut))
	})

	t.Run("case offset=6000 limit=1000", func(t *testing.T) {
		res := Copy(fileInput, fileOutput, 6000, 1000)

		fileOut, _ := os.Open(fileOutput)
		contentOut, _ := io.ReadAll(fileOut)

		fileMatch, _ := os.Open("./testdata/out_offset6000_limit1000.txt")
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
