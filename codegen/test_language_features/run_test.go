package test_standard_library

import (
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/codegen/codegen_javascript"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/external/node"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
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

		desugared, err := desugar.Desugar(*parsed)
		assert.NoError(t, err)

		typed, err := typer.TypecheckSingleFile(desugared)
		if err != nil {
			t.Fatal(type_error.Render(program, err.(*type_error.TypecheckError)))
		}
		irProgram := ir.ToIR(*typed)
		foundTests := codegen.FindTests(typed)
		t.Run("go_"+dirEntry.Name(), func(t *testing.T) {
			runTestInGolang(t, &irProgram, foundTests)
		})
		t.Run("node_"+dirEntry.Name(), func(t *testing.T) {
			runTestInNode(t, &irProgram, foundTests)
		})
	}
	assert.True(t, atLeastOneFileFound)
}

func runTestInGolang(t *testing.T, program *ir.Program, foundTests codegen.FoundTests) {
	generated := codegen_golang.GenerateProgramTest(program, foundTests)

	t.Log("Generated:")
	t.Log(generated.String())
	t.Log("--------------------------------")

	output := golang.RunCodeUnlessCached(t, generated.String())
	if strings.Contains(output, "[\u001b[31mFAILURE\u001b[0m]") {
		t.Fatal(output)
	}
	if !strings.Contains(output, "[\u001b[32mOK\u001b[0m]") {
		t.Fatal(output)
	}
	if !strings.Contains(output, "* 0 failed") {
		t.Fatal(output)
	}
}

func runTestInNode(t *testing.T, program *ir.Program, foundTests codegen.FoundTests) {
	generated := codegen_javascript.GenerateProgramTest(program, foundTests)

	t.Log("Generated:")
	t.Log(generated.String())
	t.Log("--------------------------------")

	output, err := node.RunCodeBlockingAndReturningOutputWhenFinished(t, generated.String())
	assert.NoError(t, err)
	if strings.Contains(output, "[\u001b[31mFAILURE\u001b[0m]") {
		t.Fatal(output)
	}
	if !strings.Contains(output, "[\u001b[32mOK\u001b[0m]") {
		t.Fatal(output)
	}
	if !strings.Contains(output, "* 0 failed") {
		t.Fatal(output)
	}
}
