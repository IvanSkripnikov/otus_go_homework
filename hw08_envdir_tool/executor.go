package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) <= 1 {
		return
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = makeEnv(env)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		log.Printf("Error: %v", err)
	}

	returnCode = command.ProcessState.ExitCode()

	return
}

func makeEnv(env Environment) []string {
	systemEnv := os.Environ()
	envs := make([]string, 0, len(systemEnv)+len(env))

	for key, value := range env {
		if !value.NeedRemove {
			strValue := fmt.Sprintf("%s=%s", key, value.Value)
			envs = append(envs, strValue)
		}
	}

	envs = append(envs, systemEnv...)

	return envs
}
