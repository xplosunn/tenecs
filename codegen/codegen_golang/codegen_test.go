package codegen_golang_test

import (
	"fmt"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"

	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "Hello world!"})
	})
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _ref any
		_ref = _runtime.(map[string]any)["_ref"].(map[string]any)["_new"].(func(any, any) any)([]any{[]any{"String"}}, map[string]any{"$type": "String", "value": "hello"})
		_ = _ref
		_ = _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _ref.(map[string]any)["_get"].(func(any) any)([]any{}))
		_ = _ref.(map[string]any)["_set"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "world"})
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, main__helloWorld())
}
func main__helloWorld() any {
	return func(generics1 any, _runtime any) any {
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _void any
		_void = func() any {
			var _nestedVar any
			_nestedVar = map[string]any{"$type": "Int", "value": 1}
			_ = _nestedVar
			return nil
		}
		_ = _void
		_ = main__funcTakingVoid().(func(any, any) any)([]any{}, _void)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "Hello world!"})
	})
}
func main__funcTakingVoid() any {
	return func(generics1 any, _void any) any {
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _first any
		_first = main__firstString()
		_ = _first
		var _hello any
		_hello = _first.(func(any, any, any) any)([]any{}, map[string]any{"$type": "String", "value": "hello"}, map[string]any{"$type": "String", "value": "world"})
		_ = _hello
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _hello)
	})
}
func main__firstString() any {
	return func(generics1 any, _str1 any, _str2 any) any {
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _expectIsInt any
		_expectIsInt = main__isStringOrInt().(func(any, any) any)([]any{}, map[string]any{"$type": "Int", "value": 1})
		_ = _expectIsInt
		var _expectIsString any
		_expectIsString = main__isStringOrInt().(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": ""})
		_ = _expectIsString
		_ = _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsInt)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsString)
	})
}
func main__isStringOrInt() any {
	return func(generics1 any, _arg any) any {
		return func() any {
			var __over any
			__over = _arg
			_ = __over
			if _, ok := any("Int").([]any); ok && slices.Contains(any("Int").([]any), __over.(map[string]any)["$type"]) || (!ok && ("Int" == __over.(map[string]any)["$type"])) {
				return map[string]any{"$type": "String", "value": "is int"}
			} else {
				if _, ok := any("String").([]any); ok && slices.Contains(any("String").([]any), __over.(map[string]any)["$type"]) || (!ok && ("String" == __over.(map[string]any)["$type"])) {
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

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _expectIsString any
		_expectIsString = main__isString().(func(any, any) any)([]any{[]any{"String"}}, map[string]any{"$type": "String", "value": ""})
		_ = _expectIsString
		var _expectIsOther any
		_expectIsOther = main__isString().(func(any, any) any)([]any{[]any{"Int"}}, map[string]any{"$type": "Int", "value": 1})
		_ = _expectIsOther
		_ = _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsString)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _expectIsOther)
	})
}
func main__isString() any {
	return func(generics1 any, _arg any) any {
		return func() any {
			var __over any
			__over = _arg
			_ = __over
			if _, ok := any("String").([]any); ok && slices.Contains(any("String").([]any), __over.(map[string]any)["$type"]) || (!ok && ("String" == __over.(map[string]any)["$type"])) {
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

func TestGenerateProgramMainWhenGeneric(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

leftOrRight := <L, R>(arg: L | R): String => {
  when arg {
	is L => { "is left" }
	is R => { "is right" }
  }
}

app := Main((runtime: Runtime): Void => {
  left := leftOrRight<String, Int>("")
  right := leftOrRight<String, Int>(1)
  runtime.console.log(left)
  runtime.console.log(right)
})
`
	expectedGoCode := `package main

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _left any
		_left = main__leftOrRight().(func(any, any) any)([]any{[]any{"String"}, []any{"Int"}}, map[string]any{"$type": "String", "value": ""})
		_ = _left
		var _right any
		_right = main__leftOrRight().(func(any, any) any)([]any{[]any{"String"}, []any{"Int"}}, map[string]any{"$type": "Int", "value": 1})
		_ = _right
		_ = _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _left)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _right)
	})
}
func main__leftOrRight() any {
	return func(generics1 any, _arg any) any {
		return func() any {
			var __over any
			__over = _arg
			_ = __over
			if _, ok := any(generics1.([]any)[0]).([]any); ok && slices.Contains(any(generics1.([]any)[0]).([]any), __over.(map[string]any)["$type"]) || (!ok && (generics1.([]any)[0] == __over.(map[string]any)["$type"])) {
				return map[string]any{"$type": "String", "value": "is left"}
			} else {
				if _, ok := any(generics1.([]any)[1]).([]any); ok && slices.Contains(any(generics1.([]any)[1]).([]any), __over.(map[string]any)["$type"]) || (!ok && (generics1.([]any)[1] == __over.(map[string]any)["$type"])) {
					return map[string]any{"$type": "String", "value": "is right"}
				} else {
					return nil
				}
			}
		}()
	}
}
`

	expectedRunResult := "is left\nis right\n"

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

func TestGenerateProgramMainNestedWhenGeneric(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

leftOrRight := <L, R>(arg: L | R): String => {
  run := (): String => {
    when arg {
	  is L => { "is left" }
	  is R => { "is right" }
    }
  }
  run()
}

app := Main((runtime: Runtime): Void => {
  left := leftOrRight<String, Int>("")
  right := leftOrRight<String, Int>(1)
  runtime.console.log(left)
  runtime.console.log(right)
})
`
	expectedGoCode := `package main

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		var _left any
		_left = main__leftOrRight().(func(any, any) any)([]any{[]any{"String"}, []any{"Int"}}, map[string]any{"$type": "String", "value": ""})
		_ = _left
		var _right any
		_right = main__leftOrRight().(func(any, any) any)([]any{[]any{"String"}, []any{"Int"}}, map[string]any{"$type": "Int", "value": 1})
		_ = _right
		_ = _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _left)
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, _right)
	})
}
func main__leftOrRight() any {
	return func(generics1 any, _arg any) any {
		var _run any
		_run = func(generics2 any) any {
			return func() any {
				var __over any
				__over = _arg
				_ = __over
				if _, ok := any(generics1.([]any)[0]).([]any); ok && slices.Contains(any(generics1.([]any)[0]).([]any), __over.(map[string]any)["$type"]) || (!ok && (generics1.([]any)[0] == __over.(map[string]any)["$type"])) {
					return map[string]any{"$type": "String", "value": "is left"}
				} else {
					if _, ok := any(generics1.([]any)[1]).([]any); ok && slices.Contains(any(generics1.([]any)[1]).([]any), __over.(map[string]any)["$type"]) || (!ok && (generics1.([]any)[1] == __over.(map[string]any)["$type"])) {
						return map[string]any{"$type": "String", "value": "is right"}
					} else {
						return nil
					}
				}
			}()
		}
		_ = _run
		return _run.(func(any) any)([]any{})
	}
}
`

	expectedRunResult := "is left\nis right\n"

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

