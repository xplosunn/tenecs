package parser_typer_test

import (
	"testing"
)

func TestImplementationWithInvalidType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

app := (arg: NonExistingType): A => implement A {
	public a := ""
}
`, "not found type: NonExistingType")
}

func TestImplementationWithConstructorWithSameNameAsVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

a := (): A => implement A {
	public a := ""
}
`, "duplicate variable 'a'")
}
