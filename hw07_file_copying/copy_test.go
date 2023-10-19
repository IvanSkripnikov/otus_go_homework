package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	fileInput := "./testdata/input.txt"
	fileOutput := "./testdata/test.txt"

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
}