func TestGenerateProgramMainIf(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

getMessage := (left: Boolean, leftMessage: String, rightMessage: String): String => {
  if left {
    leftMessage
  } else {
    rightMessage
  }
}

app := Main((runtime: Runtime): Void => {
  left := getMessage(true, "left", "right")
  right := getMessage(false, "left", "right")
  runtime.console.log(left)
  runtime.console.log(right)
})
`

	expectedRunResult := "left\nright\n"

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
	snaps.MatchStandaloneSnapshot(t, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generatedProgram.String())
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateProgramMainWhenList(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

match := (list: List<String> | List<Boolean>): String => {
  when list {
    is List<String> => { "string" }
    is List<Boolean> => { "boolean" }
  }
}

app := Main((runtime: Runtime): Void => {
  runtime.console.log(match(<String>[]))
  runtime.console.log(match(<Boolean>[]))
  runtime.console.log(match(["a", "b", "c"]))
  runtime.console.log(match([true, false, true]))
})
`

	expectedRunResult := "string\nboolean\nstring\nboolean\n"

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
	snaps.MatchStandaloneSnapshot(t, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generatedProgram.String())
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateProgramMainIgnoredString(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main((runtime: Runtime): Void => {
  "ignored string"
  runtime.console.log("output")
})
`
	expectedGoCode := `package main

import (
	"slices"
	"strings"
)

func main__app() any {
	return tenecs_go__Main().(func(any, any) any)([]any{}, func(generics1 any, _runtime any) any {
		_ = map[string]any{"$type": "String", "value": "ignored string"}
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any, any) any)([]any{}, map[string]any{"$type": "String", "value": "output"})
	})
}
`

	expectedRunResult := "output\n"

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

func TestGenerateProgramTestSuccess(t *testing.T) {
	program := `package test

