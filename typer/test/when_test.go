package parser_typer_test

import "testing"

func TestWhenFailOnStructWithPhantomType(t *testing.T) {
	program := `package example

struct Phantom<T>(str: String)

usage := (): String => {
  p := Phantom<Int>("hello")
  when p {
    is phantom: Phantom<Int> => {
      phantom.str
    } 
  }
}
`
	invalidProgram(t, program, "matching on a struct with generics requires the struct to have one field of that type")
}

func TestWhenFailOnGeneric(t *testing.T) {
	program := `package example

usage := <A, B>(arg: A): String => {
  when arg {
    is B => {
      "it is B"
    }
    other => {
      "other"
    }
  }
}
`
	invalidProgram(t, program, "can't match on generic")
}
