package parser_typer_test

import (
	"testing"
)

func TestStructWithConstructorWithUnknownType(t *testing.T) {
	invalidProgram(t, `
package main

struct InvalidRecord(a: Unknown)
`, "not found type: Unknown")
}
