package incus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

func ContainerExists(name string, verbose bool, startTime time.Time) (bool, string, int, string, time.Duration) {
	cmdStart := time.Now()
	cmd := exec.Command("incus", "ls", "--format=json")

	if verbose {
		fmt.Printf("[%s] Running command: %s\n", time.Since(startTime).Truncate(time.Second), cmd.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(cmdStart)

	var exitCode int
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	} else {
		exitCode = 0
	}

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	var containers []map[string]interface{}
	if err := json.Unmarshal([]byte(stdoutStr), &containers); err != nil {
		return false, "", exitCode, stderrStr, duration
	}

	for _, container := range containers {
		if container["name"] == name {
			return true, stdoutStr, exitCode, stderrStr, duration
		}
	}

	return false, stdoutStr, exitCode, stderrStr, duration
}

func RemoveContainer(name string, verbose bool, startTime time.Time) (string, int, string, time.Duration) {
	cmdStart := time.Now()
	cmd := exec.Command("incus", "rm", "--force", name)

	if verbose {
		fmt.Printf("[%s] Running command: %s\n", time.Since(startTime).Truncate(time.Second), cmd.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(cmdStart)

	var exitCode int
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	} else {
		exitCode = 0
	}

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	return stdoutStr, exitCode, stderrStr, duration
}

func WaitForContainerRemoval(name string, verbose bool, startTime time.Time) bool {
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return false
		case <-ticker.C:
			exists, _, _, _, duration := ContainerExists(name, verbose, startTime)
			if !exists {
				if verbose {
					fmt.Printf("[%s] Container %s removed successfully (took %s)\n", time.Since(startTime).Truncate(time.Second), name, duration.Truncate(time.Second))
				}
				return true
			}
			if verbose {
				fmt.Printf("[%s] Container %s still exists, removing...\n", time.Since(startTime).Truncate(time.Second), name)
			}
			RemoveContainer(name, verbose, startTime)
		}
	}
}

func LaunchContainer(imageName, containerName string, verbose bool, startTime time.Time) (string, int, string, time.Duration) {
	cmdStart := time.Now()
	cmd := exec.Command("incus", "launch", imageName, containerName)

	if verbose {
		fmt.Printf("[%s] Running command: %s\n", time.Since(startTime).Truncate(time.Second), cmd.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(cmdStart)

	var exitCode int
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	} else {
		exitCode = 0
	}

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	return stdoutStr, exitCode, stderrStr, duration
}

func WaitForContainerCreation(name string, verbose bool, startTime time.Time) bool {
	timeout := time.After(20 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return false
		case <-ticker.C:
			exists, _, _, _, duration := ContainerExists(name, verbose, startTime)
			if exists {
				if verbose {
					fmt.Printf("[%s] Container %s created successfully (took %s)\n", time.Since(startTime).Truncate(time.Second), name, duration.Truncate(time.Second))
				}
				return true
			}
			if verbose {
				fmt.Printf("[%s] Waiting for container %s to be created...\n", time.Since(startTime).Truncate(time.Second), name)
			}
		}
	}
}

func ProcessContainerCommand(filter, name string, verbose bool) {
	startTime := time.Now()

	exists, _, _, _, _ := ContainerExists(name, verbose, startTime)
	if exists {
		fmt.Printf("Container %s already exists. Removing...\n", name)
		if !WaitForContainerRemoval(name, verbose, startTime) {
			fmt.Printf("Failed to remove container %s\n", name)
			return
		}
		fmt.Printf("Container %s removed successfully\n", name)
	}

	imageName := FindImageByAlias(filter, verbose, startTime)
	if imageName == "" {
		fmt.Printf("No image found matching filter: %s\n", filter)
		return
	}

	fmt.Printf("Launching container %s from image %s\n", name, imageName)
	_, _, _, duration := LaunchContainer(imageName, name, verbose, startTime)

	if verbose {
		fmt.Printf("[%s] Launch command took %s\n", time.Since(startTime).Truncate(time.Second), duration.Truncate(time.Second))
	}

	if !WaitForContainerCreation(name, verbose, startTime) {
		fmt.Printf("Failed to create container %s\n", name)
		return
	}

	fmt.Printf("Container %s created successfully\n", name)
}
