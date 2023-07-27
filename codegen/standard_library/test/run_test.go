package test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	dirEntries, err := os.ReadDir(".")
	assert.NoError(t, err)
	atLeastOneFileFound := false
	for _, dirEntry := range dirEntries {
		if !strings.HasSuffix(dirEntry.Name(), ".10x") {
			continue
		}
		atLeastOneFileFound = true
		t.Run(dirEntry.Name(), func(t *testing.T) {
			runTest(t, dirEntry.Name())
		})
	}
	assert.True(t, atLeastOneFileFound)
}

func runTest(t *testing.T, fileName string) {
	programBytes, err := os.ReadFile(fileName)
	assert.NoError(t, err)

	program := string(programBytes)

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	if err != nil {
		t.Fatal(type_error.Render(program, err.(*type_error.TypecheckError)))
	}

	generated := codegen.Generate(true, typed)

	createFileAndRun(t, generated)
}

func createFileAndRun(t *testing.T, fileContent string) {
	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	filePath := filepath.Join(dir, "program.go")
	t.Log(filePath)

	_, err = os.Create(filePath)

	contentBytes := []byte(fileContent)
	err = os.WriteFile(filePath, contentBytes, 0644)
	assert.NoError(t, err)

	cmd := exec.Command("go", "run", filePath)
	cmd.Dir = dir
	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		t.Log(err.Error())
		t.Fatal(output)
	}
	if strings.Contains(output, "[FAILURE]") {
		t.Fatal(output)
	}
	if !strings.Contains(output, "[OK]") {
		t.Fatal(output)
	}
}
