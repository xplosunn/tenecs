package parser_typer_test

import (
	"testing"
)

func TestStructWithConstructorWithUnknownType(t *testing.T) {
	invalidProgram(t, `
package main

struct InvalidRecord(a: Unknown)
`, "not found type: Unknown (are you using an incomparable type?)")
}

func TestStructWithConstructorWithInterface(t *testing.T) {
	invalidProgram(t, `
package main

interface A {}

struct InvalidRecord(a: A)
`, "not found type: A (are you using an incomparable type?)")
}
