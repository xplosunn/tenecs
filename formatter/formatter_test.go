package formatter_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestDisplayMainProgramWithSingleExpression(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramWithSingleExpression)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
  public main := (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithAnotherFunctionTakingConsole(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramWithAnotherFunctionTakingConsole)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Main
import tenecs.os.Runtime
import tenecs.os.Console

app := (): Main => implement Main {
  public main := (runtime) => {
    mainRun(runtime.console)
  }

  mainRun := (console: Console): Void => {
    console.log("Hello world!")
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithIfElse(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramWithIfElse)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
  public main := (runtime: Runtime) => {
    if true {
      runtime.console.log("Hello world!")
    } else {
      runtime.console.log("Hello world!")
    }
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramWithVariableWithFunctionWithTypeInferred)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
  public main := (runtime: Runtime) => {
    applyToString := (f: (String) -> Void, strF: () -> String): Void => {
      f(strF())
    }
    applyToString(runtime.console.log, () => {
      "Hello World!"
    })
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayGenericStructInstance1(t *testing.T) {
	parsed, err := parser.ParseString(testcode.GenericStructInstance)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Main

struct Box<T>(
  inside: T
)

app := (): Main => implement Main {
  public main := (runtime) => {
    box := Box<String>("Hello world!")
    runtime.console.log(box.inside)
  }
}
`
	assert.Equal(t, expected, formatted)
}
