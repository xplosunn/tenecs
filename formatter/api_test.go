package formatter_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestTestCode(t *testing.T) {
	for _, code := range testcode.GetAll() {
		t.Run(code.Name, func(t *testing.T) {
			parsed, err := parser.ParseString(code.Content)
			assert.NoError(t, err)
			formatted := formatter.DisplayFileTopLevel(*parsed)
			assert.Equal(t, code.Content, formatted)
		})
	}
}

func TestDisplayMainProgramWithSingleExpression(t *testing.T) {
	code := `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => runtime.console.log("Hello world!")
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithAnotherFunctionTakingConsole(t *testing.T) {
	code := `
package main

import tenecs.go.Main
import tenecs.go.Runtime
import tenecs.go.Console

app := Main(
  main = (runtime) => {
    mainRun(runtime.console)
  }
)

mainRun := (console: Console): Void => {
  console.log("Hello world!")
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Console
import tenecs.go.Main
import tenecs.go.Runtime

app := Main(
  main = (runtime) => {
    mainRun(runtime.console)
  }
)

mainRun := (console: Console): Void => {
  console.log("Hello world!")
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithIfElse(t *testing.T) {
	code := `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    if false {
      runtime.console.log("Hello world!")
    } else {
      runtime.console.log("Hello world!")
    }
  }
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main(
  main = (runtime: Runtime) => {
    if false {
      runtime.console.log("Hello world!")
    } else {
      runtime.console.log("Hello world!")
    }
  }
)
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithIfElseIf(t *testing.T) {
	code := `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    if false {
      runtime.console.log("Hello world!")
    } else if false {
      runtime.console.log("Hello world!")
    } else if true {
      runtime.console.log("Hello world!")
    } else {
      runtime.console.log("Hello world!")
    }
  }
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main(
  main = (runtime: Runtime) => {
    if false {
      runtime.console.log("Hello world!")
    } else if false {
      runtime.console.log("Hello world!")
    } else if true {
      runtime.console.log("Hello world!")
    } else {
      runtime.console.log("Hello world!")
    }
  }
)
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	code := `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  (runtime: Runtime) => {
    applyToString := (f: (String) ~> Void, strF: () ~> String): Void => {
      f(strF())
    }
    applyToString(runtime.console.log, () => {"Hello World!"})
  }
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  applyToString := (f: (String) ~> Void, strF: () ~> String): Void => {
    f(strF())
  }
  applyToString(runtime.console.log, () => {
    "Hello World!"
  })
})
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayGenericStructInstance1(t *testing.T) {
	code := `
package main

import tenecs.go.Main

struct Box<T>(inside: T)

app := Main(
  main = (runtime) => {
    box := Box<String>("Hello world!")
    runtime.console.log(box.inside)
  }
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Main

struct Box<T>(
  inside: T
)

app := Main(
  main = (runtime) => {
    box := Box<String>("Hello world!")
    runtime.console.log(box.inside)
  }
)
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayListVariableWithEmptyList(t *testing.T) {
	code := `
package main

noStrings := [ String ] ( )
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


noStrings := [String]()
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayListVariableWithTwoElementList(t *testing.T) {
	code := `
package main

someStrings := [ String ] ( "a" , "b" )
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


someStrings := [String]("a", "b")
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayOrVariableWithEmptyList(t *testing.T) {
	code := `
package main

empty := [ String | Boolean ] ( )
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


empty := [String | Boolean]()
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayOrVariableWithTwoElementList(t *testing.T) {
	code := `
package main

hasStuff := [ Boolean | String ] ( "first", false )
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


hasStuff := [Boolean | String]("first", false)
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayBasicTypeTrue(t *testing.T) {
	code := `
package main

value := true
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


value := true
`
	assert.Equal(t, expected, formatted)
}

func TestWhenExplicitExhaustive(t *testing.T) {
	code := `
package main

asString := (arg: Boolean | String): String => {
  when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    is b: String => {
      b
    }
  }
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


asString := (arg: Boolean | String): String => {
  when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    is b: String => {
      b
    }
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestWhenOtherMultipleTypes(t *testing.T) {
	code := `
package main

yeetString := (arg: Boolean | String | Void): Boolean | Void => {
  when arg {
    is String => {
      false
    }
    other a => {
      a
    }
  }
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


yeetString := (arg: Boolean | String | Void): Boolean | Void => {
  when arg {
    is String => {
      false
    }
    other a => {
      a
    }
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestGenericIO(t *testing.T) {
	code := `
package mypackage

struct IO<A>(
  run: () ~> A,
  _map: <B>((A) ~> B) ~> IO<B>
)

make := <A>(a: () ~> A): IO<A> => {
  IO<A>(
    run = () => {
      a()
    },
    _map = <B>(f: (A) ~> B): IO<B> => {
      make<B>(() => { f(a()) })
    }
  )
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package mypackage


struct IO<A>(
  run: () ~> A,
  _map: <B>((A) ~> B) ~> IO<B>
)

make := <A>(a: () ~> A): IO<A> => {
  IO<A>(
    run = () => {
      a()
    },
    _map = <B>(f: (A) ~> B): IO<B> => {
      make<B>(() => {
        f(a())
      })
    }
  )
}
`
	assert.Equal(t, expected, formatted)
}

func TestMainProgramAnnotatedType(t *testing.T) {
	code := `
package main.program

import tenecs.go.Runtime
import tenecs.go.Main

app: Main = Main(
  main = (runtime: Runtime) => runtime.console.log("Hello world!")
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main.program

import tenecs.go.Main
import tenecs.go.Runtime

app: Main = Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)
`
	assert.Equal(t, expected, formatted)
}

func TestImportAliasMain(t *testing.T) {
	code := `
package main

import tenecs.go.Runtime as Lib
import tenecs.go.Main as App
import tenecs.string.join as concat

app := App(
  main = (runtime: Lib) => {
    runtime.console.log(concat("Hello ", "world!"))
  }
)
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.go.Main as App
import tenecs.go.Runtime as Lib
import tenecs.string.join as concat

app := App(
  main = (runtime: Lib) => {
    runtime.console.log(concat("Hello ", "world!"))
  }
)
`
	assert.Equal(t, expected, formatted)
}

func TestWhenAnnotatedVariable(t *testing.T) {
	code := `
package main

asString := (arg: Boolean | String): String => {
  result: String = when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    is s: String => {
      s
    }
  }
  result
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


asString := (arg: Boolean | String): String => {
  result: String = when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    is s: String => {
      s
    }
  }
  result
}
`
	assert.Equal(t, expected, formatted)
}

func TestGenericsInferTypeParameterPartialLeft(t *testing.T) {
	code := `
package main

pickRight := <L, R>(left: L, right: R): R => {
  right
}

usage := (): Void => {
  str := pickRight<_, String>("", "")
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


pickRight := <L, R>(left: L, right: R): R => {
  right
}

usage := (): Void => {
  str := pickRight<_, String>("", "")
}
`
	assert.Equal(t, expected, formatted)
}

func TestWFunctionCallToSplitArgumentsAcrossLines(t *testing.T) {
	parsed, err := parser.ParseString(`package main

func := (f: () ~> String, g: () ~> String): Void => {}

usage := (): Void => {
  helloWorld := (): String => { "hello world" }
  doNotSplit := func(helloWorld, helloWorld)
  alsoDoNotSplit := func(helloWorld, (): String => { "foo" })
  split := func((): String => { "foo" }, helloWorld)
  alsoSplit := func((): String => { "foo" }, (): String => { "foo" })
}

`)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


func := (f: () ~> String, g: () ~> String): Void => {}

usage := (): Void => {
  helloWorld := (): String => {
    "hello world"
  }

  doNotSplit := func(helloWorld, helloWorld)

  alsoDoNotSplit := func(helloWorld, (): String => {
    "foo"
  })

  split := func(
    (): String => {
      "foo"
    },
    helloWorld
  )

  alsoSplit := func(
    (): String => {
      "foo"
    },
    (): String => {
      "foo"
    }
  )
}
`
	assert.Equal(t, expected, formatted)
}

func TestComments(t *testing.T) {
	parsed, err := parser.ParseString(`// 1
package /* 2 */ main // 3
// 4
import /* 5 */ tenecs.list.append // 6


str /* 7 */ := /* 8 */ "valueWithNoTypeAnnotation" // 9

struct /* 10 */ Post /* 11 */ (/* 12 */ title /* 13 */ : /* 14 */ String /* 15 */, author: String /* 16 */) // 17

`)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `// 1
/* 2 */
package main

// 3
// 4
/* 5 */
import tenecs.list.append

// 6
/* 7 */
/* 8 */
str := "valueWithNoTypeAnnotation"

// 9
/* 10 */
struct Post(
  /* 11 */
  /* 12 */
  /* 13 */
  /* 14 */
  title: String,
  /* 15 */
  author: String
  /* 16 */
)
`
	assert.Equal(t, expected, formatted)
}

func TestShortcircuit(t *testing.T) {
	code := `package main


stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  strOne: String ? Int = stringOrInt()

  strTwo :? Int = stringOrInt()

  strThree: String ?= stringOrInt()

  willNotCompileButShouldFormat :?= stringOrInt()
  stringOrInt()
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	assert.Equal(t, code, formatted)
}

func TestNamedArgument(t *testing.T) {
	code := `package main


f := (a: String, b: String): String => {
  a
}

usage := (): String => {
  f("", "")
  f(
    a = "",
    ""
  )
  f(
    "",
    b = ""
  )
  f(
    a = "",
    b = ""
  )
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	assert.Equal(t, code, formatted)
}

func TestNArrowInvocationOneArg(t *testing.T) {
	code := `package main


f := (str: String): String => {
  str
}

usage := (): Void => {
  str := "foo"
  str->f()
}
`
	parsed, err := parser.ParseString(code)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	assert.Equal(t, code, formatted)
}
