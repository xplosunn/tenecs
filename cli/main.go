package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tenecs",
	Short: "Utilities of the Tenecs programming language",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1-alpha")
	},
}

var formatCmd = &cobra.Command{
	Use:   "format [FILE]",
	Short: "Format the code",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Please provide a file")
		}
		filePath := args[0]
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		fileContent := string(bytes)
		parsed, err := parser.ParseString(fileContent)
		if err != nil {
			return err
		}
		formatted := formatter.DisplayFileTopLevel(*parsed)
		err = os.WriteFile(filePath, []byte(formatted), 0644)
		if err != nil {
			return err
		}
		return nil
	},
}

var runCmd = &cobra.Command{
	Use:   "run [FILE]",
	Short: "Run the code",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Please provide a file")
		}

		filePath := args[0]
		return compileAndRun(false, filePath)
	},
}

var testCmd = &cobra.Command{
	Use:   "test [FILE]",
	Short: "Run the tests",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Please provide a file")
		}

		filePath := args[0]
		return compileAndRun(true, filePath)
	},
}

func compileAndRun(testMode bool, filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent := string(bytes)
	parsed, err := parser.ParseString(fileContent)
	if err != nil {
		return err
	}
	ast, err := typer.Typecheck(*parsed)
	if err != nil {
		return err
	}
	generated := codegen.Generate(testMode, ast)
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	generatedFilePath := filepath.Join(dir, "main.go")
	if err != nil {
		return err
	}
	_, err = os.Create(generatedFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(generatedFilePath, []byte(generated), 0644)
	if err != nil {
		return err
	}
	runCmd := exec.Command("go", "run", generatedFilePath)
	runCmd.Dir = dir
	outputBytes, err := runCmd.Output()
	if err != nil {
		return err
	}
	fmt.Print(string(outputBytes))

	return nil
}
