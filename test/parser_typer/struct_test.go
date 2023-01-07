package parser_typer_test

import "testing"

func TestStructWithConstructorEmpty(t *testing.T) {
	validProgram(t, `
package main

struct NoArgsStruct()
`)
}

func TestStructWithConstructorWithString(t *testing.T) {
	validProgram(t, `
package main

struct StringStruct(str: String)
`)
}

func TestStructWithConstructorWithBooleans(t *testing.T) {
	validProgram(t, `
package main

struct BooleanColor(red: Boolean, green: Boolean, blue: Boolean)
`)
}

func TestStructWithConstructorAnotherStruct(t *testing.T) {
	validProgram(t, `
package main

struct Outer(c: Inner)
struct Inner()
`)
	validProgram(t, `
package main

struct Inner()
struct Outer(c: Inner)
`)
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
	validProgram(t, `
package main

import tenecs.os.Main

struct Person(name: String)

implementing Main module app {
	public main := (runtime) => {
		me := Person("Author")
	}
}
`)
}
