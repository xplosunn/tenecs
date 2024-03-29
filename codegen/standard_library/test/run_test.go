package test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"os"
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

	typed, err := typer.TypecheckSingleFile(*parsed)
	if err != nil {
		t.Fatal(type_error.Render(program, err.(*type_error.TypecheckError)))
	}

	generated := codegen.GenerateProgramTest(typed)

	output := golang.RunCodeUnlessCached(t, generated)
	if strings.Contains(output, codegen.Red("FAILURE")) {
		t.Fatal(output)
	}
	if !strings.Contains(output, codegen.Green("OK")) {
		t.Fatal(output)
	}
	if !strings.Contains(output, "* 0 failed") {
		t.Fatal(output)
	}
}
