package typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestUntypecheck(t *testing.T) {
	program := `package org.my.example


struct MyStruct<B>(
  a: String,
  b: B,
  err: Void,
  list: List<String>,
  f: <A>(String, A) ~> String
)

functionGeneric := <A>(a: A): A => {
  a
}

functionGenericUsed := <A>(a: A, myStruct: MyStruct<A>): A => {
  myStruct.f("", a)
  functionGeneric<A>(a)
}

functionIf := (): Boolean => {
  if true {
    true
  } else if false {
    false
  } else {
    true
  }
}

functionSimple := (): Void => {
  result := null
  result
}

functionSimpleUsage := (): Void => {
  functionSimple()
}

functionWhen := <T>(input: T | String | Boolean): String => {
  when input {
    is Boolean => {
      "<boolean>"
    }
    is s: String => {
      s
    }
    other => {
      "<unknown>"
    }
  }
}

functionWhenShortCircuit := (f: () ~> String | Boolean): String | Boolean => {
  a: String ?= f()
  b :? Boolean = f()
  c := f()
  c
}

literalInt: Int = 1

literalStr := "string literal"

literalStrList := []("")
`

	expected := `package org.my.example


struct MyStruct<B>(
  a: String,
  b: B,
  err: Void,
  list: List<String>,
  f: <A>(String, A) ~> String
)

functionGeneric: <A>(A) ~> A = <A>(a: A): A => {
  a
}

functionGenericUsed: <A>(A, MyStruct<A>) ~> A = <A>(a: A, myStruct: MyStruct<A>): A => {
  myStruct.f<A>("", a)
  functionGeneric<A>(a)
}

functionIf: () ~> Boolean = (): Boolean => {
  if true {
    true
  } else {
    if false {
      false
    } else {
      true
    }
  }
}

functionSimple: () ~> Void = (): Void => {
  result: Void = null
  result
}

functionSimpleUsage: () ~> Void = (): Void => {
  functionSimple()
}

functionWhen: <T>(T | String | Boolean) ~> String = <T>(input: T | String | Boolean): String => {
  when input {
    is Boolean => {
      "<boolean>"
    }
    is s: String => {
      s
    }
    other => {
      "<unknown>"
    }
  }
}

functionWhenShortCircuit: (() ~> String | Boolean) ~> String | Boolean = (f: () ~> String | Boolean): String | Boolean => {
  when f() {
    is a: String => {
      when f() {
        is b: Boolean => {
          b
        }
        other b => {
          c: String | Boolean = f()
          c
        }
      }
    }
    other a => {
      a
    }
  }
}

literalInt: Int = 1

literalStr: String = "string literal"

literalStrList: List<String> = [String]("")
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	untyped := typer.Untypecheck(typed)
	formatted := formatter.DisplayFileTopLevelIgnoringComments(untyped)
	assert.Equal(t, expected, formatted)
}
