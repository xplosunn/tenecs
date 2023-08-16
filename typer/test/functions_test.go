package parser_typer_test

import (
	"testing"
)

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
`, "expected 1 params but got 2")
}

func TestInvalidMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
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
`, "empty block only allowed for Void type")
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

func TestWrongReturn(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.array.append
import tenecs.array.length
import tenecs.http.Server
import tenecs.http.newServer
import tenecs.int.plus
import tenecs.json.field
import tenecs.json.parseBoolean
import tenecs.json.parseObject2
import tenecs.json.parseString
import tenecs.os.Main
import tenecs.ref.Ref
import tenecs.ref.RefCreator

struct Todo(
  id: Int,
  title: String,
  done: Boolean
)

app := implement Main {
  public main := (runtime) => {
    runtime.console.log("Starting demo todo server")

    state := runtime.ref.new([Todo]())

    server := setupServer(runtime.ref, state)

    error := server.serve("localhost:8081", runtime.execution.blocker)
    runtime.console.log(error.message)
  }
}

setupServer := (refCreator: RefCreator, state: Ref<Array<Todo>>): Server => {
  todoParser := parseObject2(
    (title: String, done: Boolean): Todo => {
      Todo(plus(length(state.get()), 1), title, done)
    },
    field("title", parseString()),
    field("done", parseBoolean())
  )

  server := newServer(refCreator)
  server.restHandlerGet<Array<Todo>>("/todo", (responseStatusRef) => {
    state.get()
  })
}
`, "Expected tenecs.http.Server but got Void")
}
