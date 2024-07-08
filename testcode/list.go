package testcode

const List TestCodeCategory = "list"

var ListVariableWithEmptyList = Create(List, "ListVariableWithEmptyList", `
package main

noStrings := [ String ] ( )
`)

var ListVariableWithTwoElementList = Create(List, "ListVariableWithTwoElementList", `
package main

someStrings := [ String ] ( "a" , "b" )
`)

var ListOfList = Create(List, "ListOfList", `
package main

someStrings := [List<String>]([String]("a" , "b"))
`)
