package parser_typer_test

import (
	"testing"
)

func TestImplementationWithInvalidType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
  a: () -> String
}

app := (arg: NonExistingType): A => implement A {
  a := () => ""
}
`, "not found type: NonExistingType")
}

func TestImplementationWithConstructorWithSameNameAsVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
  a: () -> String
}

a := (): A => implement A {
  a := () => ""
}
`, "duplicate variable 'a'")
}
