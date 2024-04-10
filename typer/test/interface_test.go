package parser_typer_test

import (
	"testing"
)

func TestInterfaceVariablesSameName(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
  a: () -> String
  a: () -> String
}
`, "more than one variable with name 'a'")
}

func TestInterfaceWithSeparateImplementationVariableStringThatIsNotInInterface(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
  a: () -> String
}

app := (): A => implement A {
  a := () => ""
  b := () => ""
}
`, "variable b is not part of interface being implemented")
}

func TestInterfaceWithSeparateImplementationMissingVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
  a: () -> String
}

app := ():A => implement A {
	
}
`, "missing declaration for variable a")
}

func TestInterfaceWithSeparateImplementationWrongVariableType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
  a: () -> Void
}

app := (): A => implement A {
  a := () => ""
}
`, "expected type Void but found String")
}
