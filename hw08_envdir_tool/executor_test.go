package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	var commands1 = []string{"ls", "-la", "/var"}
	var commands2 = []string{"ls", "-la", "/root"}

	envs, err := ReadDir("./testdata/env")
	if err != nil {
		require.Fail(t, "Error: %v", err)
	}

	require.Equal(t, 0, RunCmd(commands1, envs))
	require.Equal(t, 2, RunCmd(commands2, envs))
}
