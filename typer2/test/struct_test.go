package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestStructWithConstructorEmpty(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorEmpty)
}

func TestStructWithConstructorWithString(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorWithString)
}

func TestStructWithConstructorWithBooleans(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorWithBooleans)
}

func TestStructWithConstructorAnotherStruct(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorAnotherStruct1)
	validProgram(t, testcode.StructWithConstructorAnotherStruct2)
}

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

func TestStructAsVariable(t *testing.T) {
	validProgram(t, testcode.StructAsVariable)
}
