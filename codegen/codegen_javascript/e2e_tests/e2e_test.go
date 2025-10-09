package e2e_tests

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen/codegen_javascript"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/type_error"
)

func Test(t *testing.T) {
	dirEntries, err := os.ReadDir(".")
	assert.NoError(t, err)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			t.Run(dirEntry.Name(), func(t *testing.T) {
				programBytes, err := os.ReadFile(dirEntry.Name() + "/in.10x")
				assert.NoError(t, err)
				program := string(programBytes)

				htmlBytes, err := os.ReadFile(dirEntry.Name() + "/out.html")
				assert.NoError(t, err)
				html := string(htmlBytes)

				parsed, err := parser.ParseString(program)
				assert.NoError(t, err)

				desugared, err := desugar.Desugar(*parsed)
				assert.NoError(t, err)

				typed, err := typer.TypecheckSingleFile(desugared)
				if err != nil {
					t.Fatal(type_error.Render(program, err.(*type_error.TypecheckError)))
				}

				irProgram := ir.ToIR(*typed)

				generatedHtml := codegen_javascript.GenerateHtmlPageForWebApp(&irProgram, ir.RefToIR(ast.Ref{
					Package: "mypage",
					Name:    "webApp",
				}), nil)
				assert.Equal(t, html, generatedHtml)
			})
		}
	}
}
