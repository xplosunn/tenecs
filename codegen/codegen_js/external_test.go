package codegen_js_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen/codegen_js"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/node"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
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
	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)
	programJs := codegen_js.GenerateProgramNonRunnable(typed)
	js := codegen_js.NodeProgramToPrintWebAppExternalGenerate(webAppPackageName, programJs, webAppVarName)
	jsOutput, err := node.RunCodeBlockingAndReturningOutputWhenFinished(t, js)
	assert.NoError(t, err)
	result, err := codegen_js.NodeProgramToPrintWebAppExternalReadOutput(jsOutput)
	assert.NoError(t, err)
	assert.Equal(t, []string{"fake_css_url.css"}, result)
}
