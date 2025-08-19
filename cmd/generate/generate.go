package main

import (
	"fmt"
	"github.com/pdok/pdok-metadata-tool/internal/app"
	clidocs "github.com/urfave/cli-docs/v3"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Generating artefacts for PDOK metadata tools")
	generateMarkdownDocs("docs/README.md")
	buildApp()
	fmt.Println("Generation success")
}

func buildApp() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory: %v\n", err)
		return
	}

	fmt.Println("Building pmt executable at: " + currentDir)

	// Run the build command from the project root
	cmd := exec.Command("go", "build", "-o", "pmt", "cmd/pmt/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build pmt: %v\n", err)
		return
	}
}

func generateMarkdownDocs(filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("failed to create %s: %v", filepath, err)
	}
	defer file.Close()

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
