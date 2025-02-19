package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/codegen/codegen_js"
	"github.com/xplosunn/tenecs/external/node"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
			fmt.Println(err.Error())
			return nil
		}
		fileContent := string(bytes)
		parsed, err := parser.ParseString(fileContent)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		formatted := formatter.DisplayFileTopLevel(*parsed)
		err = os.WriteFile(filePath, []byte(formatted), 0644)
		if err != nil {
			fmt.Println(err.Error())
			return nil
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
		compileAndRun(false, filePath)
		return nil
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
		compileAndRun(true, filePath)
		return nil
	},
}

func compileAndRun(testMode bool, filePath string) {
	files, err := getFiles(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	parsedFiles := map[string]parser.FileTopLevel{}
	fileContents := map[string]string{}
	for _, filePath := range files {
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fileContent := string(bytes)
		fileContents[filePath] = fileContent
		parsed, err := parser.ParseString(fileContent)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		parsedFiles[filePath] = *parsed
	}
	ast, err := typer.TypecheckPackages(parsedFiles)
	if err != nil {
		typecheckError := err.(*type_error.TypecheckError)
		rendered, err2 := type_error.Render(fileContents[typecheckError.File], typecheckError)
		if err2 != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(rendered)
		return
	}
	if testMode {
		foundTests := codegen.FindTests(ast)
		generated := codegen_golang.GenerateProgramTest(ast, foundTests)
		runGo(generated)
	} else {
		foundRunnables := codegen.FindRunnables(ast)
		if len(foundRunnables.GoMain) > 1 ||
			len(foundRunnables.WebWebApp) > 1 ||
			(len(foundRunnables.GoMain) > 0 && len(foundRunnables.WebWebApp) > 0) {
			fmt.Println("multiple runnables found")
			return
		} else if len(foundRunnables.GoMain) > 0 {
			targetMain := foundRunnables.GoMain[0]
			generated := codegen_golang.GenerateProgramMain(ast, targetMain)
			runGo(generated)
		} else if len(foundRunnables.WebWebApp) > 0 {
			target := foundRunnables.WebWebApp[0]

			cssFiles, err := func() ([]string, error) {
				programJs := codegen_js.GenerateProgramNonRunnable(ast)
				js := codegen_js.NodeProgramToPrintWebAppExternalGenerate(target.Package, programJs, target.Name)
				jsOutput, err := node.RunCodeBlockingAndReturningOutputWhenFinished(nil, js)
				if err != nil {
					return nil, err
				}
				result, err := codegen_js.NodeProgramToPrintWebAppExternalReadOutput(jsOutput)
				if err != nil {
					return nil, err
				}
				return result, nil
			}()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			html := codegen_js.GenerateHtmlPageForWebApp(ast, target, cssFiles)
			runWebApp(html)
		} else {
			fmt.Println("no runnables found")
			return
		}
	}
}

func getFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(p, ".10x") {
			files = append(files, p)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func runGo(generated string) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	generatedFilePath := filepath.Join(dir, "main.go")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = os.Create(generatedFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = os.WriteFile(generatedFilePath, []byte(generated), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	runCmd := exec.Command("go", "run", generatedFilePath)
	runCmd.Dir = dir
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	err = runCmd.Run()
	if err != nil {
		fmt.Println("error running " + generatedFilePath)
		fmt.Println(err.Error())
		return
	}
}

func runWebApp(html string) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	generatedFilePath := filepath.Join(dir, "index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = os.Create(generatedFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = os.WriteFile(generatedFilePath, []byte(html), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = openURL(generatedFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	println(html)
	println("opened in browser")
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// openURL opens the specified URL in the default browser of the user.
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
