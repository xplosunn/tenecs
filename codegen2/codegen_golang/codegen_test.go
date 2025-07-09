package codegen_golang_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen2/codegen_golang"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
)

func TestGenerateProgramMainHelloWorld(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any) any)(func(_runtime any) any {
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any) any)(map[string]any{"$type": "String", "value": "Hello world!"})
	})
}

func tenecs_go__Main() any {
	log := func(msg any) any {
		println(msg.(map[string]any)["value"].(string))
		return nil
	}
	console := map[string]any{
		"_log": log,
	}
	refCreator := map[string]any{
		"_new": func(value any) any {
			var ref any = value
			return map[string]any{
				"_get": func() any {
					return ref
				},
				"_set": func(value any) any {
					ref = value
					return nil
				},
				"_modify": func(f any) any {
					ref = f.(func(any) any)(ref)
					return nil
				},
			}
		},
	}
	runtime := map[string]any{
		"_console": console,
		"_ref":     refCreator,
	}
	return func(run any) any {
		return run.(func(any) any)(runtime)
	}
}

func main() {
	main__app()
}
`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	mainPackage := "main"
	generated := codegen_golang.GenerateProgramMain(&codeIR, ir.Reference{
		Name: ir.VariableName(&mainPackage, "app"),
	}).String()
	generatedFormatted := golang.Fmt(t, generated)
	assert.Equal(t, expectedGoCode, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateProgramMainWithRef(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
	ref := runtime.ref.new("hello")
	runtime.console.log(ref.get())
	ref.set("world")
	runtime.console.log(ref.get())
  }
)`
	expectedGoUserspaceCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any) any)(func(_runtime any) any {
		_ref := _runtime.(map[string]any)["_ref"].(map[string]any)["_new"].(func(any) any)(map[string]any{"$type": "String", "value": "hello"})
		_ = _ref
		_runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any) any)(_ref.(map[string]any)["_get"].(func() any)())
		_ref.(map[string]any)["_set"].(func(any) any)(map[string]any{"$type": "String", "value": "world"})
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any) any)(_ref.(map[string]any)["_get"].(func() any)())
	})
}
`

	expectedRunResult := "hello\nworld\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	mainPackage := "main"
	generatedProgram := codegen_golang.GenerateProgramMain(&codeIR, ir.Reference{
		Name: ir.VariableName(&mainPackage, "app"),
	})
	generated := generatedProgram.PackageCode + "\n\n" +
		generatedProgram.ImportsCode + "\n\n" +
		generatedProgram.UserspaceCode
	generatedFormatted := golang.Fmt(t, generated)
	assert.Equal(t, expectedGoUserspaceCode, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generatedProgram.String())
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateProgramMainHelloWorldSeparateFunction(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

helloWorld := (runtime: Runtime): Void => {
  runtime.console.log("Hello world!")
}

app := Main(helloWorld)
`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any) any)(main__helloWorld())
}
func main__helloWorld() any {
	return func(_runtime any) any {
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any) any)(map[string]any{"$type": "String", "value": "Hello world!"})
	}
}
`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	mainPackage := "main"
	generatedProgram := codegen_golang.GenerateProgramMain(&codeIR, ir.Reference{
		Name: ir.VariableName(&mainPackage, "app"),
	})
	generated := generatedProgram.PackageCode + "\n\n" +
		generatedProgram.ImportsCode + "\n\n" +
		generatedProgram.UserspaceCode
	generatedFormatted := golang.Fmt(t, generated)
	assert.Equal(t, expectedGoCode, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generatedProgram.String())
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateProgramMainNestedAssignment(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

funcTakingVoid := (void: Void): Void => {
  void
}

app := Main((runtime: Runtime): Void => {
  void := nestedVar := 1
  funcTakingVoid(void)
  runtime.console.log("Hello world!")
})
`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any) any)(func(_runtime any) any {
		_void := func() any {
			_nestedVar := map[string]any{"$type": "Int", "value": 1}
			_ = _nestedVar
			return nil
		}
		_ = _void
		main__funcTakingVoid().(func(any) any)(_void)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any) any)(map[string]any{"$type": "String", "value": "Hello world!"})
	})
}
func main__funcTakingVoid() any {
	return func(_void any) any {
		return _void
	}
}
`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	mainPackage := "main"
	generatedProgram := codegen_golang.GenerateProgramMain(&codeIR, ir.Reference{
		Name: ir.VariableName(&mainPackage, "app"),
	})
	generated := generatedProgram.PackageCode + "\n\n" +
		generatedProgram.ImportsCode + "\n\n" +
		generatedProgram.UserspaceCode
	generatedFormatted := golang.Fmt(t, generated)
	assert.Equal(t, expectedGoCode, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generatedProgram.String())
	assert.Equal(t, expectedRunResult, output)
}
