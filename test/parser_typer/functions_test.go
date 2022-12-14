package parser_typer_test

import "testing"

func TestMainProgramWithInnerFunction(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		go := (): Void => {
			runtime.console.log("Hello world!")	
		}
		go()
	}
}
`)
}

func TestMainProgramWithWrongArgCount(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime, anotherRuntime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "expected same number of arguments as interface variable (1) but found 2")
}

func TestMainProgramWithVariableWithFunction(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		output := (): String => {
			"Hello world!"
		}
		runtime.console.log(output())
	}
}
`)
}

func TestMainProgramWithVariableWithFunctionTakingFunction(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		output := (): String => {
			"Hello world!"
		}
		run := (f: () -> String): String => {
			f()
		}
		runtime.console.log(run(output))
	}
}
`)
}

func TestMainProgramWithVariableWithFunctionTakingFunctionFromStdLib(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		applyToString := (f: (String) -> Void, str: String): Void => {
			f(str)
		}
		output := (): String => {
			"Hello world!"
		}
		applyToString(runtime.console.log, output())
	}
}
`)

	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		applyToString := (f: (String) -> Void, strF: () -> String): Void => {
			f(strF())
		}
		output := (): String => {
			"Hello world!"
		}
		applyToString(runtime.console.log, output)
	}
}
`)
}

func TestMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		applyToString := (f: (String) -> Void, strF: () -> String): Void => {
			f(strF())
		}
		applyToString(runtime.console.log, () => {"Hello World!"})
	}
}
`)

	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
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

implementing Main module app {
	public main := (runtime: Runtime) => {
		output := (): String => {
			
		}
		runtime.console.log(output())
	}
}
`, "Function has return type of String but has empty body")
}

func TestMainProgramWithAnotherFunctionTakingConsole(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main
import tenecs.os.Runtime
import tenecs.os.Console

implementing Main module app {
	public main := (runtime) => {
		mainRun(runtime.console)
	}
	mainRun := (console: Console): Void => {
		console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithAnotherFunctionTakingConsoleAndMessage(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main
import tenecs.os.Runtime
import tenecs.os.Console

implementing Main module app {
	public main := (runtime) => {
		mainRun(runtime.console, "Hello world!")
	}
	mainRun := (console: Console, message: String): Void => {
		console.log(message)
	}
}
`)
}

func TestMainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main
import tenecs.os.Runtime
import tenecs.os.Console

implementing Main module app {
	public main := (runtime) => {
		mainRun(runtime.console, helloWorld())
	}
	mainRun := (console: Console, message: String): Void => {
		console.log(message)
	}
	helloWorld := (): String => {
		"Hello world!"
	}
}
`)
}

func TestMainProgramWithArgAnnotatedArg(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithArgAnnotatedWrongArg(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: String) => {
		runtime.console.log("Hello world!")
	}
}
`, "in parameter position 0 expected type tenecs.os.Runtime but you have annotated String")
}

func TestMainProgramWithArgAnnotatedReturn(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithArgAnnotatedWrongReturn(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): String => {
		runtime.console.log("Hello world!")
	}
}
`, "in return type expected type Void but you have annotated String")
}

func TestMainProgramWithArgAnnotatedArgAndReturn(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime): Void => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithAnotherFunctionTakingRuntime(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main
import tenecs.os.Runtime

implementing Main module app {
	public main := (runtime) => {
		mainRun(runtime)
	}
	mainRun := (runtime: Runtime): Void => {
		runtime.console.log("Hello world!")
	}
}
`)
}
