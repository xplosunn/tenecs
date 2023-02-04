package testcode

const Struct TestCodeCategory = "struct"

var StructWithConstructorEmpty = Create(Struct, "StructWithConstructorEmpty", `
package main

struct NoArgsStruct()
`)

var StructWithConstructorWithString = Create(Struct, "StructWithConstructorWithString", `
package main

struct StringStruct(str: String)
`)

var StructWithConstructorWithBooleans = Create(Struct, "StructWithConstructorWithBooleans", `
package main

struct BooleanColor(red: Boolean, green: Boolean, blue: Boolean)
`)

var StructWithConstructorAnotherStruct1 = Create(Struct, "StructWithConstructorAnotherStruct1", `
package main

struct Outer(c: Inner)
struct Inner()
`)

var StructWithConstructorAnotherStruct2 = Create(Struct, "StructWithConstructorAnotherStruct2", `
package main

struct Inner()
struct Outer(c: Inner)
`)

var StructAsVariable = Create(Struct, "StructAsVariable", `
package main

import tenecs.os.Main

struct Person(name: String)

app := (): Main => implement Main {
	public main := (runtime) => {
		me := Person("Author")
	}
}
`)
