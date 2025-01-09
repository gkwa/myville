package incus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ImageInfo struct {
	Aliases []struct {
		Name string `json:"name"`
	} `json:"aliases"`
	Fingerprint string `json:"fingerprint"`
}

type ImageMatch struct {
	Fingerprint string
	Aliases     []string
}

func ImageAliases(verbose bool, startTime time.Time) (map[string]ImageInfo, int, string, time.Duration) {
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

	imageMap := make(map[string]ImageInfo)
	for _, img := range images {
		imageMap[img.Fingerprint] = img
	}

	return imageMap, exitCode, stderrStr, duration
}

func RemoveImage(fingerprint string, verbose bool, startTime time.Time) (string, int, string, time.Duration) {
	cmdStart := time.Now()
	cmd := exec.Command("incus", "image", "rm", fingerprint)

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

func ProcessImageRemoveCommand(filters []string, verbose bool, dryRun bool) {
	startTime := time.Now()
	images, exitCode, stderrStr, _ := ImageAliases(verbose, startTime)

	if exitCode != 0 {
		fmt.Printf("Error getting images: %s\n", stderrStr)
		return
	}

	matchedImages := make(map[string]ImageMatch)
	for _, filter := range filters {
		regex, err := regexp.Compile(filter)
		if err != nil {
			fmt.Printf("Invalid regex pattern '%s': %s\n", filter, err)
			continue
		}

		for fingerprint, image := range images {
			for _, alias := range image.Aliases {
				if regex.MatchString(alias.Name) {
					match, exists := matchedImages[fingerprint]
					if !exists {
						match = ImageMatch{
							Fingerprint: fingerprint,
							Aliases:     make([]string, 0),
						}
					}
					match.Aliases = append(match.Aliases, alias.Name)
					matchedImages[fingerprint] = match
				}
			}
		}
	}

	if len(matchedImages) == 0 {
		fmt.Println("No images matched the provided filters")
		return
	}

	for _, match := range matchedImages {
		if dryRun {
			fmt.Printf("Would remove image %s (aliases: %s)\n", match.Fingerprint[:12], strings.Join(match.Aliases, ", "))
			continue
		}

		fmt.Printf("Removing image %s (aliases: %s)\n", match.Fingerprint[:12], strings.Join(match.Aliases, ", "))
		_, exitCode, stderrStr, _ := RemoveImage(match.Fingerprint, verbose, startTime)
		if exitCode != 0 {
			fmt.Printf("Failed to remove image: %s\n", stderrStr)
		} else if verbose {
			fmt.Printf("Successfully removed image\n")
		}
	}
}

func ReportImageAliases(filter string, verbose bool, startTime time.Time) {
	images, exitCode, stderrStr, duration := ImageAliases(verbose, startTime)

	if exitCode != 0 {
		fmt.Printf("Error: %s\n", stderrStr)
		return
	}

	var aliasNames []string
	for _, image := range images {
		for _, alias := range image.Aliases {
			aliasNames = append(aliasNames, alias.Name)
		}
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

func FindImageByAlias(filter string, verbose bool, startTime time.Time) (string, error) {
	images, exitCode, stderrStr, duration := ImageAliases(verbose, startTime)
	if exitCode != 0 {
		return "", fmt.Errorf("error getting images: %s", stderrStr)
	}

	var matchingImages []string
	for _, image := range images {
		for _, alias := range image.Aliases {
			if strings.Contains(strings.ToLower(alias.Name), strings.ToLower(filter)) {
				matchingImages = append(matchingImages, alias.Name)
			}
		}
	}

	if verbose {
		fmt.Printf("[%s] Image alias search took %s\n", time.Since(startTime).Truncate(time.Second), duration.Truncate(time.Second))
	}

	if len(matchingImages) == 0 {
		return "", errors.New("no images found matching the filter")
	} else if len(matchingImages) > 1 {
		errorMsg := fmt.Sprintf("multiple images match the filter '%s':\n", filter)
		for _, imageName := range matchingImages {
			errorMsg += imageName + "\n"
		}
		errorMsg += "please refine the filter to match a single image"
		return "", errors.New(errorMsg)
	}

	return matchingImages[0], nil
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
