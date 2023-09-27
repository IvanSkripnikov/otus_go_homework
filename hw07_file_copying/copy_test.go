package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("case error ErrOpenFile", func(t *testing.T) {
		err := Copy("random_name.txt", "testdata/2.txt", 50, 0)
		require.Truef(t, errors.Is(err, ErrOpenFile), "actual error %q", err)
	})

	t.Run("case error ErrOffsetExceedsFileSize", func(t *testing.T) {
		err := Copy("./testdata/1.txt", "./testdata/2.txt", 50, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})
}
