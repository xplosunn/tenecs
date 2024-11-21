package e2e_tests

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen/codegen_js"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"os"
	"testing"
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

				typed, err := typer.TypecheckSingleFile(*parsed)
				if err != nil {
					t.Fatal(type_error.Render(program, err.(*type_error.TypecheckError)))
				}
				generatedHtml := codegen_js.GenerateHtmlPageForWebApp(typed, "webApp")
				assert.Equal(t, html, generatedHtml)
			})
		}
	}
}
