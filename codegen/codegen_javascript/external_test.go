package codegen_javascript_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen/codegen_javascript"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/node"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
)

func TestNodeProgramToPrintWebAppExternal(t *testing.T) {
	webAppVarName := "myapp"
	webAppPackageName := "test"
	tenecsProgram := `package test

import tenecs.web.CssUrl
import tenecs.web.WebApp
import tenecs.web.HtmlElement
import tenecs.web.HtmlElementProperty

struct State()
struct Event()

myapp := WebApp<State, Event>(
  init = () => State(),
  update = update,
  view = view,
  external = [
    CssUrl("fake_css_url.css")
  ]
)

update := (model: State, event: Event): State => {
  model
}

view := (model: State): HtmlElement<Event> => {
  HtmlElement("p", <HtmlElementProperty<Event>>[], "Hello world!")
}
`
	parsed, err := parser.ParseString(tenecsProgram)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)
	irProgram := ir.ToIR(*typed)
	programJs := codegen_javascript.GenerateProgramNonRunnable(&irProgram)
	js := codegen_javascript.NodeProgramToPrintWebAppExternalGenerate(programJs.String(), ir.RefToIR(ast.Ref{
		Package: webAppPackageName,
		Name:    webAppVarName,
	}))
	jsOutput, err := node.RunCodeBlockingAndReturningOutputWhenFinished(t, js)
	assert.NoError(t, err)
	result, err := codegen_javascript.NodeProgramToPrintWebAppExternalReadOutput(jsOutput)
	assert.NoError(t, err)
	assert.Equal(t, []string{"fake_css_url.css"}, result)
}
