package main

import (
	"fmt"
	clidocs "github.com/urfave/cli-docs/v3"
	"os"
	"pdok-metadata-tool/internal/app"
)

func main() {
	generateMarkdownDocs("docs/README.md")
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
