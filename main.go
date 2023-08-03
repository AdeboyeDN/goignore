package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var language string

var rootCmd = &cobra.Command{
	Use:   "goignore",
	Short: "A lightweight CLI tool to generate .gitignore files",
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate and add .gitignore file to your project",
	Run: func(cmd *cobra.Command, args []string) {
		if language == "" {
			color.Red("Please provide a programming language.")
			return
		}

		// Check if .git repo exists, if not initialize it
		_, err := os.Stat(".git")
		if err != nil {
			color.Yellow("Initializing a new Git repository...")
			err := execCommand("git", "init")
			if err != nil {
				color.Red("Error initializing Git repository:", err)
				return
			}
		}

		// Read .gitignore template content from file
		templateContent, err := readTemplateFile(language)
		if err != nil {
			color.Red("Error reading template file:", err)
			return
		}

		// Generate and write the .gitignore file
		err = generateGitignore(templateContent)
		if err != nil {
			color.Red("Error generating .gitignore:", err)
			return
		}

		color.Green("Generated .gitignore for %s", language)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)

	newCmd.PersistentFlags().StringVarP(&language, "language", "l", "", "Programming language for .gitignore file")
	newCmd.MarkPersistentFlagRequired("language")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported programming languages",
	Run: func(cmd *cobra.Command, args []string) {
		supportedLanguages := getSupportedLanguages()
		color.Cyan("Supported programming languages:")
		for _, lang := range supportedLanguages {
			fmt.Println("-", lang)
		}
	},
}

func getSupportedLanguages() []string {
	return []string{"python", "javascript", "golang", "c++"}
}
func main() {
	rootCmd.Execute()
}

func readTemplateFile(language string) (string, error) {
	templatePath := fmt.Sprintf("templates/%s.txt", language)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func generateGitignore(content string) error {
	file, err := os.Create(".gitignore")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func execCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