import tenecs.test.UnitTest
import tenecs.test.UnitTestSuite
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.test.GoIntegrationTest
import tenecs.test.GoIntegrationTestKit
import tenecs.go.Runtime

helloWorld := (): String => {
  "hello world!"
}

_ := UnitTestSuite(
  "My Tests",
  tests = (registry: UnitTestRegistry): Void => {
    registry.test("hello world function", testCaseHelloworld)
  }
)

_ := UnitTest("unitHello", testCaseHelloworld)

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()
  expected := "hello world!"
  testkit.assert.equal<String>(result, expected)
}

_ := GoIntegrationTest(
  dependsOnVersionOf = "no_dependencies",
  "integration test",
  (testkit: GoIntegrationTestKit, runtime: Runtime): Void => {
    testkit.assert.equal<String>("hello world!", "hello world!")
  }
)
`
	expectedRunResult := fmt.Sprintf(`unit tests:
  %s unitHello
My Tests:
  %s hello world function
integration tests:
  %s integration test

Ran a total of 3 tests
  * 3 succeeded
  * 0 failed
`, "[\u001b[32mOK\u001b[0m]", "[\u001b[32mOK\u001b[0m]", "[\u001b[32mOK\u001b[0m]")

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	generated := codegen_golang.GenerateProgramTest(&codeIR, codegen.FindTests(typed))

	output := golang.RunCodeUnlessCached(t, generated.String())
	assert.Equal(t, expectedRunResult, output)

	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated.String()))
}

func TestGenerateProgramTestFail(t *testing.T) {
	program := `package test

import tenecs.test.UnitTest
import tenecs.test.UnitTestSuite
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.test.GoIntegrationTest
import tenecs.test.GoIntegrationTestKit
import tenecs.go.Runtime

helloWorld := (): String => {
  "hello world!"
}

_ := UnitTestSuite(
  "My Tests",
  tests = (registry: UnitTestRegistry): Void => {
    registry.test("hello world function", testCaseHelloworld)
  }
)

_ := UnitTest("unitHello", testCaseHelloworld)

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()
  expected := "hello world"
  testkit.assert.equal<String>(result, expected)
}

_ := GoIntegrationTest(
  dependsOnVersionOf = "no_dependencies",
  "integration test",
  (testkit: GoIntegrationTestKit, runtime: Runtime): Void => {
    testkit.assert.equal<String>("hello world!", "hello world")
  }
)
`
	expectedRunResult := fmt.Sprintf(`unit tests:
  %s unitHello
    @file.10x:27: {"$type":"String","value":"hello world!"} is not equal to {"$type":"String","value":"hello world"}
My Tests:
  %s hello world function
    @file.10x:27: {"$type":"String","value":"hello world!"} is not equal to {"$type":"String","value":"hello world"}
integration tests:
  %s integration test
    @file.10x:34: {"$type":"String","value":"hello world!"} is not equal to {"$type":"String","value":"hello world"}

Ran a total of 3 tests
  * 0 succeeded
  * 3 failed
`, "[\u001b[31mFAILURE\u001b[0m]", "[\u001b[31mFAILURE\u001b[0m]", "[\u001b[31mFAILURE\u001b[0m]")

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	generated := codegen_golang.GenerateProgramTest(&codeIR, codegen.FindTests(typed))

	output := golang.RunCodeUnlessCached(t, generated.String())
	assert.Equal(t, expectedRunResult, output)

	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated.String()))
}
