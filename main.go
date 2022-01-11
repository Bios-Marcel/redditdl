package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Bios-Marcel/redditdl/redditdl"
	"github.com/kennygrant/sanitize"
)

func main() {
	dirFlag := flag.String("dir", "", "directory that videos will be downloaded.")
	flag.Parse()

	targetDirectory, absError := filepath.Abs(*dirFlag)
	if absError != nil {
		fmt.Printf("Error resolving path '%s'\n", *dirFlag)
		os.Exit(1)
	}

	if *dirFlag != "" {
		mkdirError := os.MkdirAll(targetDirectory, 0777)
		if mkdirError != nil {
			fmt.Printf("Error creating target directory:\n\t%s\n", mkdirError)
			os.Exit(1)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fmt.Printf("Downloading '%s'\n", line)
		downloadError := downloadAndSave(targetDirectory, line)
		if downloadError != nil {
			fmt.Printf("Error downloading post '%s':\n\t%s\n", line, downloadError)
		}
	}
}

func downloadAndSave(targetFolder, postURL string) error {
	baseName := strings.TrimSuffix(sanitize.BaseName(postURL), "-")
	targetFile := path.Join(targetFolder, baseName) + ".mp4"
	output, createError := os.Create(targetFile)
	if createError != nil {
		return createError
	}
	defer output.Close()

	return redditdl.Download(postURL, output)
}
