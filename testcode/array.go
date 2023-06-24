package testcode

const Array TestCodeCategory = "array"

var ArrayVariableWithEmptyArray = Create(Array, "ArrayVariableWithEmptyArray", `
package main

noStrings := [ String ] ( )
`)

var ArrayVariableWithTwoElementArray = Create(Array, "ArrayVariableWithTwoElementArray", `
package main

someStrings := [ String ] ( "a" , "b" )
`)

var ArrayOfArray = Create(Array, "ArrayOfArray", `
package main

someStrings := [Array<String>]([String]("a" , "b"))
`)
