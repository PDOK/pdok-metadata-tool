package main

import (
	"fmt"
	clidocs "github.com/urfave/cli-docs/v3"
	"os"
	"os/exec"
	"pdok-metadata-tool/internal/app"
)

func main() {
	generateMarkdownDocs("docs/README.md")
	buildPMT()
}

func buildPMT() {
	fmt.Println("Building pmt executable...")

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory: %v\n", err)
		return
	}

	// Navigate to project root (assuming we're in cmd/generate)
	err = os.Chdir("../..")
	if err != nil {
		fmt.Printf("Failed to change to project root: %v\n", err)
		return
	}

	// Ensure we go back to the original directory when done
	defer func() {
		err := os.Chdir(currentDir)
		if err != nil {
			fmt.Printf("Failed to return to original directory: %v\n", err)
		}
	}()

	// Run the build command from the project root
	cmd := exec.Command("go", "build", "-o", "pmt", "cmd/pmt/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build pmt: %v\n", err)
		return
	}

	fmt.Println("Successfully built pmt executable in project root")
}

func generateMarkdownDocs(filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Errorf("failed to create %s: %w", filepath, err)
	}
	defer file.Close()

	md, err := clidocs.ToMarkdown(app.PDOKMetadataToolCLI)
	if err != nil {
		fmt.Errorf("failed generating CLI docs: %w", err)
	}
	_, err = file.WriteString(md)
	if err != nil {
		fmt.Errorf("failed writing CLI docs: %w", err)
	}
	fmt.Println("CLI documentation generated in: " + filepath)
}
