package desugar

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
)

func TestDesugarShortCircuitExplicit(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ShortCircuitExplicit)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  when stringOrInt() {
    is str: Int => {
      str
    }
    is str: String => {
      join(str, "!")
    }
  }
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarShortCircuitInferLeft(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ShortCircuitInferLeft)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  when stringOrInt() {
    is str: Int => {
      str
    }
    other str => {
      join(str, "!")
    }
  }
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarShortCircuitInferRight(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ShortCircuitInferRight)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  when stringOrInt() {
    is str: String => {
      join(str, "!")
    }
    other str => {
      str
    }
  }
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarShortCircuitTwice(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ShortCircuitTwice)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  when stringOrInt() {
    is str: String => {
      when stringOrInt() {
        is strAgain: Int => {
          strAgain
        }
        other strAgain => {
          join(str, strAgain)
        }
      }
    }
    other str => {
      str
    }
  }
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarArrowInvocationOneArg(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrowInvocationOneArg)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main


f := (str: String): String => {
  str
}

usage := (): String => {
  str := "foo"
  f(str)
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarArrowInvocationOneArgChain(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrowInvocationOneArgChain)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main


f := (str: String): String => {
  str
}

g := (str: String): String => {
  str
}

h := (str: String): String => {
  str
}

usage := (): String => {
  str := "foo"
  h(g(f(str)))
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarArrowInvocationTwoArg(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrowInvocationTwoArg)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main


f := (str: String, str2: String): String => {
  str
}

usage := (): String => {
  str := "foo"

  str2 := "foo"
  f(str, str2)
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarArrowInvocationThreeArg(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrowInvocationThreeArg)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main


f := (str: String, str2: String, str3: String): String => {
  str
}

usage := (): String => {
  str := "foo"

  str2 := "foo"

  str3 := "foo"
  f(str, str2, str3)
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarArrowInvocationFunctions(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ArrowInvocationFunctions)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main


struct Stringer(
  produce: () ~> String,
  take1: (String) ~> String,
  take2: (String, String) ~> String,
  new: (String) ~> Stringer,
  consume: (String) ~> Void
)

usage := (s: Stringer): Void => {
  take1 := s.take1

  take2 := s.take2

  new := s.new

  consume := s.consume
  consume(take2(new(take1(s.produce())).produce(), s.produce()))
}
`

	assert.Equal(t, expected, formatted)
}

func TestDesugarShortCircuitUnused(t *testing.T) {
	parsed, err := parser.ParseString(testcode.ShortCircuitUnused)
	assert.NoError(t, err)
	desugared, err := Desugar(*parsed)
	assert.NoError(t, err)
	result := ToParsed(desugared)
	formatted := formatter.DisplayFileTopLevel(result)

	expected := `package main

import tenecs.error.Error

maybeString := (): String | Error => {
  ""
}

usage := (): String | Error => {
  when maybeString() {
    is _unused_: Error => {
      ""
    }
    other _unused_ => {
      _unused_
    }
  }
}
`

	assert.Equal(t, expected, formatted)
}
