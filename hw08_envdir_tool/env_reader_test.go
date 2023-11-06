package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("case exists environment var", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")
		if err != nil {
			require.Fail(t, "Error: %v", err)
		}

		require.NoError(t, err)
		require.NotNil(t, envs)

		for _, key := range []string{"BAR", "EMPTY", "FOO"} {
			value, ok := envs[key]

			require.True(t, ok)
			require.False(t, value.NeedRemove)
		}
	})

	t.Run("case not exists environment var", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")
		if err != nil {
			require.Fail(t, "Error: %v", err)
		}

		require.NoError(t, err)
		require.NotNil(t, envs)

		for _, key := range []string{"UNSET"} {
			value, ok := envs[key]

			require.True(t, ok)
			require.True(t, value.NeedRemove)
			require.Empty(t, value.Value)
		}
	})
}
