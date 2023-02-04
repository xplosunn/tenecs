package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestMainProgramWithSingleExpression(t *testing.T) {
	validProgram(t, testcode.MainProgramWithSingleExpression)
}

func TestMainProgramWithInnerFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithInnerFunction)
}

func TestMainProgramWithWrongArgCount(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: Runtime, anotherRuntime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "expected same number of arguments as interface variable (1) but found 2")
}

func TestMainProgramWithVariableWithFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunction)
}

func TestMainProgramWithVariableWithFunctionTakingFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionTakingFunction)
}

func TestMainProgramWithVariableWithFunctionTakingFunctionFromStdLib(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionTakingFunctionFromStdLib1)
	validProgram(t, testcode.MainProgramWithVariableWithFunctionTakingFunctionFromStdLib2)
}

func TestMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionWithTypeInferred)

	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: Runtime) => {
		applyToString := (f: (String) -> Void, strF: () -> String): Void => {
			f(strF())
		}
		output := (): String => {
			"Hello world!"
		}
		applyToString(runtime.console.log, () => {false})
	}
}
`, "expected type String but found Boolean")
}

func TestMainProgramWithVariableWithFunctionWithWrongType(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: Runtime) => {
		output := (): String => {
			
		}
		runtime.console.log(output())
	}
}
`, "Function has return type of String but has empty body")
}

func TestMainProgramWithAnotherFunctionTakingConsole(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingConsole)
}

func TestMainProgramWithAnotherFunctionTakingConsoleAndMessage(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingConsoleAndMessage)
}

func TestMainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction)
}

func TestMainProgramWithArgAnnotatedArg(t *testing.T) {
	validProgram(t, testcode.MainProgramWithArgAnnotatedArg)
}

func TestMainProgramWithArgAnnotatedWrongArg(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: String) => {
		runtime.console.log("Hello world!")
	}
}
`, "in parameter position 0 expected type tenecs.os.Runtime but you have annotated String")
}

func TestMainProgramWithArgAnnotatedReturn(t *testing.T) {
	validProgram(t, testcode.MainProgramWithArgAnnotatedReturn)
}

func TestMainProgramWithArgAnnotatedWrongReturn(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): String => {
		runtime.console.log("Hello world!")
	}
}
`, "in return type expected type Void but you have annotated String")
}

func TestMainProgramWithArgAnnotatedArgAndReturn(t *testing.T) {
	validProgram(t, testcode.MainProgramWithArgAnnotatedArgAndReturn)
}

func TestMainProgramWithAnotherFunctionTakingRuntime(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingRuntime)
}
