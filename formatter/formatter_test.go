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

import tenecs.os.Main
import tenecs.os.Runtime

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

import tenecs.os.Console
import tenecs.os.Main
import tenecs.os.Runtime

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

import tenecs.os.Main
import tenecs.os.Runtime

app := (): Main => implement Main {
  public main := (runtime: Runtime) => {
    if false {
      runtime.console.log("Hello world!")
    } else {
      runtime.console.log("Hello world!")
    }
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithIfElseIf(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramWithIfElseIf)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Main
import tenecs.os.Runtime

app := (): Main => implement Main {
  public main := (runtime: Runtime) => {
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
}
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramWithVariableWithFunctionWithTypeInferred)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Main
import tenecs.os.Runtime

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

func TestDisplayArrayVariableWithEmptyArray(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrayVariableWithEmptyArray)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


noStrings := [String]()
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayArrayVariableWithTwoElementArray(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrayVariableWithTwoElementArray)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


someStrings := [String]("a", "b")
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayOrVariableWithEmptyArray(t *testing.T) {
	parsed, err := parser.ParseString(testcode.OrVariableWithEmptyArray)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


empty := [String | Boolean]()
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayOrVariableWithTwoElementArray(t *testing.T) {
	parsed, err := parser.ParseString(testcode.OrVariableWithTwoElementArray)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


hasStuff := [Boolean | String]("first", false)
`
	assert.Equal(t, expected, formatted)
}

func TestDisplayBasicTypeTrue(t *testing.T) {
	parsed, err := parser.ParseString(testcode.BasicTypeTrue)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


value := true
`
	assert.Equal(t, expected, formatted)
}

func TestWhenExplicitExhaustive(t *testing.T) {
	parsed, err := parser.ParseString(testcode.WhenExplicitExhaustive)
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
	parsed, err := parser.ParseString(testcode.WhenOtherMultipleTypes)
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
	parsed, err := parser.ParseString(testcode.GenericIO)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package mypackage


interface IO<A> {
  public run: () -> A
  public map: <B>((A) -> B) -> IO<B>
}

make := <A>(a: () -> A): IO<A> => implement IO<A> {
  public run := () => {
    a()
  }

  public map := <B>(f: (A) -> B): IO<B> => {
    make<B>(() => {
      f(a())
    })
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestMainProgramAnnotatedType(t *testing.T) {
	parsed, err := parser.ParseString(testcode.MainProgramAnnotatedType)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main.program

import tenecs.os.Main
import tenecs.os.Runtime

app: () -> Main = () => implement Main {
  public main := (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestImplementationWithAnnotatedVariable(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ImplementationWithAnnotatedVariable)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main


interface A {
  public a: () -> String
}

app := (): A => implement A {
  public a: () -> String = () => {
    ""
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestImportAliasMain(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ImportAliasMain)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `package main

import tenecs.os.Main as App
import tenecs.os.Runtime as Lib
import tenecs.string.join as concat

app := implement App {
  public main := (runtime: Lib) => {
    runtime.console.log(concat("Hello ", "world!"))
  }
}
`
	assert.Equal(t, expected, formatted)
}

func TestWhenAnnotatedVariable(t *testing.T) {
	parsed, err := parser.ParseString(testcode.WhenAnnotatedVariable)
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

func TestWFunctionCallToSplitArgumentsAcrossLines(t *testing.T) {
	parsed, err := parser.ParseString(`package main

func := (f: () -> String, g: () -> String): Void => {}

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


func := (f: () -> String, g: () -> String): Void => {}

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
import /* 5 */ tenecs.array.append // 6


str /* 7 */ := /* 8 */ "valueWithNoTypeAnnotation" // 9

struct /* 10 */ Post /* 11 */ (/* 12 */ title /* 13 */ : /* 14 */ String /* 15 */, author: String /* 16 */) // 17

interface /* 18 */ HelloWorld /* 19 */ { // 20
  public /* 21 */ inUppercase /* 22 */ : String // 23
  public inLowercase: String // 24
}
`)
	assert.NoError(t, err)
	formatted := formatter.DisplayFileTopLevel(*parsed)
	expected := `// 1
/* 2 */
package main

// 3
// 4
/* 5 */
import tenecs.array.append

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

// 17
/* 18 */
/* 19 */
interface HelloWorld {
  // 20
  /* 21 */
  /* 22 */
  public inUppercase: String
  // 23
  public inLowercase: String
  // 24
}
`
	assert.Equal(t, expected, formatted)
}
