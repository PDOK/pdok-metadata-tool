// Package main provides a script for generating PDOK metadata tool artefacts.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pdok/pdok-metadata-tool/v2/internal/app"
	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	clidocs "github.com/urfave/cli-docs/v3"
)

var (
	projRoot = common.GetProjectRoot()
)

func main() {
	fmt.Println("Generating artefacts for PDOK metadata tools")

	path := filepath.Join(projRoot, "docs/README.md")
	generateMarkdownDocs(path)

	if success := buildApp(); success {
		fmt.Println("Generation success")
	} else {
		fmt.Println("Generation failed")
	}
}

func buildApp() bool {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory: %v\n", err)

		return false
	}

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("Failed to read entries from current directory: %v\n", err)

		return false
	}

	const executableName = "pmt"

	found := false

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == executableName {
			found = true

			break
		}
	}

	if found {
		fmt.Printf(
			"Cannot build executable since a directory named %s exists in %s\n",
			executableName,
			currentDir,
		)

		return false
	}

	fmt.Println("Building pmt executable at: " + currentDir)
	// Run the build command from the project root
	path := filepath.Join(projRoot, "cmd/pmt/main.go")
	//nolint:gosec
	cmd := exec.Command("go", "build", "-o", executableName, path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build pmt: %v\n", err)

		return false
	}

	return true
}

func generateMarkdownDocs(filepath string) {
	//nolint:gosec
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("failed to create %s: %v", filepath, err)
	}
	defer common.SafeClose(file)

	md, err := clidocs.ToMarkdown(app.PDOKMetadataToolCLI)
	if err != nil {
		fmt.Printf("failed generating CLI docs: %v", err)
	}

	_, err = file.WriteString(md)
	if err != nil {
		fmt.Printf("failed writing CLI docs: %v", err)
	}

	fmt.Println("CLI documentation generated in: " + filepath)
}
