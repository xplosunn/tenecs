package test_standard_library

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/codegen/codegen_js"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/external/node"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
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
		programBytes, err := os.ReadFile(dirEntry.Name())
		assert.NoError(t, err)

		program := string(programBytes)

		parsed, err := parser.ParseString(program)
		if err != nil {
			t.Log("Failed to parse " + dirEntry.Name())
		}
		assert.NoError(t, err)

		typed, err := typer.TypecheckSingleFile(*parsed)
		if err != nil {
			t.Fatal(type_error.Render(program, err.(*type_error.TypecheckError)))
		}
		foundTests := codegen.FindTests(typed)
		t.Run("go_"+dirEntry.Name(), func(t *testing.T) {
			runTestInGolang(t, typed, foundTests)
		})
		t.Run("node_"+dirEntry.Name(), func(t *testing.T) {
			runTestInNode(t, typed, foundTests)
		})
	}
	assert.True(t, atLeastOneFileFound)
}

func runTestInGolang(t *testing.T, program *ast.Program, foundTests codegen.FoundTests) {
	generated := codegen_golang.GenerateProgramTest(program, foundTests)

	output := golang.RunCodeUnlessCached(t, generated)
	if strings.Contains(output, codegen_golang.Red("FAILURE")) {
		t.Fatal(output)
	}
	if !strings.Contains(output, codegen_golang.Green("OK")) {
		t.Fatal(output)
	}
	if !strings.Contains(output, "* 0 failed") {
		t.Fatal(output)
	}
}

func runTestInNode(t *testing.T, program *ast.Program, foundTests codegen.FoundTests) {
	generated := codegen_js.GenerateProgramTest(program, foundTests)

	output, err := node.RunCodeBlockingAndReturningOutputWhenFinished(t, generated)
	assert.NoError(t, err)
	if strings.Contains(output, codegen_golang.Red("FAILURE")) {
		t.Fatal(output)
	}
	if !strings.Contains(output, codegen_golang.Green("OK")) {
		t.Fatal(output)
	}
	if !strings.Contains(output, "* 0 failed") {
		t.Fatal(output)
	}
}
