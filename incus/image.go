package incus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type ImageInfo struct {
	Aliases []struct {
		Name string `json:"name"`
	} `json:"aliases"`
}

func ImageAliases(verbose bool, startTime time.Time) ([]string, int, string, time.Duration) {
	cmdStart := time.Now()
	cmd := exec.Command("incus", "image", "ls", "--format=json")

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

	var images []ImageInfo
	if err := json.Unmarshal([]byte(stdoutStr), &images); err != nil {
		return nil, exitCode, stderrStr, duration
	}

	var aliasNames []string
	for _, image := range images {
		for _, alias := range image.Aliases {
			aliasNames = append(aliasNames, alias.Name)
		}
	}

	return aliasNames, exitCode, stderrStr, duration
}

func ReportImageAliases(filter string, verbose bool, startTime time.Time) {
	aliasNames, exitCode, stderrStr, duration := ImageAliases(verbose, startTime)

	if exitCode != 0 {
		fmt.Printf("Error: %s\n", stderrStr)
		return
	}

	sort.Strings(aliasNames)

	for _, aliasName := range aliasNames {
		if filter == "" || strings.Contains(strings.ToLower(aliasName), strings.ToLower(filter)) {
			fmt.Println(aliasName)
		}
	}

	if verbose {
		fmt.Printf("[%s] Image alias retrieval took %s\n", time.Since(startTime).Truncate(time.Second), duration.Truncate(time.Second))
	}
}

func ContainerInfo(container string, verbose bool, startTime time.Time) (string, int, string, time.Duration) {
	cmdStart := time.Now()
	cmd := exec.Command("incus", "info", container)

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

func FindImageByAlias(filter string, verbose bool, startTime time.Time) string {
	aliasNames, _, _, duration := ImageAliases(verbose, startTime)

	for _, aliasName := range aliasNames {
		if strings.Contains(strings.ToLower(aliasName), strings.ToLower(filter)) {
			return aliasName
		}
	}

	if verbose {
		fmt.Printf("[%s] Image alias search took %s\n", time.Since(startTime).Truncate(time.Second), duration.Truncate(time.Second))
	}

	return ""
}

func ProcessImageCommand(filter, container string, verbose bool) {
	startTime := time.Now()

	if container == "" {
		ReportImageAliases(filter, verbose, startTime)
	} else {
		stdoutStr, exitCode, stderrStr, duration := ContainerInfo(container, verbose, startTime)
		if exitCode != 0 {
			fmt.Printf("Error: %s\n", stderrStr)
		} else {
			fmt.Println(stdoutStr)
		}

		if verbose {
			fmt.Printf("[%s] Container info retrieval took %s\n", time.Since(startTime).Truncate(time.Second), duration.Truncate(time.Second))
		}
	}
}
