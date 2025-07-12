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
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics any, _runtime any) any {
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "Hello world!"})
	})
}

func tenecs_go__Main() any {
	log := func(generics any, msg any) any {
		println(msg.(map[string]any)["value"].(string))
		return nil
	}
	console := map[string]any{
		"_log": log,
	}
	refCreator := map[string]any{
		"_new": func(generics any, value any) any {
			var ref any = value
			return map[string]any{
				"_get": func(generics any) any {
					return ref
				},
				"_set": func(generics any, value any) any {
					ref = value
					return nil
				},
				"_modify": func(generics any, f any) any {
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
	return func(generics any, run any) any {
		return run.(func(any, any) any)(generics, runtime)
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
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics any, _runtime any) any {
		_ref := _runtime.(map[string]any)["_ref"].(map[string]any)["_new"].(func(any, any) any)([]any{"String"}, map[string]any{"$type": "String", "value": "hello"})
		_ = _ref
		_runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _ref.(map[string]any)["_get"].(func(any) any)([]any{}))
		_ref.(map[string]any)["_set"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "world"})
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _ref.(map[string]any)["_get"].(func(any) any)([]any{}))
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
	return tenecs_go__Main().(func(any, any) any)([]any{}, main__helloWorld())
}
func main__helloWorld() any {
	return func(generics any, _runtime any) any {
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "Hello world!"})
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
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics any, _runtime any) any {
		_void := func() any {
			_nestedVar := map[string]any{"$type": "Int", "value": 1}
			_ = _nestedVar
			return nil
		}
		_ = _void
		main__funcTakingVoid().(func(any, any) any)([]any{}, _void)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "Hello world!"})
	})
}
func main__funcTakingVoid() any {
	return func(generics any, _void any) any {
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

func TestGenerateProgramMainTopLevelFunctionAssignedToLocalVariable(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

firstString := (str1: String, str2: String): String => {
  str1
}

app := Main((runtime: Runtime): Void => {
  first := firstString
  hello := first("hello", "world")
  runtime.console.log(hello)
})
`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics any, _runtime any) any {
		_first := main__firstString()
		_ = _first
		_hello := _first.(func(any, any, any) any)([]any{}, map[string]any{"$type": "String", "value": "hello"}, map[string]any{"$type": "String", "value": "world"})
		_ = _hello
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _hello)
	})
}
func main__firstString() any {
	return func(generics any, _str1 any, _str2 any) any {
		return _str1
	}
}
`

	expectedRunResult := "hello\n"

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

func TestGenerateProgramMainWhen(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

isStringOrInt := (arg: Int | String): String => {
  when arg {
	is Int => { "is int" }
	is String => { "is string" }
  }
}

app := Main((runtime: Runtime): Void => {
  expectIsInt := isStringOrInt(1)
  expectIsString := isStringOrInt("")
  runtime.console.log(expectIsInt)
  runtime.console.log(expectIsString)
})
`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics any, _runtime any) any {
		_expectIsInt := main__isStringOrInt().(func(any, any) any)([]any{}, map[string]any{"$type": "Int", "value": 1})
		_ = _expectIsInt
		_expectIsString := main__isStringOrInt().(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": ""})
		_ = _expectIsString
		_runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsInt)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsString)
	})
}
func main__isStringOrInt() any {
	return func(generics any, _arg any) any {
		return func() any {
			__over := _arg
			_ = __over
			if __over.(map[string]any)["$type"] == "Int" {
				return map[string]any{"$type": "String", "value": "is int"}
			} else {
				if __over.(map[string]any)["$type"] == "String" {
					return map[string]any{"$type": "String", "value": "is string"}
				} else {
					return nil
				}
			}
		}()
	}
}
`

	expectedRunResult := "is int\nis string\n"

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

func TestGenerateProgramMainWhenOtherGeneric(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

isString := <T>(arg: T | String): String => {
  when arg {
	is String => { "is string" }
	other => { "is other" }
  }
}

app := Main((runtime: Runtime): Void => {
  expectIsString := isString("")
  expectIsOther := isString(1)
  runtime.console.log(expectIsString)
  runtime.console.log(expectIsOther)
})
`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics any, _runtime any) any {
		_expectIsString := main__isString().(func(any, any) any)([]any{"String"}, map[string]any{"$type": "String", "value": ""})
		_ = _expectIsString
		_expectIsOther := main__isString().(func(any, any) any)([]any{"Int"}, map[string]any{"$type": "Int", "value": 1})
		_ = _expectIsOther
		_runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsString)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsOther)
	})
}
func main__isString() any {
	return func(generics any, _arg any) any {
		return func() any {
			__over := _arg
			_ = __over
			if __over.(map[string]any)["$type"] == "String" {
				return map[string]any{"$type": "String", "value": "is string"}
			} else {
				return map[string]any{"$type": "String", "value": "is other"}
			}
		}()
	}
}
`

	expectedRunResult := "is string\nis other\n"

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
